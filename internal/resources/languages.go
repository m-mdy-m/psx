package resources

func NormalizeLanguage(name string) string {
	if canonical, ok := languages.Aliases[name]; ok {
		return canonical
	}
	return name
}

func GetSrcFolders(lang string) []string {
	if folders, ok := languages.Folders[lang]; ok {
		return folders.Src
	}
	return []string{"src"}
}

func GetTestsFolders(lang string) []string {
	if folders, ok := languages.Folders[lang]; ok {
		return folders.Tests
	}
	return []string{"tests"}
}

func GetTestPatterns(lang string) []string {
	if patterns, ok := languages.TestPatterns[lang]; ok {
		return patterns
	}
	return []string{"test", "_test"}
}

func GetSrcFolderName(lang string) string {
	folders := GetSrcFolders(lang)
	if len(folders) > 0 {
		return folders[0]
	}
	return "src"
}

func GetTestsFolderName(lang string) string {
	folders := GetTestsFolders(lang)
	if len(folders) > 0 {
		return folders[0]
	}
	return "tests"
}