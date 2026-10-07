package command

import (
	"errors"
	"testing"

	"github.com/m-mdy-m/psx/internal/rules"
)

// A folder tree mixes files and directories; counting both as files told the user
// psx created files it had actually created directories.
func TestFixSummaryCountsFilesAndDirectoriesSeparately(t *testing.T) {
	results := []*rules.FixResult{{
		Fixed: true,
		Changes: []rules.Change{
			{Type: rules.ChangeCreateFolder, Path: "/p/config"},
			{Type: rules.ChangeCreateFile, Path: "/p/config/prod.yaml"},
			{Type: rules.ChangeCreateFile, Path: "/p/config/stage.yaml"},
			{Type: rules.ChangeCreateFile, Path: "/p/README.md"},
		},
	}}

	s := summarizeFix(results)
	if s.Files != 3 {
		t.Errorf("Files = %d, want 3", s.Files)
	}
	if s.Dirs != 1 {
		t.Errorf("Dirs = %d, want 1", s.Dirs)
	}
	if got, want := s.describe(), "3 files and 1 directory"; got != want {
		t.Errorf("describe() = %q, want %q", got, want)
	}
}

func TestFixSummaryDescriptionPluralises(t *testing.T) {
	tests := []struct {
		files, dirs int
		want        string
	}{
		{0, 0, "0 files"},
		{1, 0, "1 file"},
		{2, 0, "2 files"},
		{0, 1, "1 directory"},
		{0, 2, "2 directories"},
		{1, 1, "1 file and 1 directory"},
	}
	for _, tc := range tests {
		s := fixSummary{Files: tc.files, Dirs: tc.dirs}
		if got := s.describe(); got != tc.want {
			t.Errorf("describe(%d files, %d dirs) = %q, want %q", tc.files, tc.dirs, got, tc.want)
		}
	}
}

func TestSummarizeFixCountsOutcomesPerRule(t *testing.T) {
	results := []*rules.FixResult{
		{Fixed: true, Changes: []rules.Change{{Type: rules.ChangeCreateFile}}},
		{Skipped: true, Reason: "already present"},
		{Error: errors.New("boom")},
	}
	s := summarizeFix(results)

	if s.Fixed != 1 || s.Skipped != 1 || s.Failed != 1 {
		t.Errorf("got %+v, want Fixed=1 Skipped=1 Failed=1", s)
	}
}
