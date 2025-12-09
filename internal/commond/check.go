package commond

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/detector"
	"github.com/m-mdy-m/psx/internal/flags"
	"github.com/m-mdy-m/psx/internal/shared/logger"
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

	// === STEP 2: Detect Project Type ===
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

	// Update effective config with detected type
	cfg.Project.Type = detectionResult.Type.Primary

	// === STEP 3: Display Project Info ===
	displayProjectInfo(detectionResult)

	// === STEP 4: Run Validation ===
	logger.Verbose(fmt.Sprintf("Running %d rules...", len(cfg.Rules)))

	logger.Success("Check complete!")
	return nil
}

func displayProjectInfo(result *detector.DetectionResult) {
	fmt.Println("\nProject Information:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Type: %s\n", result.Description)

	if result.Type.Framework != "" {
		fmt.Printf("Framework: %s\n", result.Type.Framework)
	}

	if result.Type.PackageManager != "" {
		fmt.Printf("Package Manager: %s\n", result.Type.PackageManager)
	}

	fmt.Printf("Structure: %s\n", result.Type.Structure)
	fmt.Printf("Confidence: %.0f%%\n", result.Type.Confidence*100)

	if len(result.Type.Secondary) > 0 {
		fmt.Printf("Also detected: %v\n", result.Type.Secondary)
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
}


