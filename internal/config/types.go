package config
type ValidationError struct {
	Field		string
	Message		string
}
type ValidationResult struct{
	Valid		bool
	Errors		[]ValidationError
	Warnings	[]string
}

// rules structre
type RulesOptions struct {
	Enabled   bool     `yaml:"enabled"`
	Severity  string   `yaml:"severity"`
	Patterns  []string `yaml:"patterns"`
	Message   string   `yaml:"message"`
	FixHint   string   `yaml:"fix_hint"`
	Doc       string   `yaml:"doc"`
}
type GeneralType struct {
	Readme			RulesOptions	`yaml:"readme"`
	License			RulesOptions	`yaml:"license"`
	Gitignore		RulesOptions	`yaml:"gitignore"`
	Changelog		RulesOptions	`yaml:"changelog"`
}
type StructureType struct{
	Src				RulesOptions	 `yaml:"src"`
	Tests			RulesOptions	 `yaml:"tests"`
	Docs			RulesOptions	 `yaml:"docs"`
}
type DocumentationType struct{
	ADR				RulesOptions	 `yaml:"adrs"`
	Contributing	RulesOptions	 `yaml:"contributing"`
}
type CiCdType struct{
	GithubActions	RulesOptions     `yaml:"gihub_actions"`
	GitlabCi		RulesOptions	 `yaml:"gitlab_ci"`
}
type QualituType struct{
	PreCommit		RulesOptions     `yaml:"pre_commit"`
	EditorConfig    RulesOptions     `yaml:"editorconfig"`
}

type RulesType struct {
	General			GeneralType		 `yaml:"general"`
	Structure		StructureType	 `yaml:"structure"`
	Documentation	DocumentationType`yaml:"documentation"`
	CiCd			CiCdType         `yaml:"cicd"`
	Quality         QualituType		 `yaml:"quality"`
}

type ProjectType struct {
    Type	string		`yaml:"type"`
}
type Config struct {
    Version int         `yaml:"version"`
    Project ProjectType `yaml:"project"`
	Rules RulesType     `yaml:"rules"`
}

