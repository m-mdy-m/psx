package resources

import (
	"path/filepath"
)

func GetReadme(projectName, projectType string) string {
	vars := getCurrentVars()
	vars["project_name"] = projectName
	vars["project_type"] = projectType

	template := templates.Readme[projectType]
	if template == "" {
		template = templates.Readme["generic"]
	}

	return replaceVars(template, vars)
}

func GetChangelog() string {
	vars := getCurrentVars()
	return replaceVars(templates.Changelog, vars)
}

func GetContributing() string {
	return templates.Contributing
}

func GetFirstADR() string {
	vars := getCurrentVars()
	return replaceVars(templates.ADR["first"], vars)
}

func GetADRTemplate() string {
	return templates.ADR["template"]
}

func GetAPIDocs(projectType string) string {
	template := templates.APIDocs[projectType]
	if template == "" {
		template = templates.APIDocs["generic"]
	}
	return template
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

func GetEditorconfig() string {
	return templates.Editorconfig
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

func GetProjectName(projectPath string) string {
	return filepath.Base(projectPath)
}