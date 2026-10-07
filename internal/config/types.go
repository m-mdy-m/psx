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

	// Resolved at load time, not part of the YAML document.
	ProjectPath string `yaml:"-"`
	ConfigFile  string `yaml:"-"`
	ProjectType string `yaml:"-"`
	// DeclaredType is the raw project.type as written in the file, before
	// normalisation. It is what tells "the user chose generic" apart from
	// "the user chose nothing", which is the difference between trusting the
	// value and running detection.
	DeclaredType string                 `yaml:"-"`
	ActiveRules  map[string]*ActiveRule `yaml:"-"`
}

// RULES METADATA (GLOBAL)
type LanguagePatterns any

// AdditionalCheck is a secondary assertion such as "package.json:license".
type AdditionalCheck struct {
	File  string `yaml:"file"`
	Field string `yaml:"field"`
}

type FixSpec struct {
	// Path is the single file or directory to create.
	Path  string            `yaml:"path,omitempty"`
	Files map[string]string `yaml:"files,omitempty"`
	// Template names the resource template used for Path.
	Template string `yaml:"template,omitempty"`
	// Mode is the permission bitmask for created files; defaults to 0644.
	Mode  uint32   `yaml:"mode,omitempty"`
	Safe  bool     `yaml:"safe,omitempty"`
	Needs []string `yaml:"needs,omitempty"`
	// Content is an inline template body, used when no shared template exists.
	Content string `yaml:"content,omitempty"`
}

// RuleMetadata is the declarative definition of a rule, loaded from rules.yml.
type RuleMetadata struct {
	ID               string            `yaml:"id"`
	Category         string            `yaml:"category"`
	Description      string            `yaml:"description"`
	DefaultSeverity  Severity          `yaml:"severity"`
	Patterns         LanguagePatterns  `yaml:"patterns"`
	AdditionalChecks []AdditionalCheck `yaml:"additional_checks,omitempty"`
	Fix              *FixSpec          `yaml:"fix,omitempty"`
	Message          string            `yaml:"message"`
	FixHint          string            `yaml:"fix_hint"`
	DocURL           string            `yaml:"doc_url"`
}

// Fixable reports whether the rule declares an automatic fix.
func (m RuleMetadata) Fixable() bool {
	return m.Fix != nil && (m.Fix.Path != "" || len(m.Fix.Files) > 0)
}

// FileFixes returns the literal-path to template mapping for a multi-file fix.
func (m RuleMetadata) FileFixes() map[string]string {
	if m.Fix == nil {
		return nil
	}
	return m.Fix.Files
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
