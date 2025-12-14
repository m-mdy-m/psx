package rules

import (
	"fmt"
	"sync"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/detector"
	"github.com/m-mdy-m/psx/internal/checker"
	"github.com/m-mdy-m/psx/internal/logger"
)

type Engine struct {
	config    *config.Config
	detection *detector.DetectionResult
	context   *checker.Context
}

func NewEngine(cfg *config.Config, detection *detector.DetectionResult) *Engine {
	return &Engine{
		config:    cfg,
		detection: detection,
		context: &checker.Context{
			ProjectPath: cfg.Path,
			ProjectType: detection.Type.Primary,
			Detection:   detection,
		},
	}
}

func (e *Engine) Execute() (*ExecutionResult, error) {
	result := &ExecutionResult{
		Results:  make([]checker.RuleResult, 0, len(e.config.ActiveRules)),
		Summary:  checker.Summary{},
		Context:  e.context,
	}

	var wg sync.WaitGroup
	resultsChan := make(chan checker.RuleResult, len(e.config.ActiveRules))

	for _, activeRule := range e.config.ActiveRules {
		wg.Add(1)
		go func(ar *config.ActiveRule) {
			defer wg.Done()

			logger.Verbose(fmt.Sprintf("  Checking rule: %s", ar.ID))

			executor, err := GetExecutor(ar.ID)
			if err != nil {
				logger.Warning(fmt.Sprintf("  Rule %s: no executor found, skipping", ar.ID))
				return
			}

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
	result.Status = checker.StatusPassed
	if result.Summary.Errors > 0 {
		result.Status = checker.StatusFailed
	} else if result.Summary.Warnings > 0 {
		result.Status = checker.StatusWarnings
	}

	return result, nil
}
