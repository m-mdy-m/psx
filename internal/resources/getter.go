package resources

import "strings"

func GetMessage(category, key string) string {
	var table map[string]string

	switch category {
	case "check":
		table = messages.Check
	case "fix":
		table = messages.Fix
	case "errors":
		table = messages.Errors
	case "init":
		table = messages.Init
	case "detect":
		table = messages.Detect
	case "verbose":
		table = messages.Verbose
	}
	if table == nil {
		return ""
	}
	return table[key]
}

func FormatMessage(category, key string, args ...any) string {
	return Message(category, key, args...)
}

// GetLicense returns the text of a named license with the author substituted.
func GetLicense(licenseType, author string) string {
	if licenseType == "" {
		licenseType = "MIT"
	}
	if author == "" {
		author = "Your Name"
	}

	tmpl, ok := (*licenses)[licenseType]
	if !ok {
		tmpl = (*licenses)["MIT"]
	}

	vars := getCurrentVars()
	vars["fullname"] = author
	vars["author"] = author
	return replaceVars(tmpl.Content, vars)
}

func GitignoreCommonBlock() string { return gitignores.Common }

func GitignoreLanguageBlock(lang string) string {
	switch normalizeKey(lang) {
	case "nodejs":
		return gitignores.NodeJS
	case "go":
		return gitignores.Go
	default:
		return ""
	}
}

func LicenseNames() []string {
	out := make([]string, 0, len(*licenses))
	for name := range *licenses {
		out = append(out, name)
	}
	sortStrings(out)
	return out
}

// sortStrings sorts in place; kept local to avoid importing sort in many files.
func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && strings.Compare(s[j-1], s[j]) > 0; j-- {
			s[j-1], s[j] = s[j], s[j-1]
		}
	}
}
