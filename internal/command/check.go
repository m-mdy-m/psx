package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/m-mdy-m/psx/internal/cmdctx"
	"github.com/m-mdy-m/psx/internal/flags"
	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/reporter"
	"github.com/m-mdy-m/psx/internal/resources"
	"github.com/m-mdy-m/psx/internal/rules"
	"github.com/m-mdy-m/psx/internal/utils"
)

var CheckCmd = &cobra.Command{
	Use:   "check [path]",
	Short: "Validate project structure",
	Long: `Validate project structure against configured rules.

Examples:
  psx check                       # Check current directory
  psx check ./my-project          # Check specific directory
  psx check --verbose             # Show detailed information
  psx check --output json         # JSON output for CI/CD`,
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
	ctx, err := cmdctx.LoadProject(args)
	if err != nil {
		return err
	}

	f := flags.GetFlags()
	logger.Verbose(resources.FormatMessage("check", "start", ctx.Path.Abs))
	logger.Verbose(fmt.Sprintf("Project type: %s", ctx.ProjectType))
	logger.Verbose(fmt.Sprintf("Active rules: %d", len(ctx.Config.ActiveRules)))
	if ctx.Config.Path != "" {
		logger.Verbose(fmt.Sprintf("Config: %s", ctx.Config.Path))
	}
	rulesCtx := &rules.Context{
		ProjectPath: ctx.Path.Abs,
		ProjectType: ctx.ProjectType,
		ProjectInfo: ctx.ProjectInfo,
		Config:      ctx.Config,
	}
	result, err := rules.Execute(ctx.Config, rulesCtx)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	rep := reporter.New(f.Check.OutputFormat, result)
	if err := rep.Report(); err != nil {
		return fmt.Errorf("report generation failed: %w", err)
	}
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
			logger.Info(resources.GetMessage("fix", "suggest"))
		}
		os.Exit(utils.ExitFailed)
	}

	return nil
}
