package tree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// The fingerprint lets a watcher decide whether anything changed without walking
// the tree a second time, so it must agree with the index it came from.
func TestFingerprintCountsOnlyFiles(t *testing.T) {
	root := buildFS(t, map[string]string{
		"go.mod":           "module x\n",
		"cmd/main.go":      "package main\n",
		"docs/readme.md":   "# hi\n",
		"docs/adr/note.md": "n\n",
	})

	snap, err := Scan(root, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, files := snap.Fingerprint()
	if files != 4 {
		t.Errorf("Fingerprint reported %d files, want 4 (directories are not files)", files)
	}
}

func TestFingerprintIsStableAcrossRepeatedCalls(t *testing.T) {
	root := buildFS(t, map[string]string{
		"a.go":        "package a\n",
		"b/c.go":      "package c\n",
		"d/e/f/g.txt": "x\n",
	})
	snap, err := Scan(root, nil)
	if err != nil {
		t.Fatal(err)
	}

	first, count := snap.Fingerprint()
	for range 5 {
		again, againCount := snap.Fingerprint()
		if again != first || againCount != count {
			t.Fatalf("Fingerprint is not stable: %q/%d then %q/%d", first, count, again, againCount)
		}
	}
}

// Map iteration order is random in Go, so an unsorted fingerprint would make a
// watcher re-check on every poll of an unchanged tree.
func TestFingerprintDoesNotDependOnMapIterationOrder(t *testing.T) {
	files := map[string]string{}
	for i := range 40 {
		files[filepath.Join("pkg", string(rune('a'+i%26))+string(rune('0'+i/26)), "f.go")] = "package p\n"
	}
	root := buildFS(t, files)

	snap, err := Scan(root, nil)
	if err != nil {
		t.Fatal(err)
	}

	first, count := snap.Fingerprint()
	for range 20 {
		again, _ := snap.Fingerprint()
		if again != first {
			t.Fatalf("Fingerprint changed between calls on an unchanged tree: %q then %q", first, again)
		}
	}
	if count != len(files) {
		t.Errorf("file count = %d, want %d", count, len(files))
	}
}

func TestFingerprintDetectsAnEditThatKeepsFileSize(t *testing.T) {
	root := buildFS(t, map[string]string{"a.txt": "one\n"})
	snap, _ := Scan(root, nil)
	before, _ := snap.Fingerprint()

	path := filepath.Join(root, "a.txt")
	if err := os.WriteFile(path, []byte("two\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Guarantee a distinct mtime even on a coarse filesystem clock.
	future := snapTime(t, path)
	if err := os.Chtimes(path, future, future); err != nil {
		t.Fatal(err)
	}

	snap2, _ := Scan(root, nil)
	after, _ := snap2.Fingerprint()
	if after == before {
		t.Error("an edit that preserved the file size was not detected; size alone is not enough")
	}
}

func TestFingerprintDetectsARenameAndARemoval(t *testing.T) {
	root := buildFS(t, map[string]string{"a.txt": "x\n"})
	snap, _ := Scan(root, nil)
	before, _ := snap.Fingerprint()

	if err := os.Rename(filepath.Join(root, "a.txt"), filepath.Join(root, "b.txt")); err != nil {
		t.Fatal(err)
	}
	snap2, _ := Scan(root, nil)
	renamed, _ := snap2.Fingerprint()
	if renamed == before {
		t.Error("a rename was not detected")
	}

	if err := os.Remove(filepath.Join(root, "b.txt")); err != nil {
		t.Fatal(err)
	}
	snap3, _ := Scan(root, nil)
	removed, files := snap3.Fingerprint()
	if removed == renamed {
		t.Error("a removal was not detected")
	}
	if files != 0 {
		t.Errorf("file count = %d after removing everything, want 0", files)
	}
}

// An ignored directory must not reach the fingerprint, otherwise every build
// artefact re-triggers a check.
func TestFingerprintRespectsIgnoredPaths(t *testing.T) {
	root := buildFS(t, map[string]string{
		"main.go":        "package main\n",
		"dist/bundle.js": "1\n",
	})
	withDist, _ := Scan(root, nil)
	fpWithDist, _ := withDist.Fingerprint()

	withoutDist, err := Scan(root, []string{"dist/"})
	if err != nil {
		t.Fatal(err)
	}
	fpWithout, files := withoutDist.Fingerprint()

	if fpWithDist == fpWithout {
		t.Error("an ignored directory changed the fingerprint")
	}
	if files != 1 {
		t.Errorf("file count = %d, want 1 with dist/ ignored", files)
	}
}

// A silent empty snapshot for a directory that does not exist would leave a
// watcher polling a path forever.
func TestScanReportsAMissingRoot(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope")
	if _, err := Scan(missing, nil); err == nil {
		t.Error("scanning a missing directory should report an error, not an empty snapshot")
	} else if !os.IsNotExist(err) {
		t.Errorf("err = %v, want a not-exist error so callers can tell it apart", err)
	}
}

// A regular file as root would be indexed as the root entry with IsDir cleared,
// so Scan refuses it. cmdctx.ResolvePath rejects this path first in practice.
func TestScanRejectsAFileAsRoot(t *testing.T) {
	root := buildFS(t, map[string]string{"a.txt": "x\n"})
	_, err := Scan(filepath.Join(root, "a.txt"), nil)
	if err == nil {
		t.Fatal("scanning a regular file should report an error, not a snapshot whose root is a file")
	}
	if !strings.Contains(err.Error(), "not a directory") {
		t.Errorf("err = %v, want a clear 'not a directory' message", err)
	}
}

func TestEntryRecordsModTime(t *testing.T) {
	root := buildFS(t, map[string]string{"a.txt": "x\n"})
	snap, err := Scan(root, nil)
	if err != nil {
		t.Fatal(err)
	}

	e, ok := snap.entries["a.txt"]
	if !ok {
		t.Fatal("a.txt missing from the index")
	}
	if e.ModTime == 0 {
		t.Error("Entry.ModTime was not recorded; a same-size edit would go unnoticed")
	}
	if e.Size != 2 {
		t.Errorf("Size = %d, want 2", e.Size)
	}
}

func TestHashStringIsDeterministic(t *testing.T) {
	if hashString("abc") != hashString("abc") {
		t.Error("hashString is not deterministic")
	}
	if hashString("abc") == hashString("abd") {
		t.Error("hashString collided on a one-character difference")
	}
	if !strings.EqualFold(hashString("abc"), hashString("abc")) {
		t.Error("hashString output is not hex-comparable")
	}
}

// snapTime returns a timestamp safely in the future, so an edit is detected even
// on a filesystem with a coarse modification-time clock.
func snapTime(t *testing.T, path string) time.Time {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}
	return info.ModTime().Add(2 * time.Second)
}
