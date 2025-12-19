package rules

import (
	"path/filepath"
	"strings"
)

type PatternType int

const (
	PatternTypeFile PatternType = iota
	PatternTypeFolder
	PatternTypeMultiple
)

type PatternResolver struct {
	knownFolders    map[string]bool
	knownExtensions map[string]bool
}

func NewPatternResolver() *PatternResolver {
	return &PatternResolver{
		knownFolders: map[string]bool{
			"src":            true,
			"tests":          true,
			"test":           true,
			"__tests__":      true,
			"docs":           true,
			"documentation":  true,
			"scripts":        true,
			"script":         true,
			"bin":            true,
			"build":          true,
			"dist":           true,
			"adr":            true,
			"ISSUE_TEMPLATE": true,
			"workflows":      true,
			"api":            true,
			"infra":          true,
			"k8s":            true,
			"kubernetes":     true,
			"cmd":            true,
			"internal":       true,
			"pkg":            true,
			".github":        true,
			".husky":         true,
		},
		knownExtensions: map[string]bool{
			".md":         true,
			".txt":        true,
			".yml":        true,
			".yaml":       true,
			".json":       true,
			".js":         true,
			".ts":         true,
			".go":         true,
			".py":         true,
			".sh":         true,
			".bash":       true,
			".ps1":        true,
			".Dockerfile": true,
		},
	}
}

func (pr *PatternResolver) ResolveType(pattern string) PatternType {
	if strings.HasSuffix(pattern, "/") {
		return PatternTypeFolder
	}

	base := filepath.Base(pattern)

	if pr.knownFolders[base] {
		return PatternTypeFolder
	}

	ext := filepath.Ext(base)
	if ext != "" {
		return PatternTypeFile
	}

	switch base {
	case "Dockerfile", "Makefile", "LICENSE", "CHANGELOG", "CONTRIBUTING",
		"README", "CODEOWNERS", ".gitignore", ".dockerignore", ".editorconfig",
		".env.example", ".env.sample":
		return PatternTypeFile
	}

	if strings.HasPrefix(base, ".") && len(base) > 1 {
		if strings.Count(base, ".") == 1 {
			return PatternTypeFile
		}
		return PatternTypeFile
	}

	return PatternTypeFolder
}

func (pr *PatternResolver) NeedsContent(pattern string) bool {
	patternType := pr.ResolveType(pattern)
	return patternType == PatternTypeFile
}

type CreateableResource struct {
	Path    string
	Type    PatternType
	RuleID  string
	Content string
}

func (pr *PatternResolver) ResolvePattern(ruleID, pattern string, projectPath string) CreateableResource {
	patternType := pr.ResolveType(pattern)
	fullPath := filepath.Join(projectPath, pattern)

	return CreateableResource{
		Path:   fullPath,
		Type:   patternType,
		RuleID: ruleID,
	}
}

func (pr *PatternResolver) IsSpecialMultiFileRule(ruleID string) bool {
	multiFileRules := []string{
		"issue_templates",
		"adr",
		"scripts_folder",
	}

	for _, rule := range multiFileRules {
		if ruleID == rule {
			return true
		}
	}

	return false
}
