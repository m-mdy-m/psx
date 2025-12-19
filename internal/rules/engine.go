package rules

import (
	"fmt"
	"sync"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/logger"
)

type Engine struct {
	ctx    *Context
	rules  map[string]*config.ActiveRule
	checks *Checker
	fixes  *Fixer
}

func NewEngine(cfg *config.Config, ctx *Context) *Engine {
	return &Engine{
		ctx:    ctx,
		rules:  cfg.ActiveRules,
		checks: NewChecker(ctx),
		fixes:  NewFixer(ctx),
	}
}

func Execute(cfg *config.Config, ctx *Context) (*ExecutionResult, error) {
	engine := NewEngine(cfg, ctx)
	return engine.Execute()
}

func (e *Engine) Execute() (*ExecutionResult, error) {
	if len(e.rules) == 0 {
		return nil, fmt.Errorf("no active rules configured")
	}

	logger.Verbose(fmt.Sprintf("Executing %d rules...", len(e.rules)))

	// Prepare results slice
	results := make([]RuleResult, 0, len(e.rules))
	resultsChan := make(chan RuleResult, len(e.rules))

	var wg sync.WaitGroup

	// Execute rules in parallel
	for ruleID, activeRule := range e.rules {
		wg.Add(1)
		go func(id string, rule *config.ActiveRule) {
			defer wg.Done()
			result := e.checkRule(id, rule)
			resultsChan <- result
		}(ruleID, activeRule)
	}

	// Wait for all rules to complete
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results
	for result := range resultsChan {
		results = append(results, result)
	}

	// Calculate summary
	summary := e.calculateSummary(results)
	status := e.determineStatus(summary)

	return &ExecutionResult{
		Context: e.ctx,
		Results: results,
		Summary: summary,
		Status:  status,
	}, nil
}

func (e *Engine) checkRule(ruleID string, activeRule *config.ActiveRule) RuleResult {
	logger.Verbose(fmt.Sprintf("Checking: %s", ruleID))

	patterns := config.GetPatterns(activeRule.Metadata.Patterns, e.ctx.ProjectType)
	if len(patterns) == 0 {
		logger.Verbose(fmt.Sprintf("No patterns for %s in %s projects", ruleID, e.ctx.ProjectType))
		return RuleResult{
			RuleID:   ruleID,
			Passed:   true,
			Severity: activeRule.Severity,
			Message:  "Not applicable for this project type",
		}
	}
	passed := e.checks.CheckAny(patterns)

	if passed {
		return RuleResult{
			RuleID:   ruleID,
			Passed:   true,
			Severity: activeRule.Severity,
			Message:  "OK",
		}
	}

	// Failed
	return RuleResult{
		RuleID:   ruleID,
		Passed:   false,
		Severity: activeRule.Severity,
		Message:  activeRule.Metadata.Message,
		FixHint:  activeRule.Metadata.FixHint,
		DocURL:   activeRule.Metadata.DocURL,
	}
}

func (e *Engine) calculateSummary(results []RuleResult) Summary {
	summary := Summary{Total: len(results)}

	for _, result := range results {
		if result.Passed {
			summary.Passed++
		} else {
			switch result.Severity {
			case config.SeverityError:
				summary.Errors++
			case config.SeverityWarning:
				summary.Warnings++
			case config.SeverityInfo:
				summary.Info++
			}
		}
	}

	return summary
}

func (e *Engine) determineStatus(summary Summary) Status {
	if summary.Errors > 0 {
		return StatusFailed
	}
	if summary.Warnings > 0 {
		return StatusWarnings
	}
	return StatusPassed
}
