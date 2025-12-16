package resources

import (
	"path/filepath"
)

func GetReadme(info *ProjectInfo, projectType string) string {
	vars := info.ToVars()
	vars["project_type"] = projectType

	template := templates.Readme[projectType]
	if template == "" {
		template = templates.Readme["generic"]
	}

	return replaceVars(template, vars)
}

func GetChangelog(info *ProjectInfo) string {
	vars := info.ToVars()
	return replaceVars(templates.Changelog, vars)
}

func GetContributing() string {
	return templates.Contributing
}

func GetFirstADR(info *ProjectInfo) string {
	vars := info.ToVars()
	return replaceVars(docsTemplates.ADRTemplate, vars)
}

func GetADRTemplate() string {
	return docsTemplates.ADRTemplate
}

func GetAPIDocs(info *ProjectInfo, projectType string) string {
	vars := info.ToVars()
	template := templates.APIDocs[projectType]
	if template == "" {
		template = templates.APIDocs["generic"]
	}
	return replaceVars(template, vars)
}

func GetTestExample(projectType string) string {
	if example, ok := templates.TestExamples[projectType]; ok {
		return example
	}
	return ""
}

func GetTestFileName(projectType string) string {
	names := map[string]string{
		"nodejs": "example.test.js",
		"go":     "example_test.go",
		"python": "test_example.py",
		"rust":   "example_test.rs",
	}
	if name, ok := names[projectType]; ok {
		return name
	}
	return ""
}

func GetGitignore(projectType string) string {
	common := gitignores.Common

	var specific string
	switch projectType {
	case "nodejs":
		specific = gitignores.NodeJS
	case "go":
		specific = gitignores.Go
	case "python":
		specific = gitignores.Python
	case "rust":
		specific = gitignores.Rust
	case "java":
		specific = gitignores.Java
	default:
		specific = ""
	}

	return common + "\n" + specific
}

func GetLicense(licenseType string, fullname string) string {
	license, ok := (*licenses)[licenseType]
	if !ok {
		license = (*licenses)["MIT"]
	}

	vars := getCurrentVars()
	vars["fullname"] = fullname
	if fullname == "" {
		vars["fullname"] = getUserName()
	}

	return replaceVars(license.Content, vars)
}

func ListLicenses() []string {
	result := make([]string, 0, len(*licenses))
	for key := range *licenses {
		result = append(result, key)
	}
	return result
}

func GetEditorconfig(projectType string) string {
	if config, ok := qualityTools.Editorconfig[projectType]; ok {
		return config
	}
	return qualityTools.Editorconfig["generic"]
}

func GetPrettierConfig() string {
	return qualityTools.Prettier.Config
}

func GetPrettierIgnore() string {
	return qualityTools.Prettier.Ignore
}

func GetESLintConfig(useTypeScript bool) string {
	if useTypeScript {
		return qualityTools.ESLint.TypeScript
	}
	return qualityTools.ESLint.Basic
}

func GetCommitlintConfig() string {
	return qualityTools.Commitlint.Config
}

func GetHuskyPreCommit() string {
	return qualityTools.Husky.PreCommit
}

func GetHuskyCommitMsg() string {
	return qualityTools.Husky.CommitMsg
}

func GetPreCommitConfig(projectType string) string {
	if config, ok := qualityTools.PreCommit[projectType]; ok {
		return config
	}
	return qualityTools.PreCommit["generic"]
}

func GetLintStagedConfig(projectType string) string {
	if config, ok := qualityTools.LintStaged[projectType]; ok {
		return config
	}
	return ""
}

func GetGitattributes(projectType string) string {
	if config, ok := qualityTools.Gitattributes[projectType]; ok {
		return config
	}
	return qualityTools.Gitattributes["generic"]
}

func GetMakefile(projectType string) string {
	if config, ok := qualityTools.Makefile[projectType]; ok {
		return config
	}
	return ""
}

func GetDockerfile(info *ProjectInfo, projectType string) string {
	vars := info.ToVars()

	var template string
	switch projectType {
	case "nodejs":
		template = devops.Docker.NodeJS.Dockerfile
	case "go":
		template = devops.Docker.Go.Dockerfile
	case "python":
		template = devops.Docker.Python.Dockerfile
	case "rust":
		template = devops.Docker.Rust.Dockerfile
	default:
		return ""
	}

	return replaceVars(template, vars)
}

func GetDockerignore(projectType string) string {
	switch projectType {
	case "nodejs":
		return devops.Docker.NodeJS.Dockerignore
	case "go":
		return devops.Docker.Go.Dockerignore
	case "python":
		return devops.Docker.Python.Dockerignore
	case "rust":
		return devops.Docker.Rust.Dockerignore
	default:
		return ""
	}
}

func GetDockerCompose(info *ProjectInfo, withDB bool) string {
	vars := info.ToVars()

	var template string
	if withDB {
		template = devops.DockerCompose["with_db"]
	} else {
		template = devops.DockerCompose["basic"]
	}

	return replaceVars(template, vars)
}

func GetKubernetesDeployment(info *ProjectInfo) string {
	vars := info.ToVars()
	return replaceVars(devops.Kubernetes.Deployment, vars)
}

func GetKubernetesService(info *ProjectInfo) string {
	vars := info.ToVars()
	return replaceVars(devops.Kubernetes.Service, vars)
}

func GetKubernetesIngress(info *ProjectInfo) string {
	vars := info.ToVars()
	return replaceVars(devops.Kubernetes.Ingress, vars)
}

func GetNginxConfig(info *ProjectInfo, configType string) string {
	vars := info.ToVars()

	var template string
	switch configType {
	case "static":
		template = devops.Nginx["static"]
	case "reverse_proxy":
		template = devops.Nginx["reverse_proxy"]
	default:
		template = devops.Nginx["static"]
	}

	return replaceVars(template, vars)
}

func GetGitHubActionsCI(info *ProjectInfo, projectType string) string {
	vars := info.ToVars()

	var template string
	switch projectType {
	case "nodejs":
		template = devops.GitHubActions["nodejs_ci"]
	case "go":
		template = devops.GitHubActions["go_ci"]
	default:
		return ""
	}

	return replaceVars(template, vars)
}

func GetGitHubActionsDocker(info *ProjectInfo) string {
	vars := info.ToVars()
	return replaceVars(devops.GitHubActions["docker_build"], vars)
}

func GetRenovateConfig() string {
	return devops.Renovate.Config
}

func GetDependabotConfig(projectType string) string {
	// Base config
	config := devops.Dependabot["config"]

	if langConfig, ok := devops.Dependabot[projectType]; ok {
		config += "\n" + langConfig
	}

	return config
}

func GetSecurity(info *ProjectInfo) string {
	vars := info.ToVars()
	return replaceVars(docsTemplates.Security, vars)
}

func GetCodeOfConduct(info *ProjectInfo) string {
	vars := info.ToVars()
	return replaceVars(docsTemplates.CodeOfConduct, vars)
}

func GetPullRequestTemplate() string {
	return docsTemplates.PullRequestTemplate
}

func GetIssueBugReport() string {
	return docsTemplates.IssueBugReport
}

func GetIssueFeatureRequest() string {
	return docsTemplates.IssueFeatureRequest
}

func GetIssueQuestion() string {
	return docsTemplates.IssueQuestion
}

func GetIssueTemplatesConfig(info *ProjectInfo) string {
	vars := info.ToVars()
	return replaceVars(docsTemplates.IssueTemplatesConfig, vars)
}

func GetCodeowners(info *ProjectInfo) string {
	vars := info.ToVars()
	return replaceVars(docsTemplates.Codeowners, vars)
}

func GetSupport(info *ProjectInfo) string {
	vars := info.ToVars()
	return replaceVars(docsTemplates.Support, vars)
}

func GetRoadmap(info *ProjectInfo) string {
	vars := info.ToVars()
	return replaceVars(docsTemplates.Roadmap, vars)
}

func GetFunding(info *ProjectInfo) string {
	vars := info.ToVars()
	return replaceVars(docsTemplates.Funding, vars)
}

func GetFundingYML(info *ProjectInfo) string {
	vars := info.ToVars()
	return replaceVars(docsTemplates.FundingYML, vars)
}

func GetEnvExample(projectType string) string {
	if env, ok := docsTemplates.EnvExample[projectType]; ok {
		return env
	}
	return docsTemplates.EnvExample["generic"]
}

func GetInstallScript(info *ProjectInfo, projectType string) string {
	vars := info.ToVars()

	// Add install command based on project type
	switch projectType {
	case "nodejs":
		vars["install_command"] = "npm install"
		vars["start_command"] = "npm start"
	case "go":
		vars["install_command"] = "go mod download"
		vars["start_command"] = "go run ."
	case "python":
		vars["install_command"] = "pip install -r requirements.txt"
		vars["start_command"] = "python main.py"
	case "rust":
		vars["install_command"] = "cargo build"
		vars["start_command"] = "cargo run"
	default:
		vars["install_command"] = "make install"
		vars["start_command"] = "make start"
	}

	return replaceVars(scripts.Install["unix"], vars)
}

func GetSetupScript(projectType string) string {
	if script, ok := scripts.Setup[projectType]; ok {
		return script
	}
	return scripts.Setup["unix"]
}

func GetTestScript(info *ProjectInfo, projectType string) string {
	vars := info.ToVars()

	switch projectType {
	case "nodejs":
		vars["test_command"] = "npm test"
	case "go":
		vars["test_command"] = "go test ./..."
	case "python":
		vars["test_command"] = "pytest"
	case "rust":
		vars["test_command"] = "cargo test"
	default:
		vars["test_command"] = "make test"
	}

	return replaceVars(scripts.Test["unix"], vars)
}

func GetBuildScript(info *ProjectInfo, projectType string) string {
	vars := info.ToVars()

	switch projectType {
	case "nodejs":
		vars["build_dir"] = "dist"
		vars["build_command"] = "npm run build"
	case "go":
		vars["build_dir"] = "build"
		vars["build_command"] = "go build -o build/ ./..."
	case "python":
		vars["build_dir"] = "dist"
		vars["build_command"] = "python setup.py build"
	case "rust":
		vars["build_dir"] = "target/release"
		vars["build_command"] = "cargo build --release"
	default:
		vars["build_dir"] = "build"
		vars["build_command"] = "make build"
	}

	return replaceVars(scripts.Build["unix"], vars)
}

func GetDeployScript(info *ProjectInfo, projectType string) string {
	vars := info.ToVars()

	switch projectType {
	case "nodejs":
		vars["build_command"] = "npm run build"
		vars["test_command"] = "npm test"
		vars["deploy_command"] = "npm run deploy"
	case "go":
		vars["build_command"] = "make build"
		vars["test_command"] = "make test"
		vars["deploy_command"] = "make deploy"
	default:
		vars["build_command"] = "make build"
		vars["test_command"] = "make test"
		vars["deploy_command"] = "make deploy"
	}

	return replaceVars(scripts.Deploy["unix"], vars)
}

func GetReleaseScript(info *ProjectInfo, projectType string) string {
	vars := info.ToVars()

	switch projectType {
	case "nodejs":
		vars["test_command"] = "npm test"
		vars["build_command"] = "npm run build"
		vars["version_update_command"] = "npm version ${VERSION}"
	case "go":
		vars["test_command"] = "make test"
		vars["build_command"] = "make build"
		vars["version_update_command"] = "echo ${VERSION} > VERSION"
	default:
		vars["test_command"] = "make test"
		vars["build_command"] = "make build"
		vars["version_update_command"] = "echo ${VERSION} > VERSION"
	}

	return replaceVars(scripts.Release["unix"], vars)
}

func GetDockerBuildScript(info *ProjectInfo) string {
	vars := info.ToVars()
	return replaceVars(scripts.DockerBuild["unix"], vars)
}

func GetCleanScript(info *ProjectInfo, projectType string) string {
	vars := info.ToVars()

	switch projectType {
	case "nodejs":
		vars["build_dirs"] = "dist/ build/ .next/ out/"
		vars["cache_dirs"] = "node_modules/ .cache/"
	case "go":
		vars["build_dirs"] = "build/ bin/"
		vars["cache_dirs"] = "vendor/"
	case "python":
		vars["build_dirs"] = "build/ dist/ *.egg-info/"
		vars["cache_dirs"] = "__pycache__/ .pytest_cache/ .coverage"
	case "rust":
		vars["build_dirs"] = "target/"
		vars["cache_dirs"] = "Cargo.lock"
	default:
		vars["build_dirs"] = "build/ dist/"
		vars["cache_dirs"] = ".cache/"
	}

	return replaceVars(scripts.Clean["unix"], vars)
}

func GetLintScript(info *ProjectInfo, projectType string) string {
	vars := info.ToVars()

	switch projectType {
	case "nodejs":
		vars["lint_command"] = "npm run lint"
	case "go":
		vars["lint_command"] = "golangci-lint run"
	case "python":
		vars["lint_command"] = "flake8 ."
	case "rust":
		vars["lint_command"] = "cargo clippy"
	default:
		vars["lint_command"] = "make lint"
	}

	return replaceVars(scripts.Lint["unix"], vars)
}

func GetFormatScript(info *ProjectInfo, projectType string) string {
	vars := info.ToVars()

	switch projectType {
	case "nodejs":
		vars["format_command"] = "npm run format"
	case "go":
		vars["format_command"] = "gofmt -w ."
	case "python":
		vars["format_command"] = "black ."
	case "rust":
		vars["format_command"] = "cargo fmt"
	default:
		vars["format_command"] = "make format"
	}

	return replaceVars(scripts.Format["unix"], vars)
}

func GetDevScript(info *ProjectInfo, projectType string) string {
	vars := info.ToVars()

	switch projectType {
	case "nodejs":
		vars["deps_dir"] = "node_modules"
		vars["install_command"] = "npm install"
		vars["dev_command"] = "npm run dev"
	case "go":
		vars["deps_dir"] = "vendor"
		vars["install_command"] = "go mod download"
		vars["dev_command"] = "go run ."
	case "python":
		vars["deps_dir"] = "venv"
		vars["install_command"] = "pip install -r requirements.txt"
		vars["dev_command"] = "python main.py"
	case "rust":
		vars["deps_dir"] = "target"
		vars["install_command"] = "cargo fetch"
		vars["dev_command"] = "cargo run"
	default:
		vars["deps_dir"] = "vendor"
		vars["install_command"] = "make install"
		vars["dev_command"] = "make dev"
	}

	return replaceVars(scripts.Dev["unix"], vars)
}

func GetProjectName(projectPath string) string {
	return filepath.Base(projectPath)
}
