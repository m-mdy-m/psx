package rules

import (
	"fmt"
	"sort"
	"strings"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/resources"
)

// ContentGenerator turns declarative fix specs into concrete file bodies.
//
// There is no per-rule switch here: rule metadata names a template, and the
// resource registry resolves it. Adding a rule never requires editing Go.
type ContentGenerator struct {
	gen *resources.Generator
}

func NewContentGenerator(info *resources.ProjectInfo, projectType string) *ContentGenerator {
	return &ContentGenerator{gen: resources.NewGenerator(info, projectType)}
}

// Generate renders the body for a rule's single-file fix spec.
func (cg *ContentGenerator) Generate(rule *config.ActiveRule, opts resources.Options) (string, error) {
	if rule == nil || rule.Metadata.Fix == nil {
		return "", fmt.Errorf("rule %s declares no fix", ruleIDOf(rule))
	}

	spec := rule.Metadata.Fix
	name := spec.Template
	if name == "" {
		return "", fmt.Errorf("rule %s fix has no template", rule.ID)
	}

	if spec.Content != "" {
		return replaceAll(spec.Content, cg.gen.Vars()), nil
	}

	// The gitignore fix merges the common block with the language block.
	if name == resources.TmplGitignore {
		return cg.gitignore(opts)
	}

	return cg.gen.Render(name, opts)
}

// GenerateFiles renders every file declared by a multi-file fix spec.
// The result map is keyed by project-relative path.
func (cg *ContentGenerator) GenerateFiles(rule *config.ActiveRule, opts resources.Options) (map[string]string, error) {
	if rule == nil || rule.Metadata.Fix == nil {
		return nil, fmt.Errorf("rule %s declares no fix", ruleIDOf(rule))
	}

	files := rule.Metadata.Fix.Files
	if len(files) == 0 {
		return nil, nil
	}

	out := make(map[string]string, len(files))
	for _, path := range SortedPaths(files) {
		tmpl := files[path]

		var (
			body string
			err  error
		)
		switch {
		case tmpl == "" && rule.Metadata.Fix.Content != "":
			body = replaceAll(rule.Metadata.Fix.Content, cg.gen.Vars())
		case resources.IsScriptTemplate(tmpl):
			body = cg.gen.Scripts(opts)[strings.TrimPrefix(path, "scripts/")]
		default:
			body, err = cg.gen.Render(tmpl, opts)
		}
		if err != nil {
			return nil, fmt.Errorf("rule %s: template %q: %w", rule.ID, tmpl, err)
		}
		if strings.TrimSpace(body) == "" {
			return nil, fmt.Errorf("rule %s: template %q produced no content", rule.ID, tmpl)
		}
		out[path] = body
	}
	return out, nil
}

func (cg *ContentGenerator) GenerateScripts(opts resources.Options) map[string]string {
	return cg.gen.Scripts(opts)
}

func (cg *ContentGenerator) gitignore(opts resources.Options) (string, error) {
	pt := opts.ProjectType
	if pt == "" {
		pt = cg.gen.ProjectType()
	}

	var langBlock string
	switch pt {
	case "nodejs":
		langBlock = resources.GitignoreLanguageBlock("nodejs")
	case "go":
		langBlock = resources.GitignoreLanguageBlock("go")
	}

	common := resources.GitignoreCommonBlock()
	if langBlock == "" {
		return replaceAll(common, cg.gen.Vars()), nil
	}
	return replaceAll(common+"\n\n"+langBlock, cg.gen.Vars()), nil
}

func (cg *ContentGenerator) Modes(files map[string]string, spec *config.FixSpec) map[string]uint32 {
	out := make(map[string]uint32, len(files))
	for path := range files {
		if spec != nil && spec.Mode != 0 {
			out[path] = spec.Mode
			continue
		}
		out[path] = resources.ModeFromName(path)
	}
	return out
}

func SortedPaths(files map[string]string) []string {
	out := make([]string, 0, len(files))
	for p := range files {
		out = append(out, p)
	}
	sort.Strings(out)
	return out
}

func ruleIDOf(rule *config.ActiveRule) string {
	if rule == nil {
		return "<nil>"
	}
	return rule.ID
}

func replaceAll(body string, vars map[string]string) string {
	out := body
	for k, v := range vars {
		out = strings.ReplaceAll(out, "{{"+k+"}}", v)
	}
	return out
}
