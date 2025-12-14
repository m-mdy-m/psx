package commond

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/detector"
	"github.com/m-mdy-m/psx/internal/flags"
	"github.com/m-mdy-m/psx/internal/reporter"
	"github.com/m-mdy-m/psx/internal/rules"
	"github.com/m-mdy-m/psx/internal/logger"
)

var CheckCmd = &cobra.Command{
	Use: "check [path]",
	Short: "Check project structure",
	Long: `Validates project structure against configured rules.

Examples:
  psx check                  # Check current directory
  psx check ./my-project     # Check specific directory
  psx check --verbose        # Detailed output
  psx check --output json    # JSON output
  psx check --level error    # Only show errors`,
  Args: cobra.MaximumNArgs(1),
  RunE:runCommand,
}

func init(){
	f:=flags.GetFlags()
	dcf:= flags.DefaultValues.Check
	CheckCmd.Flags().StringVarP(&f.Check.OutputFormat,"output","o",dcf.OutputFormat,"output format (table|json)")
	CheckCmd.Flags().StringVar(&f.Check.ServerityLevel,"level",dcf.ServerityLevel,"")
	CheckCmd.Flags().StringVar(&f.Check.FailOn,"fail-on",dcf.FailOn,"exit code 1 when this severity found (error|warning)")
}

func runCommand(cmd *cobra.Command, args []string) error{
	// projet path
	root:="."
	if len(args)>0{
		root=args[0]
	}

	abs ,err:=filepath.Abs(root)
	if err!=nil{
		return fmt.Errorf("invalid project path: %w",err)
	}
	logger.Verbose(fmt.Sprintf("Checking project at: %s",abs))

	// chek if project exist
	if _,err := os.Stat(abs); os.IsNotExist(err){
		return fmt.Errorf("project path does not exist: %s",abs)
	}

	// load config
	logger.Verbose("Loading configuration...")
	f:=flags.GetFlags()
	cfg, err := config.Load(f.GlobalFlags.ConfigFile, abs)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	logger.Verbose(fmt.Sprintf("Configuration loaded: %d active rules", len(cfg.ActiveRules)))


	var detectionResult *detector.DetectionResult

	if cfg.Project.Type != "" {
		// Use specified type
		logger.Verbose(fmt.Sprintf("Using configured project type: %s", cfg.Project.Type))
		detectionResult, err = detector.DetectWithHint(abs, cfg.Project.Type)
		if err != nil {
			logger.Warning(fmt.Sprintf("Failed to detect as %s: %v", cfg.Project.Type, err))
			logger.Info("Falling back to auto-detection...")
			detectionResult, err = detector.Detect(abs)
			if err != nil {
				return fmt.Errorf("project detection failed: %w", err)
			}
		}
	} else {
		// Auto-detect
		logger.Verbose("Auto-detecting project type...")
		detectionResult, err = detector.Detect(abs)
		if err != nil {
			return fmt.Errorf("project detection failed: %w", err)
		}
	}

	cfg.Project.Type = detectionResult.Type.Primary
	logger.Verbose(fmt.Sprintf("Running %d rules...", len(cfg.Rules)))

	engine := rules.NewEngine(cfg,detectionResult)
	execResult,err :=engine.Execute()
	if err !=nil{
		return fmt.Errorf("rule execution failed: %w",err)
	}

	rep := reporter.New(f.Check.OutputFormat, execResult)
	if err := rep.Report(); err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	return determineExitCode(execResult, f.Check.FailOn)
}

func determineExitCode(result *rules.ExecutionResult, failOn string) error {
	f := flags.GetFlags()

	switch failOn {
	case "error":
		if result.Summary.Errors > 0 {
			// Exit with code 1 but don't show usage
			if !f.GlobalFlags.Quiet {
				fmt.Println() // Empty line before final message
				logger.Error(fmt.Sprintf("Validation failed: %d error(s) found", result.Summary.Errors))
			}
			os.Exit(1)
		}
	case "warning":
		if result.Summary.Errors > 0 || result.Summary.Warnings > 0 {
			if !f.GlobalFlags.Quiet {
				fmt.Println()
				logger.Error(fmt.Sprintf("Validation failed: %d error(s), %d warning(s) found",
					result.Summary.Errors, result.Summary.Warnings))
			}
			os.Exit(1)
		}
	}

	// Success
	if !f.GlobalFlags.Quiet {
		fmt.Println()
		logger.Success("Validation passed! ✓")
	}
	return nil
}
