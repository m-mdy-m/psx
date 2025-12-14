package logger

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/m-mdy-m/psx/internal/flags"
)
func Info(m string){
	f := flags.GetFlags()
	if f.GlobalFlags.IsQuiet(){
		fmt.Println(m)
	}
}
func Verbose(m string){
	f := flags.GetFlags()
	if f.GlobalFlags.Verbose{
		color.Cyan("→ "+m)
	}
}
func Error(m string){
	color.Red("✗ "+m)
}
func Errorf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	Error("✗ "+msg)
}
func Success(m string){
	f := flags.GetFlags()
	if !f.GlobalFlags.IsQuiet(){
		color.Green("✓ "+m)
	}
}
func Warning(m string){
	f := flags.GetFlags()
	if !f.GlobalFlags.IsQuiet(){
		color.Yellow("⚠ "+m)
	}
}
