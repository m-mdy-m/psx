// Package command wires the cobra command tree.
//
// Commands stay thin: each parses flags into a flags.Options value, delegates to
// a domain package, and returns an error. Shared setup lives in helpers here
// rather than being repeated per command.
package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/ui"
)

// Version is set at build time.
var Version = "dev"

// Process exit codes.
const (
	exitCode   = 1 // validation failed
	exitConfig = 2 // configuration or usage problem
	exitArgs   = 3 // invalid arguments
)

// globalOpts accumulates flags shared by every command.
var globalOpts struct {
	configFile string
	verbose    bool
	quiet      bool
	noColor    bool
	yes        bool
}

// NewRootCommand builds the command tree.
//
// Commands are registered from a table so adding one is a single entry rather
// than a new init function and an AddCommand call.
func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "psx",
		Short: "Validate and standardise project structure",
		Long: "psx checks that a project follows a consistent structure, and can create the\n" +
			"missing files for you.\n\n" +
			"Rules are declarative: inspect them with `psx rules list`.",
		Version:       Version,
		SilenceUsage:  true,
		SilenceErrors: true,
		// A context is threaded through so long-running commands can be
		// cancelled with Ctrl+C without touching global state.
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			configureLogging()
		},
	}

	pf := root.PersistentFlags()
	pf.StringVar(&globalOpts.configFile, "config", "", "path to a psx configuration file")
	pf.BoolVarP(&globalOpts.verbose, "verbose", "v", false, "show detailed output")
	pf.BoolVarP(&globalOpts.quiet, "quiet", "q", false, "suppress informational output")
	pf.BoolVar(&globalOpts.noColor, "no-color", false, "disable coloured output")
	pf.BoolVarP(&globalOpts.yes, "yes", "y", false, "assume yes for prompts and never prompt at all")

	for _, c := range commands() {
		root.AddCommand(c)
	}
	return root
}

// The table is the single registration point: adding a command means adding one
// entry, with no separate init function to keep in sync.
func commands() []*cobra.Command {
	return []*cobra.Command{
		newCheckCmd(),
		newFixCmd(),
		newWatchCmd(),
		newInitCmd(),
		newRulesCmd(),
		newWorkflowsCmd(),
		newDetectCmd(),
		newExplainCmd(),
	}
}

// configureLogging wires the logger to the resolved flags and environment.
func configureLogging() {
	noColor := globalOpts.noColor || !ui.IsTerminal(os.Stdout) || os.Getenv("CI") != ""
	logger.Configure(globalOpts.verbose, globalOpts.quiet, noColor)
}

// Execute runs the CLI and returns the process exit code.
//
// Errors are reported here rather than by cobra, so an unhandled failure is
// never silent: the previous version silenced cobra's error output and printed
// nothing at all.
func Execute() int {
	root := NewRootCommand()

	if err := root.Execute(); err != nil {
		var ec *exitCodedError
		if ok := asExitCoded(err, &ec); ok {
			logger.Error(ec.message)
			return ec.code
		}
		logger.Error(err.Error())
		return exitCode
	}
	return 0
}

type exitCodedError struct {
	code    int
	message string
}

func (e *exitCodedError) Error() string { return e.message }

// asExitCoded unwraps an error into an exitCodedError when possible.
func asExitCoded(err error, target **exitCodedError) bool {
	for err != nil {
		if ec, ok := err.(*exitCodedError); ok {
			*target = ec
			return true
		}
		u, ok := err.(interface{ Unwrap() error })
		if !ok {
			return false
		}
		err = u.Unwrap()
	}
	return false
}

func commandError(code int, format string, args ...any) error {
	return &exitCodedError{code: code, message: fmt.Sprintf(format, args...)}
}
