package commond

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"

	"github.com/m-mdy-m/psx/internal/cmdctx"
	"github.com/m-mdy-m/psx/internal/flags"
	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/reporter"
	"github.com/m-mdy-m/psx/internal/rules"
)

var CheckCmd = &cobra.Command{
	Use:   "check [path]",
	Short: "Validate project structure",
	Long: `Validates your project structure against configured rules.

Examples:
  psx check                    # Check current directory
  psx check ./my-project       # Check specific directory
  psx check --verbose          # Show detailed information
  psx check --output json      # JSON output
  psx check --level error      # Only show errors
  psx check --fail-on warning  # Exit with error on warnings`,
	Args: cobra.MaximumNArgs(1),
	RunE: runCheckCommand,
}

func init() {
	f := flags.GetFlags()
	dcf := flags.DefaultValues.Check

	CheckCmd.Flags().StringVarP(&f.Check.OutputFormat, "output", "o", dcf.OutputFormat,
		"output format: table | json")

	CheckCmd.Flags().StringVar(&f.Check.ServerityLevel, "level", dcf.ServerityLevel,
		"filter by severity: error | warning | info | all")

	CheckCmd.Flags().StringVar(&f.Check.FailOn, "fail-on", dcf.FailOn,
		"exit with error when: error | warning")
}

func runCheckCommand(cmd *cobra.Command, args []string) error {
	ctx, err := cmdctx.LoadProject(args)
	if err != nil {
		return err
	}

	engine := rules.NewEngine(ctx.Config,ctx.Detection)
	result, err := engine.Execute()
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	logger.Verbosef("Checked %d rules", result.Summary.Total)

	// Generate report
	rep := reporter.New(ctx.Flags.Check.OutputFormat, result)
	if err := rep.Report(); err != nil {
		return fmt.Errorf("report failed: %w", err)
	}

	// Determine exit code
	return determineExitCode(result, ctx.Flags.Check.FailOn)
}

func determineExitCode(result *rules.ExecutionResult, failOn string) error {
	f := flags.GetFlags()

	hasErrors := result.Summary.Errors > 0
	hasWarnings := result.Summary.Warnings > 0

	shouldFail := false

	switch failOn {
		case "error":
			shouldFail = hasErrors
		case "warning":
			shouldFail = hasErrors || hasWarnings
	}

	if shouldFail {
		if !f.GlobalFlags.Quiet {
			fmt.Println()
			logger.Info("Run 'psx fix' to fix issues automatically")
		}
		os.Exit(1)
	}

	return nil
}

