package commond

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/m-mdy-m/psx/internal/flags"
	"github.com/m-mdy-m/psx/internal/resources"
	"github.com/m-mdy-m/psx/internal/shared"
)

var root = &cobra.Command{
	Use:   "psx",
	Short: "PSX - Project Structure Checker",
	Long: `PSX validates and standardizes project structures across different
programming languages and frameworks.

Examples:
  psx check                  # Validate current project
  psx fix --interactive      # Fix issues with confirmation
  psx project show           # Show cached project info

Documentation: https://github.com/m-mdy-m/psx`,
	Version:           "", // Set in Exec()
	SilenceUsage:      true,
	SilenceErrors:     true,
	PersistentPreRun:  preRun,
	PersistentPostRun: postRun,
}

func Exec(version string) error {
	root.Version = version

	// Handle help and version specially
	if len(os.Args) == 1 {
		resources.HelpMain()
		return nil
	}

	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "--help", "-h", "help":
			resources.HelpMain()
			return nil
		case "--version", "-v", "version":
			shared.Version(version)
			return nil
		}
	}

	return root.Execute()
}

func preRun(cmd *cobra.Command, args []string) {
	f := flags.GetFlags()

	// Disable colors if requested or not a TTY
	if f.GlobalFlags.NoColor || !isTerminal() {
		color.NoColor = true
	}

	// CI detection
	if os.Getenv("CI") != "" {
		f.GlobalFlags.NoColor = true
		color.NoColor = true
	}
}

func postRun(cmd *cobra.Command, args []string) {
	// Cleanup if needed
}

func initGlobalFlags() {
	f := flags.GetFlags()
	df := flags.DefaultValues.GlobalFlags

	root.PersistentFlags().StringVar(&f.GlobalFlags.ConfigFile, "config", df.ConfigFile,
		"config file path")

	root.PersistentFlags().BoolVarP(&f.GlobalFlags.Verbose, "verbose", "v", df.Verbose,
		"verbose output")

	root.PersistentFlags().BoolVarP(&f.GlobalFlags.Quiet, "quiet", "q", df.Quiet,
		"quiet mode")

	root.PersistentFlags().BoolVar(&f.GlobalFlags.NoColor, "no-color", df.NoColor,
		"disable colors")
}

func isTerminal() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

func init() {
	initGlobalFlags()

	// Add all commands
	root.AddCommand(CheckCmd)
	root.AddCommand(FixCmd)
	root.AddCommand(ProjectCmd)
}
