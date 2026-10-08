// Package ui handles terminal interaction.
//
// Everything that reads stdin or writes a progress indicator lives here, so the
// rest of the codebase can stay pure and testable. When stdin is not a
// terminal, prompts return their default instead of blocking.
package ui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Stdin and Stdout are overridable for tests.
var (
	Stdin  io.Reader = os.Stdin
	Stdout io.Writer = os.Stdout

	reader *bufio.Reader
)

// initReader lazily wraps stdin so a test can swap Stdin first.
func initReader() {
	if reader == nil {
		reader = bufio.NewReader(Stdin)
	}
}

// Reset drops the buffered reader, for tests that replace Stdin.
func Reset() { reader = nil }

// IsInteractive reports whether a human is present to answer questions.
func IsInteractive() bool {
	if os.Getenv("CI") != "" {
		return false
	}
	if os.Getenv("PSX_NON_INTERACTIVE") != "" {
		return false
	}
	return IsTerminal(os.Stdin) && IsTerminal(os.Stdout)
}

func IsTerminal(f *os.File) bool {
	info, err := f.Stat()
	if err != nil {
		// Without a stat result there is no basis for claiming a terminal,
		// and the previous behaviour dereferenced a nil FileInfo here.
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

func Confirm(question string, def bool) bool {
	if !IsInteractive() {
		return def
	}

	suffix := "[y/N]"
	if def {
		suffix = "[Y/n]"
	}
	fmt.Fprintf(Stdout, "%s %s: ", question, suffix)

	answer, ok := readLine()
	if !ok {
		return def
	}
	return parseYesNo(answer, def)
}

// parseYesNo reads a yes/no answer, falling back to def for anything it does not
// recognise. It is separated from Confirm so the decision can be tested without
// a terminal.
func parseYesNo(answer string, def bool) bool {
	switch strings.ToLower(strings.TrimSpace(answer)) {
	case "y", "yes", "1":
		return true
	case "n", "no", "0":
		return false
	default:
		return def
	}
}

type ConfirmAll int

const (
	// ConfirmNo skips just this item.
	ConfirmNo ConfirmAll = iota
	// ConfirmYes applies to this item only.
	ConfirmYes
	// ConfirmAllYes applies to this and every remaining item.
	ConfirmAllYes
	// ConfirmAllNo skips this and every remaining item.
	ConfirmAllNo
	// ConfirmQuit stops processing entirely.
	ConfirmQuit
)

func Prompt(question string, choices []string) int {
	if len(choices) == 0 {
		return -1
	}
	if !IsInteractive() {
		return 0
	}

	fmt.Fprintln(Stdout, question)
	for i, c := range choices {
		fmt.Fprintf(Stdout, "  %d) %s\n", i+1, c)
	}
	fmt.Fprintf(Stdout, "  choice [1]: ")

	answer, ok := readLine()
	if !ok {
		return 0
	}
	return parseChoice(answer, len(choices))
}

// parseChoice maps a 1-based answer onto a 0-based index, defaulting to the
// first choice. It rejects trailing junk so "2x" is a typo rather than a choice.
func parseChoice(answer string, choices int) int {
	if choices <= 0 {
		return -1
	}
	answer = strings.TrimSpace(answer)
	if answer == "" {
		return 0
	}
	for _, r := range answer {
		if r < '0' || r > '9' {
			return 0
		}
	}
	n, err := strconv.Atoi(answer)
	if err != nil || n < 1 || n > choices {
		return 0
	}
	return n - 1
}

func Choose(question string, choices []string) string {
	return choose(choices, Prompt(question, choices))
}

func choose(choices []string, i int) string {
	if i < 0 || i >= len(choices) {
		return ""
	}
	return choices[i]
}

func Input(question, def string) string {
	if !IsInteractive() {
		return def
	}
	if def != "" {
		fmt.Fprintf(Stdout, "%s [%s]: ", question, def)
	} else {
		fmt.Fprintf(Stdout, "%s: ", question)
	}

	answer, ok := readLine()
	if !ok {
		return def
	}
	return parseInput(answer, def)
}

// parseInput trims the answer and substitutes def for a blank one.
func parseInput(answer, def string) string {
	answer = strings.TrimSpace(answer)
	if answer == "" {
		return def
	}
	return answer
}

func readLine() (string, bool) {
	initReader()
	line, err := reader.ReadString('\n')
	if err != nil && line == "" {
		// EOF or a read error: fall back to defaults rather than spinning.
		return "", false
	}
	return strings.TrimSpace(line), true
}
