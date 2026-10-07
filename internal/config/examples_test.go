package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
)

// examplesDir holds the shipped example configurations.
const examplesDir = "../../examples"

// exampleFiles lists the example configs, excluding the baseline file and the
// index, neither of which is a configuration.
func exampleFiles(t *testing.T) []string {
	t.Helper()

	entries, err := os.ReadDir(examplesDir)
	if err != nil {
		t.Fatalf("read examples dir: %v", err)
	}

	var out []string
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".yml") {
			continue
		}
		// baseline.txt records accepted violations; it is neither a config nor YAML.
		if name == "baseline.txt" {
			continue
		}
		out = append(out, name)
	}
	return out
}

// TestExamplesAreValidConfigs loads every example.
//
// An example that does not load is worse than no example: it is the first thing a
// user copies, so it has to work.
func TestExamplesAreValidConfigs(t *testing.T) {
	files := exampleFiles(t)
	if len(files) == 0 {
		t.Fatal("no example configs found")
	}

	meta := GetRulesMetadata().Rules

	for _, name := range files {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(examplesDir, name)

			var cfg Config
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read: %v", err)
			}
			if err := yaml.Unmarshal(data, &cfg); err != nil {
				t.Fatalf("parse: %v", err)
			}

			if res := Validate(&cfg); !IsValid(res) {
				for _, e := range res.Errors {
					t.Errorf("invalid: [%s] %s", e.Field, e.Message)
				}
				t.Fatal("example does not validate")
			}

			if len(cfg.Rules) == 0 {
				t.Error("example enables no rules")
			}

			for id, sev := range cfg.Rules {
				if _, ok := meta[id]; !ok {
					t.Errorf("unknown rule id %q; run `psx rules` for the list", id)
				}
				if _, err := ParseSeverity(sev, config0Severity()); err != nil {
					t.Errorf("rule %s: %v", id, err)
				}
			}
		})
	}
}

// config0Severity is the fallback severity for validating an example's values.
func config0Severity() Severity { return SeverityWarning }

// TestExamplesUseOnlyKnownSeverities is a stricter pass over severity strings,
// including the boolean form used to disable a rule.
func TestExamplesUseOnlyKnownSeverities(t *testing.T) {
	for _, name := range exampleFiles(t) {
		data, err := os.ReadFile(filepath.Join(examplesDir, name))
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}

		var raw struct {
			Rules map[string]any `yaml:"rules"`
		}
		if err := yaml.Unmarshal(data, &raw); err != nil {
			t.Fatalf("%s: %v", name, err)
		}

		for id, v := range raw.Rules {
			switch val := v.(type) {
			case string:
				if _, err := ParseSeverity(val, SeverityInfo); err != nil {
					t.Errorf("%s: rule %s: %v", name, id, err)
				}
			case bool:
				if val {
					t.Errorf("%s: rule %s: true is not a severity; omit the rule or use false", name, id)
				}
			default:
				t.Errorf("%s: rule %s: severity must be a string or false, got %T", name, id, v)
			}
		}
	}
}

// TestExampleIgnorePatternsAreValid checks each ignore entry is a usable pattern.
//
// An entry that is only whitespace, or a bare "*", silently disables the whole scan.
func TestExampleIgnorePatternsAreValid(t *testing.T) {
	for _, name := range exampleFiles(t) {
		data, err := os.ReadFile(filepath.Join(examplesDir, name))
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}

		var cfg struct {
			Ignore []string `yaml:"ignore"`
			Custom *struct {
				Files []struct {
					Path string `yaml:"path"`
				} `yaml:"files"`
				Folders []struct {
					Path string `yaml:"path"`
				} `yaml:"folders"`
			} `yaml:"custom"`
		}
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			t.Fatalf("%s: %v", name, err)
		}

		for i, p := range cfg.Ignore {
			if strings.TrimSpace(p) == "" {
				t.Errorf("%s: ignore[%d] is empty", name, i)
			}
			if p == "*" || p == "**" {
				t.Errorf("%s: ignore[%d] %q would exclude the entire project", name, i, p)
			}
			if filepath.IsAbs(p) || strings.HasPrefix(p, "/") && strings.Count(p, "/") == 1 && len(p) > 1 &&
				!strings.ContainsAny(p, "*?[]") && p != "/" {
				// A single leading slash anchors to the root, which is valid.
				continue
			}
		}

		// Custom paths must be relative and stay inside the project.
		if cfg.Custom != nil {
			for _, f := range cfg.Custom.Files {
				assertRelativeExamplePath(t, name, f.Path)
			}
			for _, f := range cfg.Custom.Folders {
				assertRelativeExamplePath(t, name, f.Path)
			}
		}
	}
}

// assertRelativeExamplePath rejects a custom path that escapes the project.
func assertRelativeExamplePath(t *testing.T, file, path string) {
	t.Helper()

	switch {
	case strings.TrimSpace(path) == "":
		t.Errorf("%s: custom path is empty", file)
	case filepath.IsAbs(path) || strings.HasPrefix(path, "/") || strings.HasPrefix(path, `\`):
		t.Errorf("%s: custom path %q must be relative", file, path)
	case strings.Contains(path, ".."):
		t.Errorf("%s: custom path %q escapes the project", file, path)
	}
}

// TestBaselineFileIsUsable checks the shipped baseline names real rules.
//
// A baseline is a promise that the debt is known; a typo in it leaves a rule
// failing with no way to tell why.
func TestBaselineFileIsUsable(t *testing.T) {
	path := filepath.Join(examplesDir, "baseline.txt")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Skipf("no baseline example: %v", err)
	}

	meta := GetRulesMetadata().Rules
	seen := 0

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		seen++
		if _, ok := meta[line]; !ok {
			t.Errorf("baseline lists unknown rule %q", line)
		}
	}

	if seen == 0 {
		t.Error("baseline lists no rules")
	}
}
