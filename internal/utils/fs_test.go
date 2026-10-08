package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileExists(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "a.txt")
	if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	ok, info := FileExists(path)
	if !ok {
		t.Fatal("FileExists said a file it just wrote does not exist")
	}
	if info == nil {
		t.Fatal("FileExists returned exists=true with a nil FileInfo")
	}
	if info.Size() != 1 {
		t.Errorf("Size = %d, want 1", info.Size())
	}
}

// The second return must be nil when the path is absent, or a caller that
// ignores the bool dereferences nothing.
func TestFileExistsReturnsNoInfoForAMissingPath(t *testing.T) {
	ok, info := FileExists(filepath.Join(t.TempDir(), "nope"))
	if ok {
		t.Error("FileExists reported a missing path as existing")
	}
	if info != nil {
		t.Errorf("FileExists returned info %v for a missing path, want nil", info)
	}
}

func TestFileExistsOnADirectory(t *testing.T) {
	ok, info := FileExists(t.TempDir())
	if !ok {
		t.Error("a directory exists")
	}
	if info == nil || !info.IsDir() {
		t.Error("FileExists should report a directory as a directory")
	}
}

func TestCreateFileMakesItsParentDirectories(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a", "b", "c.txt")
	if err := CreateFile(path, "hello"); err != nil {
		t.Fatalf("CreateFile: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("the file was not created: %v", err)
	}
	if string(got) != "hello" {
		t.Errorf("content = %q, want hello", got)
	}
}

func TestCreateFileOverwrites(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a.txt")
	if err := CreateFile(path, "first"); err != nil {
		t.Fatal(err)
	}
	if err := CreateFile(path, "second"); err != nil {
		t.Fatalf("CreateFile on an existing file: %v", err)
	}

	got, _ := os.ReadFile(path)
	if string(got) != "second" {
		t.Errorf("content = %q, want second", got)
	}
}

func TestCreateFileWithEmptyContent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "empty.txt")
	if err := CreateFile(path, ""); err != nil {
		t.Fatalf("CreateFile with no content: %v", err)
	}
	ok, info := FileExists(path)
	if !ok {
		t.Fatal("an empty file should still exist")
	}
	if info.Size() != 0 {
		t.Errorf("Size = %d, want 0", info.Size())
	}
}

func TestCreateFileReportsAFailedWrite(t *testing.T) {
	// A path whose parent is a regular file cannot be created.
	root := t.TempDir()
	blocker := filepath.Join(root, "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := CreateFile(filepath.Join(blocker, "child.txt"), "x")
	if err == nil {
		t.Fatal("CreateFile should fail when the parent is a file")
	}
	if !strings.Contains(err.Error(), "failed to") {
		t.Errorf("err = %q, want a message naming what failed", err)
	}
}

func TestCreateDirIsIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a", "b")
	if err := CreateDir(path); err != nil {
		t.Fatalf("CreateDir: %v", err)
	}
	if err := CreateDir(path); err != nil {
		t.Errorf("CreateDir on an existing directory should succeed: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		t.Errorf("Stat(%s) = %v, %v; want a directory", path, info, err)
	}
}

func TestCreateDirReportsAFailure(t *testing.T) {
	root := t.TempDir()
	blocker := filepath.Join(root, "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := CreateDir(filepath.Join(blocker, "child")); err == nil {
		t.Error("CreateDir should fail when the parent is a file")
	}
}

func TestIsDirEmpty(t *testing.T) {
	root := t.TempDir()
	empty := filepath.Join(root, "empty")
	if err := CreateDir(empty); err != nil {
		t.Fatal(err)
	}

	ok, err := IsDirEmpty(empty)
	if err != nil {
		t.Fatalf("IsDirEmpty: %v", err)
	}
	if !ok {
		t.Error("a freshly created directory is empty")
	}

	full := filepath.Join(root, "full")
	if err := CreateFile(filepath.Join(full, "a.txt"), "x"); err != nil {
		t.Fatal(err)
	}
	ok, err = IsDirEmpty(full)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("a directory holding a file is not empty")
	}
}

func TestIsDirEmptyReportsAMissingPath(t *testing.T) {
	ok, err := IsDirEmpty(filepath.Join(t.TempDir(), "nope"))
	if err == nil {
		t.Error("a missing path should be reported as an error, not as empty")
	}
	if ok {
		t.Error("a missing path is not an empty directory")
	}
}

// IsDirEmpty on a file is a caller mistake, and it must surface rather than
// quietly claiming the file is an empty directory.
func TestIsDirEmptyOnAFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a.txt")
	if err := CreateFile(path, "x"); err != nil {
		t.Fatal(err)
	}
	if _, err := IsDirEmpty(path); err == nil {
		t.Error("passing a file to IsDirEmpty should report an error")
	}
}
