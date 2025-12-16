package resources

type LanguagesConfig struct {
	Aliases      map[string]string          `yaml:"aliases"`
	Folders      map[string]LanguageFolders `yaml:"folders"`
	TestPatterns map[string][]string        `yaml:"test_patterns"`
}

type LanguageFolders struct {
	Src   []string `yaml:"src"`
	Tests []string `yaml:"tests"`
}

type MessagesConfig struct {
	Exit    map[string]string `yaml:"exit"`
	Errors  map[string]string `yaml:"errors"`
	Help    map[string]string `yaml:"help"`
	Check   map[string]string `yaml:"check"`
	Fix     map[string]string `yaml:"fix"`
	Init    map[string]string `yaml:"init"`
	Verbose map[string]string `yaml:"verbose"`
}

type GitignoresConfig struct {
	Common string `yaml:"common"`
	NodeJS string `yaml:"nodejs"`
	Go     string `yaml:"go"`
	Python string `yaml:"python"`
	Rust   string `yaml:"rust"`
	Java   string `yaml:"java"`
}

type LicensesConfig map[string]LicenseTemplate

type LicenseTemplate struct {
	Name    string `yaml:"name"`
	Content string `yaml:"content"`
}

type TemplatesConfig struct {
	Readme       map[string]string `yaml:"readme"`
	Changelog    string            `yaml:"changelog"`
	Contributing string            `yaml:"contributing"`
	ADR          map[string]string `yaml:"adr"`
	APIDocs      map[string]string `yaml:"api_docs"`
	TestExamples map[string]string `yaml:"test_examples"`
	Editorconfig string            `yaml:"editorconfig"`
}

type QualityToolsConfig struct {
	Editorconfig  map[string]string `yaml:"editorconfig"`
	Prettier      PrettierConfig    `yaml:"prettier"`
	ESLint        ESLintConfig      `yaml:"eslint"`
	Commitlint    CommitlintConfig  `yaml:"commitlint"`
	Husky         HuskyConfig       `yaml:"husky"`
	PreCommit     map[string]string `yaml:"pre_commit"`
	LintStaged    map[string]string `yaml:"lint_staged"`
	Gitattributes map[string]string `yaml:"gitattributes"`
	Makefile      map[string]string `yaml:"makefile"`
}

type PrettierConfig struct {
	Config string `yaml:"config"`
	Ignore string `yaml:"ignore"`
}

type ESLintConfig struct {
	Basic      string `yaml:"basic"`
	TypeScript string `yaml:"typescript"`
}

type CommitlintConfig struct {
	Config string `yaml:"config"`
}

type HuskyConfig struct {
	PreCommit string `yaml:"pre_commit"`
	CommitMsg string `yaml:"commit_msg"`
}

type DevOpsConfig struct {
	Docker        DockerConfig      `yaml:"docker"`
	DockerCompose map[string]string `yaml:"docker_compose"`
	Kubernetes    KubernetesConfig  `yaml:"kubernetes"`
	Nginx         map[string]string `yaml:"nginx"`
	GitHubActions map[string]string `yaml:"github_actions"`
	Renovate      RenovateConfig    `yaml:"renovate"`
	Dependabot    map[string]string `yaml:"dependabot"`
}

type DockerConfig struct {
	NodeJS DockerLanguageConfig `yaml:"nodejs"`
	Go     DockerLanguageConfig `yaml:"go"`
	Python DockerLanguageConfig `yaml:"python"`
	Rust   DockerLanguageConfig `yaml:"rust"`
}

type DockerLanguageConfig struct {
	Dockerfile   string `yaml:"dockerfile"`
	Dockerignore string `yaml:"dockerignore"`
}

type KubernetesConfig struct {
	Deployment string `yaml:"deployment"`
	Service    string `yaml:"service"`
	Ingress    string `yaml:"ingress"`
}

type RenovateConfig struct {
	Config string `yaml:"config"`
}

type DocsTemplatesConfig struct {
	Security             string            `yaml:"security"`
	CodeOfConduct        string            `yaml:"code_of_conduct"`
	PullRequestTemplate  string            `yaml:"pull_request_template"`
	IssueBugReport       string            `yaml:"issue_bug_report"`
	IssueFeatureRequest  string            `yaml:"issue_feature_request"`
	IssueQuestion        string            `yaml:"issue_question"`
	IssueTemplatesConfig string            `yaml:"issue_templates_config"`
	Codeowners           string            `yaml:"codeowners"`
	Support              string            `yaml:"support"`
	Roadmap              string            `yaml:"roadmap"`
	Funding              string            `yaml:"funding"`
	FundingYML           string            `yaml:"funding_yml"`
	ADRTemplate          string            `yaml:"adr_template"`
	EnvExample           map[string]string `yaml:"env_example"`
}

type ScriptsConfig struct {
	Install     ScriptPlatformConfig `yaml:"install"`
	Setup       map[string]string    `yaml:"setup"`
	Test        ScriptPlatformConfig `yaml:"test"`
	Build       ScriptPlatformConfig `yaml:"build"`
	Deploy      ScriptPlatformConfig `yaml:"deploy"`
	Release     ScriptPlatformConfig `yaml:"release"`
	DockerBuild ScriptPlatformConfig `yaml:"docker_build"`
	Clean       ScriptPlatformConfig `yaml:"clean"`
	Lint        ScriptPlatformConfig `yaml:"lint"`
	Format      ScriptPlatformConfig `yaml:"format"`
	Dev         ScriptPlatformConfig `yaml:"dev"`
}

type ScriptPlatformConfig map[string]string

type ProjectInfo struct {
	Name          string
	Description   string
	Author        string
	Email         string
	GitHubUser    string
	RepoName      string
	RepoURL       string
	License       string
	Domain        string
	DockerImage   string
	SupportEmail  string
	SecurityEmail string
	CurrentDir    string
}
