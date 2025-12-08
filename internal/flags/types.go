package flags
type GlobalFlags struct {
    ConfigFile string
    Verbose    bool
    Quiet      bool
    NoColor    bool
}
type Check struct {
	OutputFormat     string
	ServerityLevel   string
	FailOn			 string
}
type Flags struct {
    GlobalFlags  GlobalFlags
	Check        Check
}

var DefaultValues = Flags{
    GlobalFlags: GlobalFlags{
        ConfigFile: "",
        Verbose:    false,
        Quiet:      false,
        NoColor:    false,
    },
	Check:Check{
		OutputFormat:	"table",
		ServerityLevel: "all",
		FailOn:			"error",
	},
}

