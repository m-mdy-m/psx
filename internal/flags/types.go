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

type Fix struct {
	Interactive   bool
	DryRun        bool
	RuleID        string
	All           bool
	CreateBackups bool
}

type Init struct {
	Template string
	Minimal  bool
	Force    bool
}

type Rules struct {
	Category string
	JSON     bool
}

type Flags struct {
	GlobalFlags GlobalFlags
	Check       Check
	Fix         Fix
	Init        Init
	Rules       Rules
}

var DefaultValues = Flags{
	GlobalFlags: GlobalFlags{
		ConfigFile: "",
		Verbose:    false,
		Quiet:      false,
		NoColor:    false,
	},
	Check: Check{
		OutputFormat:   "table",
		ServerityLevel: "all",
		FailOn:         "error",
	},
	Fix: Fix{
		Interactive:   true,
		DryRun:        false,
		RuleID:        "",
		All:           false,
		CreateBackups: false,
	},
	Init: Init{
		Template: "",
		Minimal:  false,
		Force:    false,
	},
	Rules: Rules{
		Category: "",
		JSON:     false,
	},
}
