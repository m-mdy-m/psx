package rules

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/resources"
)

// A fix is expressed entirely by rule metadata, so this type contains no
// per-rule knowledge.
type Fixer struct {
	ctx    *Context
	gen    *ContentGenerator
	opts   resources.Options
	dryRun bool
	force  bool
}

func NewFixer(ctx *Context, fixCtx *FixContext, opts resources.Options) *Fixer {
	dryRun := fixCtx != nil && fixCtx.DryRun
	force := fixCtx != nil && fixCtx.Force
	return &Fixer{
		ctx:    ctx,
		gen:    NewContentGenerator(ctx.ProjectInfo, ctx.ProjectType),
		opts:   opts,
		dryRun: dryRun,
		force:  force,
	}
}

// Fix applies a single rule. An unknown rule id is an error; a rule that
// cannot be fixed is reported through the result.
func Fix(cfg *config.Config, fixCtx *FixContext, ruleID string, opts resources.Options) (*FixResult, error) {
	rule, ok := cfg.ActiveRules[ruleID]
	if !ok {
		return nil, fmt.Errorf("unknown rule: %s", ruleID)
	}
	return NewFixer(fixCtx.Context, fixCtx, opts).fix(ruleID, rule), nil
}

// FixAll applies fixes for the given rules and reports one result per rule.
//
// Failures are recorded in the returned results rather than swallowed, so the
// caller can surface them.
func FixAll(cfg *config.Config, fixCtx *FixContext, ruleIDs []string, opts resources.Options) []*FixResult {
	f := NewFixer(fixCtx.Context, fixCtx, opts)

	out := make([]*FixResult, 0, len(ruleIDs))
	for _, id := range ruleIDs {
		rule, ok := cfg.ActiveRules[id]
		if !ok {
			out = append(out, &FixResult{RuleID: id, Error: fmt.Errorf("unknown rule: %s", id)})
			continue
		}
		out = append(out, f.fix(id, rule))
	}

	// User-declared entries are applied last, after the built-in rules.
	if cfg.Custom != nil {
		handler := CustomHandler{root: fixCtx.Context.ProjectPath}
		out = append(out, handler.Apply(cfg.Custom, fixCtx)...)
	}
	return out
}

func (f *Fixer) fix(ruleID string, rule *config.ActiveRule) *FixResult {
	res := &FixResult{RuleID: ruleID}

	if !rule.Metadata.Fixable() {
		res.Skipped = true
		res.Reason = "no automatic fix available"
		return res
	}

	logger.Verbose(fmt.Sprintf("Fixing: %s", ruleID))

	if files := rule.Metadata.FileFixes(); len(files) > 0 {
		return f.fixFiles(res, rule, files)
	}
	return f.fixOne(res, rule)
}

func (f *Fixer) fixOne(res *FixResult, rule *config.ActiveRule) *FixResult {
	spec := rule.Metadata.Fix
	rel := strings.TrimSuffix(spec.Path, "/")

	if reason, done := f.skipExisting(rel); done {
		res.Skipped = true
		res.Reason = reason
		return res
	}

	isDir := strings.HasSuffix(spec.Path, "/")
	if isDir {
		full := filepath.Join(f.ctx.ProjectPath, filepath.FromSlash(rel))
		res.Changes = []Change{{
			Type:        ChangeCreateFolder,
			Path:        full,
			Description: "create " + rel + "/",
		}}
		if f.dryRun {
			res.Fixed = true
			return res
		}
		if err := os.MkdirAll(full, 0o755); err != nil {
			res.Error = fmt.Errorf("create dir %s: %w", rel, err)
			return res
		}
		res.Fixed = true
		return res
	}

	body, err := f.gen.Generate(rule, f.opts)
	if err != nil {
		// A template that does not exist for this project type is not a
		// failure: the rule simply does not apply here.
		if resources.IsNotApplicable(err) {
			res.Skipped = true
			res.Reason = err.Error()
			return res
		}
		res.Error = err
		return res
	}
	if strings.TrimSpace(body) == "" {
		res.Error = fmt.Errorf("template produced an empty body; refusing to create %s", rel)
		return res
	}

	mode := spec.Mode
	if mode == 0 {
		mode = resources.ModeFromName(rel)
	}
	return f.writeOne(res, rel, body, mode)
}

func (f *Fixer) fixFiles(res *FixResult, rule *config.ActiveRule, files map[string]string) *FixResult {
	bodies, err := f.gen.GenerateFiles(rule, f.opts)
	if err != nil {
		if resources.IsNotApplicable(err) {
			res.Skipped = true
			res.Reason = err.Error()
			return res
		}
		res.Error = err
		return res
	}

	modes := f.gen.Modes(bodies, rule.Metadata.Fix)
	var changes []Change
	var failures []string

	for _, rel := range SortedPaths(bodies) {
		if reason, done := f.skipExisting(rel); done {
			logger.Verbose(fmt.Sprintf("Skipped %s: %s", rel, reason))
			continue
		}
		body := bodies[rel]
		if strings.TrimSpace(body) == "" {
			failures = append(failures, rel+": empty template body")
			continue
		}
		change, err := f.write(rel, body, modes[rel])
		if err != nil {
			failures = append(failures, rel+": "+err.Error())
			continue
		}
		changes = append(changes, change...)
	}

	res.Changes = changes
	if len(failures) > 0 {
		// Partial success is still reported, alongside the failures.
		res.Fixed = len(changes) > 0
		res.Error = errors.New(strings.Join(failures, "; "))
		return res
	}
	if len(changes) == 0 {
		res.Skipped = true
		res.Reason = "all target paths already exist"
		return res
	}
	res.Fixed = true
	return res
}

func (f *Fixer) writeOne(res *FixResult, rel, body string, mode uint32) *FixResult {
	changes, err := f.write(rel, body, mode)
	if err != nil {
		res.Error = err
		return res
	}
	res.Changes = changes
	res.Fixed = true
	return res
}

func (f *Fixer) write(rel, body string, mode uint32) ([]Change, error) {
	full := filepath.Join(f.ctx.ProjectPath, filepath.FromSlash(rel))

	if f.dryRun {
		return []Change{{
			Type:        ChangeCreateFile,
			Path:        full,
			Description: "create " + rel,
			Content:     Preview(body, 10),
			Mode:        mode,
		}}, nil
	}

	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return nil, fmt.Errorf("create parent of %s: %w", rel, err)
	}
	if err := os.WriteFile(full, []byte(body), os.FileMode(mode)); err != nil {
		return nil, fmt.Errorf("write %s: %w", rel, err)
	}
	return []Change{{
		Type:        ChangeCreateFile,
		Path:        full,
		Description: "created " + rel,
		Mode:        mode,
	}}, nil
}

// skipExisting reports whether a path must be left alone.
//
// A non-empty target is never overwritten unless --force is set; an existing
// empty file or directory is filled in.
func (f *Fixer) skipExisting(rel string) (string, bool) {
	full := filepath.Join(f.ctx.ProjectPath, filepath.FromSlash(rel))

	info, err := os.Lstat(full)
	if err != nil {
		return "", false // absent: nothing to preserve
	}
	if f.force {
		return "", false
	}
	if info.IsDir() {
		entries, err := os.ReadDir(full)
		if err == nil && len(entries) > 0 {
			return "directory is not empty", true
		}
		return "", false
	}
	if info.Size() > 0 {
		return "file already has content", true
	}
	return "", false
}

func Preview(content string, maxLines int) string {
	if content == "" {
		return ""
	}
	lines := strings.Split(content, "\n")
	if len(lines) <= maxLines {
		return content
	}
	return fmt.Sprintf("%s\n... (%d more lines)",
		strings.Join(lines[:maxLines], "\n"), len(lines)-maxLines)
}

func sortedRuleIDs(m map[string]*config.ActiveRule) []string {
	out := make([]string, 0, len(m))
	for id := range m {
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}
