package commond

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/m-mdy-m/psx/internal/flags"

)

var root=&cobra.Command{
	Use: "psx",
	Short: "PSX -Porject Structure X",
	Long: `PSX is a CLI tool that validates and stadradizes project structures across different programing languages`,
	Version: "", // write with IdFlags
	PersistentPreRun: func (cmd * cobra.Command, args [] string) {
		gf:= flags.GetFlags().GlobalFlags

		if  gf.IsColorEnabled() {
			color.NoColor= true
		}
	},
}

func Exce(v string) error{
	root.Version = v
	return root.Execute()
}

// intial global flgs for every ccomand
func initGlobalFlags(){
	f:=flags.GetFlags()

	defaultValues := flags.DefaultValues.GlobalFlags

	// config flag
	root.PersistentFlags().StringVar(&f.GlobalFlags.ConfigFile,"config",defaultValues.ConfigFile,"config file path (deault: auto discover psx.yml)")

	// verbose flag
	root.PersistentFlags().BoolVar(&f.GlobalFlags.Verbose,"verbose",defaultValues.Verbose,"enable veerbose output")

	// quiet
	root.PersistentFlags().BoolVar(&f.GlobalFlags.Quiet,"quiet",defaultValues.Quiet,"enable quiet mode")

	// color
	root.PersistentFlags().BoolVar(&f.GlobalFlags.NoColor,"no-color",defaultValues.NoColor,"disable colored output")
}
func init() {
    initGlobalFlags()
	root.AddCommand(CheckCmd)
}

