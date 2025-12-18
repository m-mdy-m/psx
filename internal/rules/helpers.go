package rules

import "github.com/m-mdy-m/psx/internal/config"

func fileRule(id, category string, severity config.Severity, patterns map[string][]string, fix FixSpec) Rule {
	return Rule{
		ID:       id,
		Type:     RuleTypeFile,
		Category: category,
		Severity: severity,
		Patterns: patterns,
		FixSpec:  fix,
	}
}

func folderRule(id, category string, severity config.Severity, patterns map[string][]string, fix FixSpec) Rule {
	return Rule{
		ID:       id,
		Type:     RuleTypeFolder,
		Category: category,
		Severity: severity,
		Patterns: patterns,
		FixSpec:  fix,
	}
}

func multiRule(id, category string, severity config.Severity, patterns map[string][]string, fix FixSpec) Rule {
	return Rule{
		ID:       id,
		Type:     RuleTypeMulti,
		Category: category,
		Severity: severity,
		Patterns: patterns,
		FixSpec:  fix,
	}
}

func (c *Checker) successResult(rule Rule, foundPattern string) CheckResult {
	return CheckResult{
		RuleID:   rule.ID,
		Passed:   true,
		Severity: rule.Severity,
		Message:  "Found: " + foundPattern,
	}
}

func (c *Checker) failResult(rule Rule, message string) CheckResult {
	return CheckResult{
		RuleID:   rule.ID,
		Passed:   false,
		Severity: rule.Severity,
		Message:  message,
		FixHint:  "Run: psx fix --rule " + rule.ID,
	}
}
