package command

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/m-mdy-m/psx/internal/flags"
	"github.com/m-mdy-m/psx/internal/resources"
	"github.com/m-mdy-m/psx/internal/utils"
)

var rootCmd = &cobra.Command{
	Use:   "psx",
	Short: "PSX - Project Structure Checker",
	Long: `PSX validates and standardizes project structures across different
programming languages and frameworks.

Examples:
  psx check                  # Validate current project
  psx fix --interactive      # Fix issues with confirmation
  psx project show           # Show cached project info

Documentation: https://github.com/m-mdy-m/psx`,
	Version:          "", // Set in Exec()
	SilenceUsage:     true,
	SilenceErrors:    true,
	PersistentPreRun: preRun,
}

func Execute(version string) error {
	rootCmd.Version = version
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "--help", "-h", "help":
			resources.GetMessage("help", "main")
			return nil
		case "--version", "-v", "version":
			utils.Version(version)
			return nil
		}
	}

	return rootCmd.Execute()
}

func init() {
	initGlobalFlags()
	rootCmd.AddCommand(CheckCmd)
	rootCmd.AddCommand(FixCmd)
}

func initGlobalFlags() {
	f := flags.GetFlags()
	df := flags.DefaultValues.GlobalFlags

	rootCmd.PersistentFlags().StringVar(&f.GlobalFlags.ConfigFile, "config", df.ConfigFile,
		"config file path")

	rootCmd.PersistentFlags().BoolVarP(&f.GlobalFlags.Verbose, "verbose", "v", df.Verbose,
		"verbose output")

	rootCmd.PersistentFlags().BoolVarP(&f.GlobalFlags.Quiet, "quiet", "q", df.Quiet,
		"quiet mode")

	rootCmd.PersistentFlags().BoolVar(&f.GlobalFlags.NoColor, "no-color", df.NoColor,
		"disable colors")
}

func preRun(cmd *cobra.Command, args []string) {
	f := flags.GetFlags()
	if f.GlobalFlags.NoColor || !isTerminal() {
		color.NoColor = true
	}
	if os.Getenv("CI") != "" {
		f.GlobalFlags.NoColor = true
		color.NoColor = true
	}
}

func isTerminal() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
