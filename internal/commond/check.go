package commond

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/detector"
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
	// Get project path
	root := "."
	if len(args) > 0 {
		root = args[0]
	}

	abs, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	f := flags.GetFlags()

	// Validate path
	if _, err := os.Stat(abs); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", abs)
	}

	logger.Verbose("Analyzing project...")
	logger.Verbosef("Path: %s", abs)

	// Load config
	logger.Verbose("Loading configuration...")
	cfg, err := config.Load(f.GlobalFlags.ConfigFile, abs)
	if err != nil {
		return fmt.Errorf("config load failed: %w", err)
	}
	logger.Verbosef("Loaded %d rules", len(cfg.ActiveRules))

	// Detect project type
	logger.Verbose("Detecting project type...")
	var detectionResult *detector.DetectionResult

	if cfg.Project.Type != "" {
		logger.Verbosef("Using configured type: %s", cfg.Project.Type)
		detectionResult, err = detector.DetectWithHint(abs, cfg.Project.Type)
		if err != nil {
			logger.Warning("Could not verify type, using auto-detection")
			detectionResult, err = detector.Detect(abs)
			if err != nil {
				return fmt.Errorf("detection failed: %w", err)
			}
		}
	} else {
		detectionResult, err = detector.Detect(abs)
		if err != nil {
			return fmt.Errorf("detection failed: %w", err)
		}
	}

	cfg.Project.Type = detectionResult.Type.Primary
	logger.Verbosef("Detected: %s", detectionResult.Type.Primary)

	// Run validation
	logger.Verbose("Running validation...")
	engine := rules.NewEngine(cfg, detectionResult)
	execResult, err := engine.Execute()
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	logger.Verbosef("Checked %d rules", execResult.Summary.Total)

	// Generate report
	rep := reporter.New(f.Check.OutputFormat, execResult)
	if err := rep.Report(); err != nil {
		return fmt.Errorf("report failed: %w", err)
	}

	// Determine exit code
	return determineExitCode(execResult, f.Check.FailOn)
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

