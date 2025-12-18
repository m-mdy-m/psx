package utils

import (
	"fmt"

	"github.com/fatih/color"
)

// Exit codes
const (
	ExitSuccess = 0
	ExitFailed  = 1
	ExitConfig  = 2
	ExitFS      = 3
	ExitArgs    = 4
)
const (
	SeverityErrorIcon   = "✗"
	SeverityWarningIcon = "⚠"
	SeverityInfoIcon    = "ℹ"
)

func SeverityColor(level string) *color.Color {
	switch level {
	case "error":
		return color.New(color.FgRed)
	case "warn":
		return color.New(color.FgYellow)
	default:
		return color.New(color.FgCyan)
	}
}

func Version(version string) {
	fmt.Printf("PSX version %s\n", version)
}
