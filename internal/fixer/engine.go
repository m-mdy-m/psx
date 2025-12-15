package fixer

import (
	"fmt"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/rules"
)

type Engine struct {
	context *FixContext
	fixes   map[string]func(*FixContext) (*FixResult, error)
}

func NewEngine(ctx *FixContext) *Engine {
	engine := &Engine{
		context: ctx,
		fixes:   make(map[string]func(*FixContext) (*FixResult, error)),
	}

	engine.registerFixers()
	return engine
}

func (e *Engine) registerFixers() {
	// General
	e.fixes["readme"] = FixReadme
	e.fixes["license"] = FixLicense
	e.fixes["gitignore"] = FixGitignore
	e.fixes["changelog"] = FixChangelog

	// Structure
	e.fixes["src_folder"] = FixSrcFolder
	e.fixes["tests_folder"] = FixTestsFolder
	e.fixes["docs_folder"] = FixDocsFolder

	// Documentation
	e.fixes["adr"] = FixADR
	e.fixes["contributing"] = FixContributing
	e.fixes["api_docs"] = FixAPIDocsFolder
}

func (e *Engine) CanFix(ruleID string) bool {
	_, ok := e.fixes[ruleID]
	return ok
}

func (e *Engine) Fix(ruleID string) (*FixResult, error) {
	fixer, ok := e.fixes[ruleID]
	if !ok {
		return nil, fmt.Errorf("no fixer available for rule: %s", ruleID)
	}

	logger.Verbose(fmt.Sprintf("Applying fix: %s", ruleID))
	return fixer(e.context)
}

func (e *Engine) FixAll(failedRules []string) (*FixPlan, error) {
	plan := &FixPlan{
		Fixes: []*FixResult{},
	}

	for _, ruleID := range failedRules {
		if !e.CanFix(ruleID) {
			logger.Verbose(fmt.Sprintf("No fixer for: %s", ruleID))
			continue
		}

		result, err := e.Fix(ruleID)
		if err != nil {
			logger.Warning(fmt.Sprintf("Fix failed for %s: %v", ruleID, err))
			continue
		}

		plan.Fixes = append(plan.Fixes, result)
		plan.TotalChanges += len(result.Changes)
	}

	return plan, nil
}

func GetFixableFails(execResult *rules.ExecutionResult, cfg *config.Config) []string {
	fixable := []string{}

	for _, result := range execResult.Results {
		if result.Passed {
			continue
		}

		// Check if active in config
		if _, ok := cfg.ActiveRules[result.RuleID]; !ok {
			continue
		}

		fixable = append(fixable, result.RuleID)
	}

	return fixable
}

func GenerateSummary(plan *FixPlan) FixSummary {
	summary := FixSummary{}

	summary.Total = len(plan.Fixes)
	summary.Changes = plan.TotalChanges

	for _, fix := range plan.Fixes {
		if fix.Fixed {
			summary.Fixed++
		} else if fix.Skipped {
			summary.Skipped++
		} else if fix.Error != nil {
			summary.Failed++
		}
	}

	return summary
}

