package flags

import (
	"fmt"
	"time"
)

// GlobalFlags apply to every command.
type GlobalFlags struct {
	ConfigFile string
	Verbose    bool
	Quiet      bool
	NoColor    bool
	Yes        bool
}

// Report output formats.
const (
	FormatTable    = "table"
	FormatCompact  = "compact"
	FormatJSON     = "json"
	FormatNDJSON   = "ndjson"
	FormatSARIF    = "sarif"
	FormatGitHub   = "github"
	FormatJUnit    = "junit"
	FormatMarkdown = "markdown"
)

// Check holds options for the check command.
type Check struct {
	OutputFormat string
	Level        string
	FailOn       string
	Only         []string
	Category     []string
	ShowPassed   bool
	ShowSkipped  bool
	BaselineFile string
}

// Fix holds options for the fix command.
type Fix struct {
	RuleID        string
	DryRun        bool
	Force         bool
	Interactive   bool
	CreateBackups bool
	Answers       map[string]string
	Watch         bool
}

// Prompting requires an explicit request and a real terminal, so a run in CI or
// a piped shell can never block on stdin.
func (f Fix) InteractiveWanted() bool { return f.Interactive }

// Init holds options for the init command.
type Init struct {
	ProjectType string
	Kind        string
	License     string
	Force       bool
	Minimal     bool
}

// Watch holds options for the watch command.
type Watch struct {
	Interval time.Duration
	Debounce time.Duration
	Fix      bool
	Once     bool
	Format   string
}

// Options is the full parsed option set for a single invocation.
type Options struct {
	Global GlobalFlags
	Check  Check
	Fix    Fix
	Init   Init
	Watch  Watch

	// Path is the project directory to operate on.
	Path string
	// Version is the build version, set by the entry point.
	Version string
}

// Defaults returns the default option set.
func Defaults() Options {
	return Options{
		Check: Check{
			OutputFormat: FormatTable,
			Level:        "all",
			FailOn:       "error",
		},
		Fix: Fix{
			Interactive: false,
			Answers:     map[string]string{},
		},
		Watch: Watch{
			Interval: 2 * time.Second,
			Debounce: 300 * time.Millisecond,
			Format:   FormatNDJSON,
		},
	}
}

// Validate checks option consistency before any work starts.
func (o *Options) Validate() error {
	switch o.Check.OutputFormat {
	case FormatTable, FormatCompact, FormatJSON, FormatNDJSON,
		FormatSARIF, FormatGitHub, FormatJUnit, FormatMarkdown:
	default:
		return errf("unsupported output format %q", o.Check.OutputFormat)
	}

	switch o.Check.Level {
	case "all", "error", "warning", "info":
	default:
		return errf("unsupported severity level %q", o.Check.Level)
	}

	switch o.Check.FailOn {
	case "error", "warning", "none":
	default:
		return errf("unsupported --fail-on value %q", o.Check.FailOn)
	}

	if o.Fix.DryRun && o.Fix.Force {
		return errf("--dry-run and --force cannot be combined")
	}
	if o.Fix.Watch && !o.Fix.DryRun {
		// Watch owns its own loop; a single-shot fix must not also watch.
		o.Fix.Watch = false
	}
	return nil
}

func (o *Options) Verbose() bool { return o.Global.Verbose && !o.Global.Quiet }

func (o *Options) machineFormat() bool {
	switch o.Check.OutputFormat {
	case FormatJSON, FormatNDJSON, FormatSARIF, FormatJUnit:
		return true
	default:
		return false
	}
}

// errf builds a validation error.
func errf(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}
