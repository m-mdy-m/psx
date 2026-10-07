package rules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/m-mdy-m/psx/internal/config"
)

func customFolderEnv(t *testing.T, structure map[string]any) *fixEnv {
	t.Helper()
	env := newFixEnv(t, map[string]string{}, "generic", nil)
	env.cfg.Custom = &config.CustomConfig{
		Folders: []config.CustomFolder{{Path: "config", Structure: structure}},
	}
	return env
}

// A key with no children is a file. Creating it as a directory made every
// `production.yaml: {}` entry a directory named like a config file, which nothing
// downstream could read.
func TestCustomFolderLeafIsAFileNotADirectory(t *testing.T) {
	env := customFolderEnv(t, map[string]any{
		"base": map[string]any{
			"production.yaml": map[string]any{},
			"staging.yaml":    "",
		},
	})

	results := CustomHandler{root: env.root}.Apply(env.cfg.Custom, &FixContext{Context: env.ctx})
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("apply folders: %v", r.Error)
		}
	}

	for _, rel := range []string{"config/base/production.yaml", "config/base/staging.yaml"} {
		info, err := os.Stat(filepath.Join(env.root, rel))
		if err != nil {
			t.Fatalf("%s should exist: %v", rel, err)
		}
		if info.IsDir() {
			t.Errorf("%s was created as a directory, want a file", rel)
		}
	}
}

// A key that has children is a directory, and must still be created as one.
func TestCustomFolderBranchIsADirectory(t *testing.T) {
	env := customFolderEnv(t, map[string]any{
		"base": map[string]any{
			"production.yaml": map[string]any{},
		},
	})

	if err := applyOne(t, env, "config/base"); err != nil {
		t.Fatalf("apply folders: %v", err)
	}

	info, err := os.Stat(filepath.Join(env.root, "config", "base"))
	if err != nil {
		t.Fatalf("config/base should exist: %v", err)
	}
	if !info.IsDir() {
		t.Error("config/base has children, so it should be a directory")
	}
}

func TestCustomFolderLeafContentIsWritten(t *testing.T) {
	env := customFolderEnv(t, map[string]any{
		"Makefile": "all:\n\techo hi\n",
	})

	if err := applyOne(t, env, "config/Makefile"); err != nil {
		t.Fatalf("apply folders: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(env.root, "config", "Makefile"))
	if err != nil {
		t.Fatalf("config/Makefile should exist: %v", err)
	}
	if string(got) != "all:\n\techo hi\n" {
		t.Errorf("content = %q, want the configured text", got)
	}
}

// Creation order must not depend on Go's map iteration.
func TestCustomFolderCreationOrderIsStable(t *testing.T) {
	structure := map[string]any{
		"z": map[string]any{"b.txt": "", "a.txt": ""},
		"a": map[string]any{"d.txt": "", "c.txt": ""},
		"m": map[string]any{"f.txt": "", "e.txt": ""},
	}

	first := flattenStructure("config", structure)
	for range 20 {
		got := flattenStructure("config", structure)
		if len(got) != len(first) {
			t.Fatalf("length changed between runs: %d then %d", len(first), len(got))
		}
		for i := range first {
			if got[i].Path != first[i].Path {
				t.Fatalf("order changed at %d: %q then %q", i, first[i].Path, got[i].Path)
			}
		}
	}
}

func applyOne(t *testing.T, env *fixEnv, rel string) error {
	t.Helper()
	results := CustomHandler{root: env.root}.Apply(env.cfg.Custom, &FixContext{Context: env.ctx})
	for _, r := range results {
		if r.Error != nil {
			return r.Error
		}
		_ = rel
	}
	return nil
}
