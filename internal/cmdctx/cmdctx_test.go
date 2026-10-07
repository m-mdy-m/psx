package cmdctx

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// stage creates a project directory with the given files.
func stage(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for name, body := range files {
		full := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// TestConfigFlagOverridesDiscovery is the regression test for --config being
// parsed and then ignored.
//
// The flag was registered on the root command and threaded into the options
// struct, but cmdctx.Load passed "" to the loader, so every invocation used
// whichever config it discovered instead.
func TestConfigFlagOverridesDiscovery(t *testing.T) {
	// A psx.yml in the project that enables a single rule.
	root := stage(t, map[string]string{
		"go.mod":  "module demo\n",
		"psx.yml": "version: 1\nproject:\n  type: go\nrules:\n  readme: error\n",
	})

	// The flag points somewhere else entirely.
	explicit := filepath.Join(t.TempDir(), "other.yml")
	body := "version: 1\nproject:\n  type: nodejs\nrules:\n  license: error\n"
	if err := os.WriteFile(explicit, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	ctx, err := Load(root, explicit, ModeReadOnly)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if ctx.Config.ConfigFile != explicit {
		t.Errorf("ConfigFile = %q, want %q", ctx.Config.ConfigFile, explicit)
	}
	if ctx.ProjectType != "nodejs" {
		t.Errorf("ProjectType = %q, want nodejs from the explicit config", ctx.ProjectType)
	}
	if _, ok := ctx.Config.ActiveRules["license"]; !ok {
		t.Error("the explicit config's rules were not applied")
	}
	if _, ok := ctx.Config.ActiveRules["readme"]; ok {
		t.Error("the discovered psx.yml was used despite --config")
	}
}

// TestDiscoveredConfigIsUsedWhenNoFlag covers the normal path.
func TestDiscoveredConfigIsUsedWhenNoFlag(t *testing.T) {
	root := stage(t, map[string]string{
		"go.mod":  "module demo\n",
		"psx.yml": "version: 1\nproject:\n  type: go\nrules:\n  readme: error\n",
	})

	ctx, err := Load(root, "", ModeReadOnly)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !strings.HasSuffix(ctx.Config.ConfigFile, "psx.yml") {
		t.Errorf("ConfigFile = %q, want the discovered psx.yml", ctx.Config.ConfigFile)
	}
	if ctx.ProjectType != "go" {
		t.Errorf("ProjectType = %q, want go", ctx.ProjectType)
	}
}

// TestMissingConfigFileIsAnError keeps a typo in --config from silently
// falling back to the defaults.
func TestMissingConfigFileIsAnError(t *testing.T) {
	root := stage(t, map[string]string{"go.mod": "module demo\n"})

	_, err := Load(root, filepath.Join(root, "nope.yml"), ModeReadOnly)
	if err == nil {
		t.Fatal("expected an error for a config file that does not exist")
	}
	if !strings.Contains(err.Error(), "nope.yml") {
		t.Errorf("error should name the missing file, got: %v", err)
	}
}

// TestReadOnlyModeNeverWrites pins the guarantee that check does not touch disk.
func TestReadOnlyModeNeverWrites(t *testing.T) {
	root := stage(t, map[string]string{"go.mod": "module demo\n"})

	if _, err := Load(root, "", ModeReadOnly); err != nil {
		t.Fatalf("Load: %v", err)
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		var names []string
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Errorf("read-only load created files: %v", names)
	}

	if _, err := os.Stat(filepath.Join(root, ".psx-project.yml")); err == nil {
		t.Error("read-only load wrote .psx-project.yml")
	}
}

// TestInvalidConfigFailsLoudly ensures a broken config does not degrade to
// the defaults, which would silently check the wrong rule set.
func TestInvalidConfigFailsLoudly(t *testing.T) {
	root := stage(t, map[string]string{"go.mod": "module demo\n"})
	bad := filepath.Join(t.TempDir(), "bad.yml")

	if err := os.WriteFile(bad, []byte("version: 1\nrules:\n  readme: nonsense\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(root, bad, ModeReadOnly)
	if err == nil {
		t.Fatal("expected an error for an invalid severity")
	}
}

// TestExplicitTypeIsNormalised confirms an alias in a config file resolves.
func TestExplicitTypeIsNormalised(t *testing.T) {
	root := stage(t, map[string]string{"package.json": "{}"})
	cfgPath := filepath.Join(root, "psx.yml")

	if err := os.WriteFile(cfgPath,
		[]byte("version: 1\nproject:\n  type: ts\nrules:\n  readme: error\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	ctx, err := Load(root, "", ModeReadOnly)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if ctx.ProjectType != "nodejs" {
		t.Errorf("ProjectType = %q, want nodejs (ts must normalise)", ctx.ProjectType)
	}
	if ctx.Config.Project.Type != "nodejs" {
		t.Errorf("Config.Project.Type = %q, want nodejs", ctx.Config.Project.Type)
	}
}

// TestDetectionRunsWhenTypeIsUnset is the regression test for detection never
// running.
//
// NormalizeProjectType("") returns "generic", and buildConfig wrote that back onto
// cfg.Project.Type. cmdctx then saw a non-empty declared type, concluded the user had
// asked for generic, and skipped detection — so every project with an unset
// project.type was reported as generic regardless of its manifests.
func TestDetectionRunsWhenTypeIsUnset(t *testing.T) {
	cases := []struct {
		file string
		want string
	}{
		{"package.json", "nodejs"},
		{"go.mod", "go"},
		{"Cargo.toml", "rust"},
		{"pyproject.toml", "python"},
	}

	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			root := stage(t, map[string]string{
				tc.file:   "manifest\n",
				"psx.yml": "version: 1\nproject:\n  type: \"\"\nrules:\n  readme: error\n",
			})

			ctx, err := Load(root, "", ModeReadOnly)
			if err != nil {
				t.Fatalf("Load: %v", err)
			}
			if !ctx.Detected {
				t.Error("Detected = false, want true when project.type is unset")
			}
			if ctx.ProjectType != tc.want {
				t.Errorf("ProjectType = %q, want %q (from %s)", ctx.ProjectType, tc.want, tc.file)
			}
		})
	}
}

// TestExplicitGenericIsNotOverriddenByDetection confirms an explicit choice wins.
func TestExplicitGenericIsNotOverriddenByDetection(t *testing.T) {
	root := stage(t, map[string]string{
		"package.json": "{}\n",
		"psx.yml":      "version: 1\nproject:\n  type: generic\nrules:\n  readme: error\n",
	})

	ctx, err := Load(root, "", ModeReadOnly)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if ctx.Detected {
		t.Error("Detected = true, want false when project.type is set explicitly")
	}
	if ctx.ProjectType != "generic" {
		t.Errorf("ProjectType = %q, want generic", ctx.ProjectType)
	}
}

// TestAutoPlaceholderTriggersDetection covers the explicit "auto" spelling.
func TestAutoPlaceholderTriggersDetection(t *testing.T) {
	root := stage(t, map[string]string{
		"go.mod":  "module demo\n",
		"psx.yml": "version: 1\nproject:\n  type: auto\nrules:\n  readme: error\n",
	})

	ctx, err := Load(root, "", ModeReadOnly)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if ctx.ProjectType != "go" {
		t.Errorf("ProjectType = %q, want go; 'auto' must trigger detection", ctx.ProjectType)
	}
}

// TestConfigPathIsNotReportedAsTheProjectDirectory guards the verbose output,
// which used to print the project directory under the label "Config".
func TestConfigPathIsNotReportedAsTheProjectDirectory(t *testing.T) {
	root := stage(t, map[string]string{"go.mod": "module demo\n"})

	ctx, err := Load(root, "", ModeReadOnly)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if ctx.Config.ConfigFile == ctx.Path {
		t.Error("ConfigFile points at the project directory, not the config file")
	}
	if ctx.Config.ProjectPath != ctx.Path {
		t.Errorf("ProjectPath = %q, want %q", ctx.Config.ProjectPath, ctx.Path)
	}
}
