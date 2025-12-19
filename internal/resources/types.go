package resources

type ProjectInfo struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Author      string `yaml:"author"`
	Email       string `yaml:"email"`
	GitHubUser  string `yaml:"github_user,omitempty"`
	RepoName    string `yaml:"repo_name,omitempty"`
	License     string `yaml:"license"`

	// Derived fields (not in YAML)
	RepoURL     string `yaml:"-"`
	Domain      string `yaml:"-"`
	DockerImage string `yaml:"-"`
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
	APIDocs      map[string]string `yaml:"api_docs"`
}

type QualityToolsConfig struct {
	Editorconfig map[string]string `yaml:"editorconfig"`
	PreCommit    map[string]string `yaml:"pre_commit"`
}

type DevOpsConfig struct {
	Docker        DockerConfig      `yaml:"docker"`
	DockerCompose map[string]string `yaml:"docker_compose"`
	CICD          CICDConfig        `yaml:"cicd"`
}
type CICDConfig struct {
	GitHubActions map[string]string `yaml:"github_actions"`
	GitLabCI      map[string]string `yaml:"gitlab_ci"`
}


type DockerConfig struct {
	NodeJS  DockerLanguageConfig `yaml:"nodejs"`
	Go      DockerLanguageConfig `yaml:"go"`
	Generic DockerLanguageConfig `yaml:"generic"`
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
	ADRTemplates         map[string]string `yaml:"adr"`
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

type LanguagesConfig struct {
	Aliases map[string]string `yaml:"aliases"`
}
