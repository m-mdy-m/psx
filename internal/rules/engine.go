package rules

import (
	"fmt"
	"sync"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/detector"
	"github.com/m-mdy-m/psx/internal/shared/logger"
)

// Engine executes rules against a project
type Engine struct {
	config    *config.Config
	detection *detector.DetectionResult
	context   *Context
}

// NewEngine creates a new rule engine
func NewEngine(cfg *config.Config, detection *detector.DetectionResult) *Engine {
	return &Engine{
		config:    cfg,
		detection: detection,
		context: &Context{
			ProjectPath: cfg.Path,
			ProjectType: detection.Type.Primary,
			Detection:   detection,
		},
	}
}

// Execute runs all enabled rules
func (e *Engine) Execute() (*ExecutionResult, error) {
	result := &ExecutionResult{
		Results:  make([]RuleResult, 0, len(e.config.ActiveRules)),
		Summary:  Summary{},
		Context:  e.context,
	}

	// Execute rules in parallel
	var wg sync.WaitGroup
	resultsChan := make(chan RuleResult, len(e.config.ActiveRules))

	for _, activeRule := range e.config.ActiveRules {
		wg.Add(1)
		go func(ar *config.ActiveRule) {
			defer wg.Done()

			logger.Verbose(fmt.Sprintf("  Checking rule: %s", ar.ID))

			// Get the rule executor
			executor, err := GetExecutor(ar.ID)
			if err != nil {
				logger.Warning(fmt.Sprintf("  Rule %s: no executor found, skipping", ar.ID))
				return
			}

			// Execute the rule
			ruleResult := executor.Execute(e.context, &ar.Metadata)
			ruleResult.RuleID = ar.ID
			ruleResult.Severity = ar.Severity

			resultsChan <- ruleResult
		}(activeRule)
	}

	// Wait for all rules to complete
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results
	for ruleResult := range resultsChan {
		result.Results = append(result.Results, ruleResult)

		// Update summary
		if ruleResult.Passed {
			result.Summary.Passed++
		} else {
			switch ruleResult.Severity {
			case config.SeverityError:
				result.Summary.Errors++
			case config.SeverityWarning:
				result.Summary.Warnings++
			case config.SeverityInfo:
				result.Summary.Info++
			}
		}
	}

	// Determine overall status
	result.Summary.Total = len(result.Results)
	result.Status = StatusPassed
	if result.Summary.Errors > 0 {
		result.Status = StatusFailed
	} else if result.Summary.Warnings > 0 {
		result.Status = StatusWarnings
	}

	return result, nil
}
