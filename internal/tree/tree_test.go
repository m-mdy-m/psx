package tree

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

// buildFS materialises a fstest.MapFS into a real temp dir, since Scan walks the OS.
func buildFS(t *testing.T, files map[string]string) string {
	t.Helper()
	fsys := fstest.MapFS{}
	for name, body := range files {
		fsys[name] = &fstest.MapFile{Data: []byte(body)}
	}
	root := t.TempDir()
	for name, f := range fsys {
		full := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, f.Data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestScanIndexesFilesAndDirs(t *testing.T) {
	root := buildFS(t, map[string]string{
		"README.md":               "hello",
		"src/index.ts":            "x",
		"docs/adr/0001-init.md":   "y",
		"docs/architecture/x.md":  "y",
		"node_modules/left-pad/a": "junk",
		".git/config":             "[core]",
	})

	snap, err := Scan(root, []string{"node_modules/", ".git/"})
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}

	if !snap.Exists("README.md") {
		t.Error("README.md should be indexed")
	}
	if !snap.Exists("src") {
		t.Error("src should be indexed as a dir")
	}
	if snap.Exists("node_modules") {
		t.Error("ignored dir must not be indexed")
	}
	if snap.Exists(".git") {
		t.Error("ignored dir must not be indexed")
	}
	if got := snap.Size("README.md"); got != 5 {
		t.Errorf("Size(README.md) = %d, want 5", got)
	}
}

func TestScanSkipsEmptyDirs(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "empty"), 0o755); err != nil {
		t.Fatal(err)
	}

	snap, err := Scan(root, nil)
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}

	if snap.NonEmpty("empty") {
		t.Error("empty dir must not satisfy a non-empty requirement")
	}
	if !snap.HasDir("empty") {
		t.Error("empty dir must still be reported as existing")
	}
}

func TestGlobDoubleStarRecurses(t *testing.T) {
	root := buildFS(t, map[string]string{
		"top_test.go":                   "x",
		"internal/rules/engine_test.go": "x",
		"a/b/c/deep_test.go":            "x",
		"internal/rules/engine.go":      "x",
	})

	snap, err := Scan(root, nil)
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}

	// The old filepath.Glob returned 0 matches for this pattern.
	got := snap.Glob("**/*_test.go")
	if len(got) != 3 {
		t.Errorf("Glob(**/*_test.go) = %v, want 3 matches", got)
	}
	if !snap.Match("**/*_test.go") {
		t.Error("Match(**/*_test.go) should be true")
	}
	if snap.Match("**/*.rb") {
		t.Error("Match(**/*.rb) should be false")
	}
}

func TestGlobBraceExpansion(t *testing.T) {
	root := buildFS(t, map[string]string{
		"src/a.test.js": "x",
		"src/a.test.ts": "x",
		"src/a.spec.ts": "x",
		"src/plain.js":  "x",
	})

	snap, _ := Scan(root, nil)
	if got := snap.Glob("**/*.{test,spec}.{js,ts}"); len(got) != 3 {
		t.Errorf("brace glob = %v, want 3", got)
	}
}

func TestGlobIsSortedForStableOutput(t *testing.T) {
	root := buildFS(t, map[string]string{
		"z_test.go": "x", "a_test.go": "x", "m/n_test.go": "x",
	})

	snap, _ := Scan(root, nil)
	first := snap.Glob("**/*_test.go")
	for i := 0; i < 20; i++ {
		got := snap.Glob("**/*_test.go")
		for j := range got {
			if got[j] != first[j] {
				t.Fatalf("Glob is not deterministic: %v vs %v", got, first)
			}
		}
	}
}

func TestIgnorePatterns(t *testing.T) {
	files := map[string]string{
		"dist/a.js":          "x",
		"vendor/a.js":        "x",
		"app.log":            "x",
		"tmpdir/a.js":        "x",
		"src/generated/x.js": "x",
		"snapshots/ui.snap":  "x",
		"src/index.ts":       "x",
	}

	cases := []struct {
		name   string
		ignore []string
		gone   string // path that must be excluded
		kept   string // path that must survive
	}{
		{"exact dir with slash", []string{"dist/"}, "dist/a.js", "src/index.ts"},
		{"exact dir no slash", []string{"dist"}, "dist/a.js", "src/index.ts"},
		{"wildcard", []string{"*.log"}, "app.log", "src/index.ts"},
		{"prefix wildcard", []string{"tmp*"}, "tmpdir/a.js", "src/index.ts"},
		{"nested glob", []string{"src/generated/**"}, "src/generated/x.js", "src/index.ts"},
		{"deep star", []string{"**/*.snap"}, "snapshots/ui.snap", "src/index.ts"},
		{"dir glob does not hide siblings", []string{"dist/*"}, "dist/a.js", "src/index.ts"},
		{"unrelated pattern", []string{"vendor/"}, "vendor/a.js", "dist/a.js"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			snap, err := Scan(buildFS(t, files), tc.ignore)
			if err != nil {
				t.Fatalf("Scan: %v", err)
			}
			if snap.Exists(tc.gone) {
				t.Errorf("%s should be excluded by %v", tc.gone, tc.ignore)
			}
			if !snap.Exists(tc.kept) {
				t.Errorf("%s must survive ignore %v", tc.kept, tc.ignore)
			}
		})
	}
}

func TestNegatedIgnore(t *testing.T) {
	root := buildFS(t, map[string]string{
		"logs/keep.txt": "x",
		"logs/drop.txt": "x",
	})

	// gitignore cannot re-include a file under an excluded directory, so the
	// broad rule must target the directory's contents rather than the dir.
	snap, err := Scan(root, []string{"logs/*", "!logs/keep.txt"})
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if !snap.Exists("logs/keep.txt") {
		t.Error("negated pattern must re-include the file")
	}
	if snap.Exists("logs/drop.txt") {
		t.Error("logs/drop.txt should stay ignored")
	}
}

func TestHasFileRejectsEmptyFile(t *testing.T) {
	root := buildFS(t, map[string]string{"empty.md": "", "full.md": "x"})
	snap, _ := Scan(root, nil)

	if snap.HasFile("empty.md") {
		t.Error("HasFile must report a zero-byte file as absent")
	}
	if !snap.HasFile("full.md") {
		t.Error("HasFile must accept a file with content")
	}
	if snap.HasDir("nope") {
		t.Error("HasDir must be false for a missing path")
	}
	if snap.HasDir("full.md") {
		t.Error("HasDir must be false for a file")
	}
}
