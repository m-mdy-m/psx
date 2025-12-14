package logger

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/m-mdy-m/psx/internal/flags"
)

func Info(m string) {
	f := flags.GetFlags()
	if f.GlobalFlags.IsQuiet() {
		return
	}
	fmt.Println(m)
}

func Infof(format string, args ...any) {
	Info(fmt.Sprintf(format, args...))
}

func Success(m string) {
	f := flags.GetFlags()
	if f.GlobalFlags.IsQuiet() {
		return
	}

	if f.GlobalFlags.NoColor {
		fmt.Println("✓", m)
	} else {
		color.Green("✓ %s", m)
	}
}

func Successf(format string, args ...any) {
	Success(fmt.Sprintf(format, args...))
}

func Warning(m string) {
	f := flags.GetFlags()
	if f.GlobalFlags.IsQuiet() {
		return
	}

	if f.GlobalFlags.NoColor {
		fmt.Println("⚠", m)
	} else {
		color.Yellow("⚠ %s", m)
	}
}

func Warningf(format string, args ...any) {
	Warning(fmt.Sprintf(format, args...))
}

func Error(m string) {
	f := flags.GetFlags()

	if f.GlobalFlags.NoColor {
		fmt.Fprintln(os.Stderr, "✗", m)
	} else {
		color.New(color.FgRed).Fprintln(os.Stderr, "✗", m)
	}
}

func Errorf(format string, args ...any) {
	Error(fmt.Sprintf(format, args...))
}

// Verbose prints debug info (only with --verbose)
func Verbose(m string) {
	f := flags.GetFlags()
	if !f.GlobalFlags.Verbose || f.GlobalFlags.IsQuiet() {
		return
	}

	if f.GlobalFlags.NoColor {
		fmt.Println("→", m)
	} else {
		color.Cyan("→ %s", m)
	}
}

func Verbosef(format string, args ...any) {
	Verbose(fmt.Sprintf(format, args...))
}

func Fatal(m string) {
	Error(m)
	os.Exit(1)
}

func Fatalf(format string, args ...any) {
	Fatal(fmt.Sprintf(format, args...))
}

