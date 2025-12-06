package logger
// pkg.dev/github.com/fatih/color
import (
	"fmt"

	"github.com/fatih/color"
)
func Info(m string){
	fmt.Println(m)
}
func Verbose(m string){
	color.Cyan("→ "+m)
}
func Error(m string){
	color.Red("✗ "+m)
}
func Success(m string){
	color.Green("✓ "+m)
}
func Warning(m string){
	color.Yellow("⚠ "+m)
}
