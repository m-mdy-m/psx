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
	"github.com/m-mdy-m/psx/internal/shared"
)

var CheckCmd = &cobra.Command{
	Use:   "check [path]",
	Short: "Validate project structure",
	Long: `Validate project structure against configured rules.

PSX checks for:
- Essential files (README, LICENSE, .gitignore)
- Project structure (src/, tests/, docs/ folders)
- Documentation requirements
- CI/CD configuration
- Code quality tools

Examples:
  psx check                       # Check current directory
  psx check ./my-project          # Check specific directory
  psx check --verbose             # Show detailed information
  psx check --output json         # JSON output for CI/CD
  psx check --level error         # Only show errors
  psx check --fail-on warning     # Fail on warnings too`,
	Args: cobra.MaximumNArgs(1),
	RunE: runCheckCommand,
}

func init() {
	f := flags.GetFlags()
	df := flags.DefaultValues.Check

	CheckCmd.Flags().StringVarP(&f.Check.OutputFormat, "output", "o", df.OutputFormat,
		"output format: table | json")

	CheckCmd.Flags().StringVar(&f.Check.ServerityLevel, "level", df.ServerityLevel,
		"filter by severity: error | warning | info | all")

	CheckCmd.Flags().StringVar(&f.Check.FailOn, "fail-on", df.FailOn,
		"exit with error on: error | warning")
}

func runCheckCommand(cmd *cobra.Command, args []string) error {
	// Load project context
	ctx, err := cmdctx.LoadProject(args)
	if err != nil {
		return err
	}

	f := flags.GetFlags()

	// Verbose output
	logger.Verbose(shared.VerboseCheckStart(ctx.Path.Abs))
	logger.Verbose(shared.VerboseDetected(ctx.Detection.Type.Primary))
	logger.Verbose(shared.VerboseRulesLoaded(len(ctx.Config.ActiveRules)))
	logger.Verbose(shared.VerboseConfigLoaded(ctx.Config.Path))

	// Execute rules
	engine := rules.NewEngine(ctx.Config, ctx.Detection)
	result, err := engine.Execute()
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Generate report
	rep := reporter.New(f.Check.OutputFormat, result)
	if err := rep.Report(); err != nil {
		return fmt.Errorf("report generation failed: %w", err)
	}

	// Determine exit code
	return determineExitCode(result, f.Check.FailOn)
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
	default:
		shouldFail = hasErrors
	}

	if shouldFail {
		if !f.GlobalFlags.Quiet && result.Summary.Errors > 0 {
			fmt.Println()
			logger.Info("Run 'psx fix' to fix issues automatically")
		}
		os.Exit(shared.ExitFailed)
	}

	return nil
}

