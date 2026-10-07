package rules

import (
	"fmt"
	"sort"
	"strings"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/tree"
)

type Engine struct {
	ctx    *Context
	rules  map[string]*config.ActiveRule
	checks *Checker
}

func NewEngine(cfg *config.Config, ctx *Context, snap *tree.Snapshot) *Engine {
	return &Engine{
		ctx:    ctx,
		rules:  cfg.ActiveRules,
		checks: NewChecker(snap),
	}
}

// Execute scans the project once and evaluates every active rule.
func Execute(cfg *config.Config, ctx *Context) (*ExecutionResult, error) {
	snap, err := tree.Scan(ctx.ProjectPath, cfg.Ignore)
	if err != nil {
		return nil, fmt.Errorf("failed to scan project: %w", err)
	}
	return ExecuteSnapshot(cfg, ctx, snap)
}

// ExecuteSnapshot evaluates every active rule against an existing snapshot.
// Callers that already hold a snapshot (such as watch mode) reuse it instead of
// walking the tree again.
func ExecuteSnapshot(cfg *config.Config, ctx *Context, snap *tree.Snapshot) (*ExecutionResult, error) {
	return NewEngine(cfg, ctx, snap).Execute()
}

func (e *Engine) Execute() (*ExecutionResult, error) {
	if len(e.rules) == 0 {
		return nil, fmt.Errorf("no active rules configured")
	}

	logger.Verbose(fmt.Sprintf("Executing %d rules", len(e.rules)))

	ids := make([]string, 0, len(e.rules))
	for id := range e.rules {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	results := make([]RuleResult, 0, len(ids))
	for _, id := range ids {
		results = append(results, e.checkRule(id, e.rules[id]))
	}

	summary := summarize(results)
	return &ExecutionResult{
		Context: e.ctx,
		Results: results,
		Summary: summary,
		Status:  outcome(summary),
	}, nil
}

// checkRule evaluates one rule. A rule whose patterns do not resolve for the
// project type is reported as skipped, never as a pass.
func (e *Engine) checkRule(ruleID string, rule *config.ActiveRule) RuleResult {
	logger.Verbose(fmt.Sprintf("Checking: %s", ruleID))

	res := RuleResult{RuleID: ruleID, Severity: rule.Severity}

	patterns := config.ResolvePatterns(rule.Metadata.Patterns, e.ctx.ProjectType)
	if len(patterns) == 0 {
		res.Status = StatusSkipped
		res.Message = fmt.Sprintf("Not applicable to %s projects", e.ctx.ProjectType)
		return res
	}

	evidence, ok := e.checks.CheckAny(patterns)
	if ok {
		res.Status = StatusPassed
		res.Message = "OK"
		res.Evidence = evidence
		return res
	}

	res.Status = StatusFailed
	res.Message = rule.Metadata.Message
	res.FixHint = rule.Metadata.FixHint
	res.DocURL = rule.Metadata.DocURL
	res.Evidence = strings.Join(patterns, ", ")
	return res
}

// Recount recomputes Summary and Status from Results. Anything that rewrites a
// result must call this, or the counts and the verdict go stale.
func Recount(res *ExecutionResult) {
	res.Summary = summarize(res.Results)
	res.Status = outcome(res.Summary)
}

// summarize counts results by status and by severity of the failures.
func summarize(results []RuleResult) Summary {
	s := Summary{Total: len(results)}
	for _, r := range results {
		switch r.Status {
		case StatusPassed:
			s.Passed++
		case StatusSkipped:
			s.Skipped++
		case StatusFailed:
			s.Failed++
			switch r.Severity {
			case config.SeverityError:
				s.Errors++
			case config.SeverityWarning:
				s.Warnings++
			case config.SeverityInfo:
				s.Info++
			}
		}
	}
	return s
}

// outcome derives the run verdict, ignoring skipped rules.
func outcome(s Summary) Outcome {
	switch {
	case s.Errors > 0:
		return OutcomeFailed
	case s.Warnings > 0:
		return OutcomeWarnings
	default:
		return OutcomePassed
	}
}

func isDirPattern(p string) bool { return strings.HasSuffix(p, "/") }

func hasGlobMeta(p string) bool { return strings.ContainsAny(p, "*?[") }
