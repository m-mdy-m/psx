package shared
import (
	"sync"
)
type GlobalFlags struct {
	ConfigFile string
	Verbose bool
	Quiet bool
	NoColor bool
}


// singletoon
var (
	flags *GlobalFlags
	once sync.Once
)
func GetGlobalFlags() *GlobalFlags{
	once.Do(func(){
		flags = &GlobalFlags{
			ConfigFile:"",
			Verbose:false,
			Quiet:false,
			NoColor:false,
		}
	})
	return flags
}
// getters
func (g*GlobalFlags) GetConfigFile() string{
	return g.ConfigFile
}
func (g*GlobalFlags) IsVerbose() bool{
	return g.Verbose && !g.Quiet
}
func (g*GlobalFlags) IsQuiet() bool{
	return g.Quiet
}
func (g*GlobalFlags) IsColorEnabled() bool{
	if g.NoColor{
		return false
	}
	return true
}
// setters
func (g*GlobalFlags) SetConfigFile(path string){
	g.ConfigFile = path
}
func (g*GlobalFlags) SetVerbose(v bool){
	g.Verbose = v
}
func (g*GlobalFlags) SetQuiet(q bool){
	g.Quiet = q
	if q{
		g.Verbose = false
	}
}
func (g*GlobalFlags) SetNoColor(c bool){
	g.NoColor = c
}

