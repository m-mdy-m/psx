package detect

import (
	"os"
	"path/filepath"
	"testing"
)

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

func TestDetectLanguageFromItsManifest(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		want    string
		wantsig string
	}{
		{"go", "go.mod", "go", "go.mod"},
		{"rust", "Cargo.toml", "rust", "Cargo.toml"},
		{"python", "pyproject.toml", "python", "pyproject.toml"},
		{"python setup", "setup.py", "python", "setup.py"},
		{"python requirements", "requirements.txt", "python", "requirements.txt"},
		{"node", "package.json", "nodejs", "package.json"},
		{"deno", "deno.json", "nodejs", "deno.json"},
		{"php", "composer.json", "php", "composer.json"},
		{"ruby", "Gemfile", "ruby", "Gemfile"},
		{"java maven", "pom.xml", "java", "pom.xml"},
		{"java gradle", "build.gradle", "java", "build.gradle"},
		{"c", "CMakeLists.txt", "c", "CMakeLists.txt"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p, err := Detect(stage(t, map[string]string{tc.file: "x"}), nil)
			if err != nil {
				t.Fatal(err)
			}
			if p.ProjectType != tc.want {
				t.Errorf("ProjectType = %q, want %q", p.ProjectType, tc.want)
			}
			if len(p.Signals) == 0 || p.Signals[0] != tc.wantsig {
				t.Errorf("Signals = %v, want %q first", p.Signals, tc.wantsig)
			}
			if p.Confidence < 0.7 {
				t.Errorf("Confidence = %v, a recognised language should not be low", p.Confidence)
			}
		})
	}
}

func TestDetectFallsBackToGeneric(t *testing.T) {
	p, err := Detect(stage(t, map[string]string{"notes.txt": "hello"}), nil)
	if err != nil {
		t.Fatal(err)
	}
	if p.ProjectType != "generic" {
		t.Errorf("ProjectType = %q, want generic", p.ProjectType)
	}
	if p.Kind != "app" {
		t.Errorf("Kind = %q, want app", p.Kind)
	}
	if p.Confidence >= 0.7 {
		t.Errorf("Confidence = %v, an unrecognised project should not claim confidence", p.Confidence)
	}
	if len(p.Signals) != 0 {
		t.Errorf("Signals = %v, want none for an unrecognised project", p.Signals)
	}
}

func TestDetectEmptyDirectory(t *testing.T) {
	p, err := Detect(t.TempDir(), nil)
	if err != nil {
		t.Fatalf("detecting an empty directory should succeed, not error: %v", err)
	}
	if p.ProjectType != "generic" || p.Kind != "app" {
		t.Errorf("got %+v, want the generic app fallback", p)
	}
}

// The markers list claims to be ordered by specificity. It is not: nodejs scores
// 0.95 but sits below python's 0.70, so a repository carrying both manifests was
// reported as Python.
func TestTheMostConfidentMarkerWins(t *testing.T) {
	p, err := Detect(stage(t, map[string]string{
		"requirements.txt": "flask\n",
		"package.json":     "{}",
	}), nil)
	if err != nil {
		t.Fatal(err)
	}
	if p.ProjectType != "nodejs" {
		t.Errorf("ProjectType = %q, want nodejs: package.json (0.95) outranks requirements.txt (0.70)",
			p.ProjectType)
	}
}

func TestTypescriptRaisesNodeConfidence(t *testing.T) {
	p, err := Detect(stage(t, map[string]string{
		"package.json":  "{}",
		"tsconfig.json": "{}",
	}), nil)
	if err != nil {
		t.Fatal(err)
	}
	if p.ProjectType != "nodejs" {
		t.Errorf("ProjectType = %q, want nodejs", p.ProjectType)
	}
	if p.Confidence < 0.95 {
		t.Errorf("Confidence = %v, tsconfig should raise it above the bare package.json score", p.Confidence)
	}
	found := false
	for _, s := range p.Signals {
		if s == "tsconfig.json" {
			found = true
		}
	}
	if !found {
		t.Errorf("Signals = %v, want tsconfig.json recorded as the reason", p.Signals)
	}
}

func TestKindFromDirectoryLayout(t *testing.T) {
	tests := []struct {
		name  string
		files map[string]string
		want  string
	}{
		{"microservice", map[string]string{"go.mod": "m", "services/api/main.go": "package api"}, "microservice"},
		{"monorepo apps", map[string]string{"go.mod": "m", "apps/web/main.go": "p"}, "monorepo"},
		{"monorepo packages", map[string]string{"go.mod": "m", "packages/core/lib.go": "p"}, "monorepo"},
		{"monorepo crates", map[string]string{"Cargo.toml": "x", "crates/core/src/lib.rs": "p"}, "monorepo"},
		{"plugin extensions", map[string]string{"package.json": "{}", "extensions/a/ext.json": "{}"}, "plugin"},
		{"plugin plugins", map[string]string{"package.json": "{}", "plugins/a/ext.json": "{}"}, "plugin"},
		{"plugin addons", map[string]string{"package.json": "{}", "addons/a/ext.json": "{}"}, "plugin"},
		{"cli", map[string]string{"go.mod": "m", "cmd/tool/main.go": "package main"}, "cli"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p, err := Detect(stage(t, tc.files), nil)
			if err != nil {
				t.Fatal(err)
			}
			if p.Kind != tc.want {
				t.Errorf("Kind = %q, want %q", p.Kind, tc.want)
			}
		})
	}
}

func TestWorkspaceFileBeatsDirectoryLayout(t *testing.T) {
	p, err := Detect(stage(t, map[string]string{
		"go.mod":           "m",
		"go.work":          "go 1.25",
		"cmd/tool/main.go": "package main",
	}), nil)
	if err != nil {
		t.Fatal(err)
	}
	if p.Kind != "monorepo" {
		t.Errorf("Kind = %q, want monorepo: go.work is an explicit statement, a cmd/ directory is a hint",
			p.Kind)
	}
}

func TestKindFromARootMainFile(t *testing.T) {
	p, err := Detect(stage(t, map[string]string{"go.mod": "m", "main.go": "package main"}), nil)
	if err != nil {
		t.Fatal(err)
	}
	if p.Kind != "cli" {
		t.Errorf("Kind = %q, want cli for a root-level main.go", p.Kind)
	}
}

func TestKindFromARootIndexFile(t *testing.T) {
	for _, entry := range []string{"index.js", "index.ts"} {
		t.Run(entry, func(t *testing.T) {
			p, err := Detect(stage(t, map[string]string{"package.json": "{}", entry: ""}), nil)
			if err != nil {
				t.Fatal(err)
			}
			if p.Kind != "library" {
				t.Errorf("Kind = %q, want library for a root-level %s", p.Kind, entry)
			}
		})
	}
}

// services/ implies a microservice, which is a stronger signal than the bare
// app fallback, so a root index file must not downgrade it to a library.
func TestKindMarkersWinOverEntryPointHeuristics(t *testing.T) {
	p, err := Detect(stage(t, map[string]string{
		"package.json":      "{}",
		"index.js":          "",
		"services/api/i.js": "",
	}), nil)
	if err != nil {
		t.Fatal(err)
	}
	if p.Kind != "microservice" {
		t.Errorf("Kind = %q, want microservice: services/ outranks a root index file", p.Kind)
	}
}

func TestWorkspaceGlobsAreOnlyForMultiPackageLayouts(t *testing.T) {
	tests := []struct {
		kind string
		want bool
	}{
		{"monorepo", true},
		{"microservice", true},
		{"plugin", true},
		{"app", false},
		{"cli", false},
		{"library", false},
	}
	for _, tc := range tests {
		t.Run(tc.kind, func(t *testing.T) {
			got := workspaceGlobs(&Profile{Kind: tc.kind, ProjectType: "nodejs"})
			if (len(got) > 0) != tc.want {
				t.Errorf("workspaceGlobs(%s) = %v, want present=%v", tc.kind, got, tc.want)
			}
		})
	}
}

func TestWorkspaceGlobsMatchTheDetectedLanguage(t *testing.T) {
	tests := []struct {
		projectType string
		kind        string
		want        []string
	}{
		{"go", "monorepo", []string{"./..."}},
		{"rust", "monorepo", []string{"crates/*"}},
		{"nodejs", "monorepo", []string{"apps/*", "packages/*"}},
		{"python", "monorepo", []string{"packages/*"}},
	}
	for _, tc := range tests {
		t.Run(tc.projectType, func(t *testing.T) {
			got := workspaceGlobs(&Profile{Kind: tc.kind, ProjectType: tc.projectType})
			if len(got) != len(tc.want) {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
			for i := range tc.want {
				if got[i] != tc.want[i] {
					t.Errorf("got %v, want %v", got, tc.want)
				}
			}
		})
	}
}

// A microservice workspace adds its services to the language's own layout.
func TestMicroserviceGlobsIncludeServices(t *testing.T) {
	got := workspaceGlobs(&Profile{Kind: "microservice", ProjectType: "go"})
	found := false
	for _, g := range got {
		if g == "services/*" {
			found = true
		}
	}
	if !found {
		t.Errorf("globs = %v, want services/* included", got)
	}
}

func TestWorkspaceGlobsAreSorted(t *testing.T) {
	got := workspaceGlobs(&Profile{Kind: "microservice", ProjectType: "nodejs"})
	for i := 1; i < len(got); i++ {
		if got[i-1] > got[i] {
			t.Errorf("globs %v are not sorted; output would differ between runs", got)
			break
		}
	}
}

func TestKindsAndKnown(t *testing.T) {
	kinds := Kinds()
	if len(kinds) == 0 {
		t.Fatal("Kinds() returned nothing")
	}
	for _, k := range kinds {
		if !Known(k) {
			t.Errorf("Known(%q) = false, but Kinds lists it", k)
		}
	}
	for _, k := range []string{"", "CLI", "app ", "unknown"} {
		if Known(k) {
			t.Errorf("Known(%q) = true, want false", k)
		}
	}
}

// Kinds is the vocabulary the --kind flag validates against, so it must not
// drift from what detection and init actually produce.
func TestKindsCoverEveryKindDetectionCanProduce(t *testing.T) {
	layouts := []map[string]string{
		{"notes.txt": "x"},
		{"go.mod": "m", "main.go": "package main"},
		{"package.json": "{}", "index.js": ""},
		{"go.mod": "m", "cmd/tool/main.go": "p"},
		{"go.mod": "m", "apps/web/main.go": "p"},
		{"go.mod": "m", "services/api/main.go": "p"},
		{"package.json": "{}", "extensions/a/e.json": "{}"},
	}
	for _, files := range layouts {
		p, err := Detect(stage(t, files), nil)
		if err != nil {
			t.Fatal(err)
		}
		if !Known(p.Kind) {
			t.Errorf("detected kind %q from %v, which Kinds() does not list", p.Kind, files)
		}
	}
}

func TestDetectHonoursIgnorePatterns(t *testing.T) {
	root := stage(t, map[string]string{
		"go.mod":            "module x\n",
		"vendor/lib/go.mod": "module lib\n",
	})
	p, err := Detect(root, []string{"vendor/"})
	if err != nil {
		t.Fatal(err)
	}
	if p.ProjectType != "go" {
		t.Errorf("ProjectType = %q, want go", p.ProjectType)
	}
}

func TestDetectOnAMissingDirectory(t *testing.T) {
	if _, err := Detect(filepath.Join(t.TempDir(), "nope"), nil); err == nil {
		t.Error("detecting a missing directory should report an error")
	}
}

func TestDetectionIsDeterministic(t *testing.T) {
	files := map[string]string{
		"go.mod":                        "m",
		"cmd/tool/main.go":              "package main",
		"services/api/main.go":          "package api",
		"packages/core/lib.go":          "package core",
		"apps/web/main.go":              "package web",
		"extensions/vscode/ext.json":    "{}",
		"plugins/other/ext.json":        "{}",
		"addons/third/ext.json":         "{}",
		"crates/rust/Cargo.toml":        "x",
		"node_modules/dep/package.json": "{}",
	}
	root := stage(t, files)

	first, err := Detect(root, nil)
	if err != nil {
		t.Fatal(err)
	}
	for range 20 {
		again, err := Detect(root, nil)
		if err != nil {
			t.Fatal(err)
		}
		if again.ProjectType != first.ProjectType || again.Kind != first.Kind {
			t.Fatalf("detection is not deterministic: %+v then %+v", first, again)
		}
		if len(again.Signals) != len(first.Signals) {
			t.Fatalf("signal count changed between runs: %v then %v", first.Signals, again.Signals)
		}
		for i := range first.Signals {
			if again.Signals[i] != first.Signals[i] {
				t.Fatalf("signal order changed: %v then %v", first.Signals, again.Signals)
			}
		}
	}
}

func TestHasPrefixDir(t *testing.T) {
	dirs := map[string]bool{"cmd/tool": true, "src": true}
	if !hasPrefixDir(dirs, "cmd/") {
		t.Error("hasPrefixDir should match a nested directory under cmd/")
	}
	if hasPrefixDir(dirs, "services/") {
		t.Error("hasPrefixDir matched a prefix that is absent")
	}
	if hasPrefixDir(nil, "cmd/") {
		t.Error("hasPrefixDir on an empty set should be false")
	}
}

// A directory named go.mod is not a Go module manifest.
func TestADirectoryDoesNotCountAsAManifest(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "go.mod"), 0o755); err != nil {
		t.Fatal(err)
	}
	p, err := Detect(root, nil)
	if err != nil {
		t.Fatal(err)
	}
	if p.ProjectType != "generic" {
		t.Errorf("ProjectType = %q, want generic: a directory named go.mod is not a manifest", p.ProjectType)
	}
}
