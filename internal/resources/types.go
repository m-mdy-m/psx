package resources

// Languages config structure
type LanguagesConfig struct {
	Aliases      map[string]string            `yaml:"aliases"`
	Folders      map[string]LanguageFolders   `yaml:"folders"`
	TestPatterns map[string][]string          `yaml:"test_patterns"`
}

type LanguageFolders struct {
	Src   []string `yaml:"src"`
	Tests []string `yaml:"tests"`
}

// Messages config structure
type MessagesConfig struct {
	Exit    map[string]string           `yaml:"exit"`
	Errors  map[string]string           `yaml:"errors"`
	Help    map[string]string           `yaml:"help"`
	Check   map[string]string           `yaml:"check"`
	Fix     map[string]string           `yaml:"fix"`
	Init    map[string]string           `yaml:"init"`
	Verbose map[string]string           `yaml:"verbose"`
}

// Gitignores config structure
type GitignoresConfig struct {
	Common string `yaml:"common"`
	NodeJS string `yaml:"nodejs"`
	Go     string `yaml:"go"`
	Python string `yaml:"python"`
	Rust   string `yaml:"rust"`
	Java   string `yaml:"java"`
}

// Licenses config structure
type LicensesConfig map[string]LicenseTemplate

type LicenseTemplate struct {
	Name    string `yaml:"name"`
	Content string `yaml:"content"`
}

// Templates config structure
type TemplatesConfig struct {
	Readme       map[string]string `yaml:"readme"`
	Changelog    string            `yaml:"changelog"`
	Contributing string            `yaml:"contributing"`
	ADR          map[string]string `yaml:"adr"`
	APIDocs      map[string]string `yaml:"api_docs"`
	TestExamples map[string]string `yaml:"test_examples"`
	Editorconfig string            `yaml:"editorconfig"`
}