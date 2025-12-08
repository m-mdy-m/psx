package commond

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/m-mdy-m/psx/internal/config"
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
	config.Load(f.GlobalFlags.ConfigFile,abs)
	return nil
}

