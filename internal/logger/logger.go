package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

// Out and Err are the output streams, overridable in tests.
var (
	Out io.Writer = os.Stdout
	Err io.Writer = os.Stderr
)

var (
	verbose bool
	quiet   bool
	noColor bool
)

// Configure sets the logging behaviour from the global flags.
func Configure(v, q, nc bool) {
	verbose, quiet, noColor = v, q, nc
	color.NoColor = nc
}

func VerboseEnabled() bool { return verbose && !quiet }

func QuietEnabled() bool { return quiet }

func NoColorEnabled() bool { return noColor }

func infoColor() *color.Color { return color.New(color.FgBlue) }
func successColor() *color.Color {
	return color.New(color.FgGreen)
}
func warnColor() *color.Color {
	return color.New(color.FgYellow)
}
func errColor() *color.Color  { return color.New(color.FgRed) }
func verbColor() *color.Color { return color.New(color.FgCyan) }
func dimColor() *color.Color  { return color.New(color.Faint) }

// Info writes an informational line to stderr.
func Info(m string) {
	if quiet {
		return
	}
	fmt.Fprintln(Err, infoColor().Sprint("›"), m)
}

// Infof writes a formatted informational line to stderr.
func Infof(format string, args ...any) { Info(fmt.Sprintf(format, args...)) }

// Success writes a success line to stderr.
func Success(m string) {
	if quiet {
		return
	}
	fmt.Fprintln(Err, successColor().Sprint("✓"), m)
}

// Successf writes a formatted success line to stderr.
func Successf(format string, args ...any) { Success(fmt.Sprintf(format, args...)) }

// Warning writes a warning line to stderr. Warnings survive quiet mode only when
// they indicate a problem, so they are suppressed like other info output.
func Warning(m string) {
	if quiet {
		return
	}
	fmt.Fprintln(Err, warnColor().Sprint("⚠"), m)
}

// Warningf writes a formatted warning line to stderr.
func Warningf(format string, args ...any) { Warning(fmt.Sprintf(format, args...)) }

// Error writes an error line to stderr. Errors are never suppressed by quiet.
func Error(m string) {
	fmt.Fprintln(Err, errColor().Sprint("✗"), m)
}

// Errorf writes a formatted error line to stderr and returns it as an error, so
// call sites can both report and propagate in one statement.
func Errorf(format string, args ...any) error {
	err := fmt.Errorf(format, args...)
	Error(err.Error())
	return err
}

// Verbose writes a diagnostic line to stderr when verbose output is enabled.
func Verbose(m string) {
	if !VerboseEnabled() {
		return
	}
	fmt.Fprintln(Err, verbColor().Sprint("→"), m)
}

// Verbosef writes a formatted diagnostic line to stderr.
func Verbosef(format string, args ...any) { Verbose(fmt.Sprintf(format, args...)) }

// Step writes a progress line to stdout. It is reserved for interactive flows
// where the user is watching, not for the final report.
func Step(m string) {
	if quiet {
		return
	}
	fmt.Fprintln(Out, dimColor().Sprint("›"), m)
}

// Plain writes text to stdout with no decoration.
func Plain(s string) { fmt.Fprintln(Out, s) }

// Fatal writes an error line and exits with a failure code.
func Fatal(m string) {
	Error(m)
	os.Exit(1)
}

// Fatalf writes a formatted error line and exits with a failure code.
func Fatalf(format string, args ...any) {
	Fatal(fmt.Sprintf(format, args...))
}
