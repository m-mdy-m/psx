package resources

// === Template Getters ===

func GetReadme(info *ProjectInfo, projectType string) string {
	if info == nil {
		info = getDefaultProjectInfo()
	}
	vars := info.ToVars()
	template := getTemplate(templates.Readme, projectType)
	return replaceVars(template, vars)
}

func GetChangelog(info *ProjectInfo) string {
	if info == nil {
		info = getDefaultProjectInfo()
	}
	vars := info.ToVars()
	return replaceVars(templates.Changelog, vars)
}

func GetContributing() string {
	return templates.Contributing
}

func GetGitignore(projectType string) string {
	common := gitignores.Common
	specific := getTemplate(map[string]string{
		"nodejs": gitignores.NodeJS,
		"go":     gitignores.Go,
	}, projectType)

	if specific != "" {
		return common + "\n\n" + specific
	}
	return common
}

func GetLicense(licenseType, author string) string {
	if licenseType == "" {
		licenseType = "MIT"
	}
	if author == "" {
		author = "Your Name"
	}

	license, ok := (*licenses)[licenseType]
	if !ok {
		license = (*licenses)["MIT"]
	}

	vars := getCurrentVars()
	vars["fullname"] = author

	return replaceVars(license.Content, vars)
}

func GetEditorconfig(projectType string) string {
	return getTemplate(qualityTools.Editorconfig, projectType)
}

func GetDockerfile(info *ProjectInfo, projectType string) string {
	if info == nil {
		info = getDefaultProjectInfo()
	}
	vars := info.ToVars()

	var template string
	switch projectType {
	case "nodejs":
		template = devops.Docker.NodeJS.Dockerfile
	case "go":
		template = devops.Docker.Go.Dockerfile
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
	default:
		return ""
	}
}

func GetSecurity(info *ProjectInfo) string {
	if info == nil {
		info = getDefaultProjectInfo()
	}
	vars := info.ToVars()
	return replaceVars(docsTemplates.Security, vars)
}

func GetCodeOfConduct(info *ProjectInfo) string {
	if info == nil {
		info = getDefaultProjectInfo()
	}
	vars := info.ToVars()
	return replaceVars(docsTemplates.CodeOfConduct, vars)
}

func GetMessage(category, key string) string {
	switch category {
	case "check":
		if msg, ok := messages.Check[key]; ok {
			return msg
		}
	case "fix":
		if msg, ok := messages.Fix[key]; ok {
			return msg
		}
	case "errors":
		if msg, ok := messages.Errors[key]; ok {
			return msg
		}
	case "help":
		if msg, ok := messages.Help[key]; ok {
			return msg
		}
	}
	return ""
}

func GetPullRequestTemplate(info *ProjectInfo) string {
	if info == nil {
		info = getDefaultProjectInfo()
	}
	vars := info.ToVars()
	return replaceVars(docsTemplates.PullRequestTemplate, vars)
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
	if info == nil {
		info = getDefaultProjectInfo()
	}
	vars := info.ToVars()
	return replaceVars(docsTemplates.IssueTemplatesConfig, vars)
}

func GetCodeowners(info *ProjectInfo) string {
	if info == nil {
		info = getDefaultProjectInfo()
	}
	vars := info.ToVars()
	return replaceVars(docsTemplates.Codeowners, vars)
}

func GetADRFirst(info *ProjectInfo) string {
	if info == nil {
		info = getDefaultProjectInfo()
	}
	vars := info.ToVars()
	if t, ok := templates.ADR["first"]; ok {
		return replaceVars(t, vars)
	}
	return ""
}

func GetADRTemplate(info *ProjectInfo) string {
	if info == nil {
		info = getDefaultProjectInfo()
	}
	vars := info.ToVars()
	return replaceVars(docsTemplates.ADRTemplate, vars)
}

func GetPreCommit(projectType string) string {
	return getTemplate(qualityTools.PreCommit, projectType)
}

func GetDockerCompose(info *ProjectInfo, projectType string) string {
	if info == nil {
		info = getDefaultProjectInfo()
	}
	vars := info.ToVars()

	template := ""
	if t, ok := devops.DockerCompose["with_db"]; ok {
		template = t
	} else if t, ok := devops.DockerCompose["basic"]; ok {
		template = t
	}

	return replaceVars(template, vars)
}

func GetGitHubAction(projectType string) string {
	switch projectType {
	case "nodejs":
		if t, ok := devops.GitHubActions["nodejs_ci"]; ok {
			return t
		}
	case "go":
		if t, ok := devops.GitHubActions["go_ci"]; ok {
			return t
		}
	}

	if t, ok := devops.GitHubActions["nodejs_ci"]; ok {
		return t
	}

	return ""
}

func GetScripts(info *ProjectInfo, projectType string) map[string]string {
	if info == nil {
		info = getDefaultProjectInfo()
	}

	vars := info.ToVars()
	result := make(map[string]string)

	// Install script
	if s := getScriptTemplate(scripts.Install, projectType); s != "" {
		result["install.sh"] = replaceVars(s, vars)
	}

	// Setup scripts
	if setupScripts, ok := scripts.Setup[projectType]; ok {
		result["setup.sh"] = replaceVars(setupScripts, vars)
	}

	// Test script
	if s := getScriptTemplate(scripts.Test, projectType); s != "" {
		result["test.sh"] = replaceVars(s, vars)
	}

	// Build script
	if s := getScriptTemplate(scripts.Build, projectType); s != "" {
		result["build.sh"] = replaceVars(s, vars)
	}

	// Deploy script
	if s := getScriptTemplate(scripts.Deploy, projectType); s != "" {
		result["deploy.sh"] = replaceVars(s, vars)
	}

	// Release script
	if s := getScriptTemplate(scripts.Release, projectType); s != "" {
		result["release.sh"] = replaceVars(s, vars)
	}

	// Dev script
	if s := getScriptTemplate(scripts.Dev, projectType); s != "" {
		result["dev.sh"] = replaceVars(s, vars)
	}

	// Lint script
	if s := getScriptTemplate(scripts.Lint, projectType); s != "" {
		result["lint.sh"] = replaceVars(s, vars)
	}

	// Format script
	if s := getScriptTemplate(scripts.Format, projectType); s != "" {
		result["format.sh"] = replaceVars(s, vars)
	}

	// Clean script
	if s := getScriptTemplate(scripts.Clean, projectType); s != "" {
		result["clean.sh"] = replaceVars(s, vars)
	}

	// Docker build script
	if s := getScriptTemplate(scripts.DockerBuild, projectType); s != "" {
		result["docker-build.sh"] = replaceVars(s, vars)
	}

	return result
}
