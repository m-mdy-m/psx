package config

type ValidationError struct {
	Field   string
	Message string
}
type ValidationResult struct {
	Valid    bool
	Errors   []ValidationError
	Warnings []string
}

// rules structre
// can be: "error","warning","info", or false (disbled)
type Severity string
type RulesSeverity any

type ProjectType struct {
	Type string `yaml:"type"`
}

type FixConfig struct {
	Interactive bool `yaml:"interactive"`
	Backup      bool `yaml:"backup"`
}

type ActiveRule struct {
	ID       string
	Metadata RuleMetadata
	Severity Severity
}

type Config struct {
	Version int                      `yaml:"version"`
	Project ProjectType              `yaml:"project"`
	Rules   map[string]RulesSeverity `yaml:"rules"`
	Ignore  []string                 `yaml:"ignore,omitempty"`
	Fix     FixConfig                `yaml:"fix,omitempty"`
	Custom  *CustomConfig            `yaml:"custom,omitempty"`

	// not in yml file
	Path        string                 `yaml:"-"`
	ActiveRules map[string]*ActiveRule `yaml:"-"`
}

// RULES METADATA (GLOBAL)
type LanguagePatterns any

// AdditionalCheck represents a check in a specific file
type AdditionalCheck struct {
	File  string // e.g., "package.json"
	Field string // e.g., "license"
}

// RuleMetadata contains all information about a rule
type RuleMetadata struct {
	ID               string           `yaml:"id"`
	Category         string           `yaml:"category"`
	Description      string           `yaml:"description"`
	DefaultSeverity  Severity         `yaml:"severity"`
	Patterns         LanguagePatterns `yaml:"patterns"` // []string or LanguagePatterns
	AdditionalChecks []string         `yaml:"additional_checks,omitempty"`
	Message          string           `yaml:"message"`
	FixHint          string           `yaml:"fix_hint"`
	DocURL           string           `yaml:"doc_url"`
}

// RulesMetadata contains all rule definitions
type RulesMetadata struct {
	Rules map[string]RuleMetadata `yaml:"rules"`
}

type RuleConfig struct {
	Metadata RuleMetadata
	Severity *Severity
}

type CustomConfig struct {
	Files   []CustomFile   `yaml:"files"`
	Folders []CustomFolder `yaml:"folders"`
}
type CustomFile struct {
	Path    string `yaml:"path"`
	Content string `yaml:"content"`
}
type CustomFolder struct {
	Path      string                 `yaml:"path"`
	Structure map[string]interface{} `yaml:"structure"`
}
