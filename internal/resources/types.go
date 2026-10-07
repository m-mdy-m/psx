package resources

// ProjectInfo describes the project being scaffolded.
// RepoURL, Domain and DockerImage are derived, never read from YAML.
type ProjectInfo struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Author      string `yaml:"author"`
	Email       string `yaml:"email"`
	GitHubUser  string `yaml:"github_user"`
	RepoName    string `yaml:"repo_name"`
	License     string `yaml:"license"`

	RepoURL     string `yaml:"-"`
	Domain      string `yaml:"-"`
	DockerImage string `yaml:"-"`
}

// MessagesConfig holds user-facing strings.
type MessagesConfig struct {
	Errors  map[string]string `yaml:"errors"`
	Check   map[string]string `yaml:"check"`
	Fix     map[string]string `yaml:"fix"`
	Init    map[string]string `yaml:"init"`
	Detect  map[string]string `yaml:"detect"`
	Verbose map[string]string `yaml:"verbose"`
}

// LicensesConfig maps a SPDX id to its full text.
type LicensesConfig map[string]LicenseTemplate

// LicenseTemplate is one license body.
type LicenseTemplate struct {
	Name    string `yaml:"name"`
	Content string `yaml:"content"`
}

// TemplatesConfig holds documentation scaffolds.
type TemplatesConfig struct {
	Readme       map[string]string `yaml:"readme"`
	APIDocs      map[string]string `yaml:"api_docs"`
	Changelog    string            `yaml:"changelog"`
	Contributing string            `yaml:"contributing"`
	DocsIndex    string            `yaml:"docs_index"`
}

// GitignoresConfig holds ignore scaffolds plus the environment template.
type GitignoresConfig struct {
	Common     string            `yaml:"common"`
	NodeJS     string            `yaml:"nodejs"`
	Go         string            `yaml:"go"`
	EnvExample map[string]string `yaml:"env_example"`
}

// QualityToolsConfig holds editor and hook configuration scaffolds.
type QualityToolsConfig struct {
	Editorconfig  map[string]string `yaml:"editorconfig"`
	PreCommit     map[string]string `yaml:"pre_commit"`
	Gitattributes string            `yaml:"gitattributes"`
	Makefile      map[string]string `yaml:"makefile"`
}

// DevOpsConfig holds container, CI and dependency-management scaffolds.
type DevOpsConfig struct {
	Docker          DockerConfig      `yaml:"docker"`
	DockerCompose   map[string]string `yaml:"docker_compose"`
	CICD            CICDConfig        `yaml:"cicd"`
	Dependabot      map[string]string `yaml:"dependabot"`
	Renovate        map[string]string `yaml:"renovate"`
	NodeVersion     map[string]string `yaml:"node_version"`
	ReleaseWorkflow map[string]string `yaml:"release_workflow"`
	Nginx           map[string]string `yaml:"nginx"`
	HelmChart       map[string]string `yaml:"helm_chart"`
	Kubernetes      map[string]string `yaml:"kubernetes"`
}

// CICDConfig holds pipeline scaffolds per platform.
type CICDConfig struct {
	GitHubActions map[string]string `yaml:"github_actions"`
	GitLabCI      map[string]string `yaml:"gitlab_ci"`
}

// DockerConfig holds Dockerfiles and ignore files per language.
type DockerConfig struct {
	NodeJS  DockerLanguageConfig `yaml:"nodejs"`
	Go      DockerLanguageConfig `yaml:"go"`
	Generic DockerLanguageConfig `yaml:"generic"`
}

// DockerLanguageConfig is a Dockerfile plus its dockerignore.
type DockerLanguageConfig struct {
	Dockerfile   string `yaml:"dockerfile"`
	Dockerignore string `yaml:"dockerignore"`
}

// DocsTemplatesConfig holds governance and policy documents.
type DocsTemplatesConfig struct {
	Security             string            `yaml:"security"`
	CodeOfConduct        string            `yaml:"code_of_conduct"`
	Support              string            `yaml:"support"`
	PullRequestTemplate  string            `yaml:"pull_request_template"`
	IssueBugReport       string            `yaml:"issue_bug_report"`
	IssueFeatureRequest  string            `yaml:"issue_feature_request"`
	IssueQuestion        string            `yaml:"issue_question"`
	IssueTemplatesConfig string            `yaml:"issue_templates_config"`
	Codeowners           string            `yaml:"codeowners"`
	Funding              string            `yaml:"funding"`
	Roadmap              string            `yaml:"roadmap"`
	Architecture         string            `yaml:"architecture"`
	Runbook              string            `yaml:"runbook"`
	OpenAPI              string            `yaml:"openapi"`
	ADRTemplates         map[string]string `yaml:"adr"`
}

// ScriptsConfig holds developer script scaffolds keyed by platform or language.
type ScriptsConfig struct {
	Install     map[string]string `yaml:"install"`
	Setup       map[string]string `yaml:"setup"`
	Test        map[string]string `yaml:"test"`
	Build       map[string]string `yaml:"build"`
	Release     map[string]string `yaml:"release"`
	DockerBuild map[string]string `yaml:"docker_build"`
	Clean       map[string]string `yaml:"clean"`
}

// ProjectScriptsConfig holds the language-neutral scripts a fix scaffolds.
type ProjectScriptsConfig struct {
	Setup string `yaml:"script_setup"`
	Test  string `yaml:"script_test"`
	Build string `yaml:"script_build"`
	Clean string `yaml:"script_clean"`
}

// ScriptPlatformConfig is a per-platform script map.
type ScriptPlatformConfig map[string]string

// LanguagesConfig carries language aliases and per-language metadata.
type LanguagesConfig struct {
	Aliases   map[string]string   `yaml:"aliases"`
	Languages map[string]Language `yaml:"languages"`
}

// Language is the metadata block for one language.
type Language struct {
	Name            string           `yaml:"name"`
	PackageManagers []PackageManager `yaml:"package_managers"`
	TestPatterns    []string         `yaml:"test_patterns"`
	SrcPatterns     []string         `yaml:"src_patterns"`
	BuildPatterns   []string         `yaml:"build_patterns"`
	DistPatterns    []string         `yaml:"dist_patterns"`
	CachePatterns   []string         `yaml:"cache_patterns"`
	InstallCommands []string         `yaml:"install_commands"`
	TestCommands    []string         `yaml:"test_commands"`
	BuildCommands   []string         `yaml:"build_commands"`
	StartCommands   []string         `yaml:"start_commands"`
}

// PackageManager describes one dependency manager for a language.
type PackageManager struct {
	Name       string `yaml:"name"`
	File       string `yaml:"file"`
	Lock       string `yaml:"lock"`
	InstallCmd string `yaml:"install_cmd"`
	Priority   int    `yaml:"priority"`
}
