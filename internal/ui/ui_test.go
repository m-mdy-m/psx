package ui

import (
	"strings"
	"testing"
)

// The parsing below is the part that decides what an answer means. It is split
// out from the terminal handling precisely so it can be tested: under `go test`
// IsInteractive is always false, so a prompt would otherwise never be reached.

func TestParseYesNo(t *testing.T) {
	tests := []struct {
		answer string
		def    bool
		want   bool
	}{
		{"", false, false},
		{"", true, true},
		{"y", false, true},
		{"Y", false, true},
		{"yes", false, true},
		{"YES", false, true},
		{"Yes", true, true},
		{"  yes  ", false, true},
		{"n", true, false},
		{"N", true, false},
		{"no", true, false},
		{"NO", true, false},
		{"n", false, false},
		{"maybe", true, true},
		{"maybe", false, false},
		{"1", false, true},
		{"0", true, false},
	}
	for _, tc := range tests {
		if got := parseYesNo(tc.answer, tc.def); got != tc.want {
			t.Errorf("parseYesNo(%q, %v) = %v, want %v", tc.answer, tc.def, got, tc.want)
		}
	}
}

func TestParseChoice(t *testing.T) {
	tests := []struct {
		answer  string
		choices int
		want    int
	}{
		{"1", 3, 0},
		{"2", 3, 1},
		{"3", 3, 2},
		{" 2 ", 3, 1},
		{"", 3, 0},
		{"0", 3, 0},
		{"4", 3, 0},
		{"-1", 3, 0},
		{"abc", 3, 0},
		{"", 0, -1},
		{"1", 0, -1},
	}
	for _, tc := range tests {
		if got := parseChoice(tc.answer, tc.choices); got != tc.want {
			t.Errorf("parseChoice(%q, %d) = %d, want %d", tc.answer, tc.choices, got, tc.want)
		}
	}
}

// A trailing number must not be silently truncated: typing "2x" is a typo, not a
// request for choice 2.
func TestParseChoiceRejectsTrailingJunk(t *testing.T) {
	for _, answer := range []string{"2x", "1abc", "2 3", "1."} {
		if got := parseChoice(answer, 3); got != 0 {
			t.Errorf("parseChoice(%q, 3) = %d, want 0 (the default)", answer, got)
		}
	}
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		answer string
		def    string
		want   string
	}{
		{"notes", "", "notes"},
		{"", "demo", "demo"},
		{"   ", "demo", "demo"},
		{"  notes  ", "", "notes"},
		{"\tnotes\t", "", "notes"},
		{"  ", "", ""},
	}
	for _, tc := range tests {
		if got := parseInput(tc.answer, tc.def); got != tc.want {
			t.Errorf("parseInput(%q, %q) = %q, want %q", tc.answer, tc.def, got, tc.want)
		}
	}
}

func TestParseInputTrimsSurroundingSpaceOnly(t *testing.T) {
	// Interior spacing is part of the answer, not noise.
	if got := parseInput("a  b", ""); got != "a  b" {
		t.Errorf("parseInput trimmed interior space: %q", got)
	}
}

func TestChooseMapsAnIndexBackToItsLabel(t *testing.T) {
	choices := []string{"alpha", "beta", "gamma"}
	for i, want := range choices {
		if got := choose(choices, i); got != want {
			t.Errorf("choose(%v, %d) = %q, want %q", choices, i, got, want)
		}
	}
	if got := choose(choices, -1); got != "" {
		t.Errorf("choose(-1) = %q, want the empty string", got)
	}
	if got := choose(choices, len(choices)); got != "" {
		t.Errorf("choose(%d) = %q, want the empty string", len(choices), got)
	}
	if got := choose(nil, 0); got != "" {
		t.Errorf("choose(nil, 0) = %q, want the empty string", got)
	}
}

// A closed stdin must not block: every prompt returns its default instead.
func TestPromptsReturnDefaultsWithoutATerminal(t *testing.T) {
	t.Setenv("PSX_NON_INTERACTIVE", "1")
	if IsInteractive() {
		t.Fatal("PSX_NON_INTERACTIVE should force the non-interactive path")
	}

	if !Confirm("proceed?", true) {
		t.Error("Confirm should return its default of true when nobody can answer")
	}
	if Confirm("proceed?", false) {
		t.Error("Confirm should return its default of false when nobody can answer")
	}
	if got := Input("name", "demo"); got != "demo" {
		t.Errorf("Input = %q, want the default %q", got, "demo")
	}
	if got := Prompt("pick", []string{"a", "b"}); got != 0 {
		t.Errorf("Prompt = %d, want the first choice", got)
	}
	if got := Choose("pick", []string{"a", "b"}); got != "a" {
		t.Errorf("Choose = %q, want the first choice", got)
	}
}

func TestPromptAndChooseRejectAnEmptyChoiceList(t *testing.T) {
	if got := Prompt("pick", nil); got != -1 {
		t.Errorf("Prompt with no choices = %d, want -1", got)
	}
	if got := Choose("pick", nil); got != "" {
		t.Errorf("Choose with no choices = %q, want the empty string", got)
	}
}

func TestIsInteractiveIsFalseUnderCI(t *testing.T) {
	t.Setenv("PSX_NON_INTERACTIVE", "")
	t.Setenv("CI", "true")
	if IsInteractive() {
		t.Error("CI should force the non-interactive path")
	}
}

func TestIsTerminalRejectsAPipe(t *testing.T) {
	t.Parallel()
	// A regular file is not a character device.
	path := t.TempDir() + "/plain.txt"
	if err := writeFileHelper(path, "x"); err != nil {
		t.Fatal(err)
	}
	f, err := openFileHelper(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if IsTerminal(f) {
		t.Error("a regular file was reported as a terminal")
	}
}

func TestReadLineHandlesEOF(t *testing.T) {
	Stdin = newStringReader("")
	Reset()
	line, ok := readLine()
	if ok {
		t.Errorf("readLine on empty input returned ok with %q", line)
	}
	if line != "" {
		t.Errorf("readLine on empty input = %q, want empty", line)
	}
}

func TestReadLineHandlesAPartialFinalLine(t *testing.T) {
	Stdin = newStringReader("no trailing newline")
	Reset()
	line, ok := readLine()
	if !ok {
		t.Error("a final line without a newline should still be read")
	}
	if line != "no trailing newline" {
		t.Errorf("readLine = %q, want the whole line", line)
	}
}

func TestReadLineReadsOneLineAtATime(t *testing.T) {
	Stdin = newStringReader("first\nsecond\n")
	Reset()

	if line, _ := readLine(); line != "first" {
		t.Errorf("first readLine = %q, want first", line)
	}
	if line, _ := readLine(); line != "second" {
		t.Errorf("second readLine = %q, want second", line)
	}
}

// Reset must make the next read pick up the replaced Stdin, or a test that swaps
// input would silently keep reading the old one.
func TestResetPicksUpAReplacedStdin(t *testing.T) {
	Stdin = newStringReader("old\n")
	Reset()
	if _, _ = readLine(); true {
		// consume the old reader
	}

	Stdin = newStringReader("new\n")
	Reset()
	if line, _ := readLine(); line != "new" {
		t.Errorf("after Reset, readLine = %q, want new", line)
	}
}

func TestReadLineIsSafeOnWhitespaceOnlyInput(t *testing.T) {
	Stdin = newStringReader("   \n")
	Reset()
	line, ok := readLine()
	if !ok {
		t.Error("a blank line is still a line")
	}
	if strings.TrimSpace(line) != "" {
		t.Errorf("readLine = %q, want it trimmed to empty", line)
	}
}
