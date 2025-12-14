package checker

import (
	"github.com/m-mdy-m/psx/internal/detector"
	"github.com/m-mdy-m/psx/internal/config"
)
type Status string
const (
	StatusPassed   Status = "passed"
	StatusWarnings Status = "warnings"
	StatusFailed   Status = "failed"
)
type Context struct {
	ProjectPath string
	ProjectType string
	Detection   *detector.DetectionResult
}
type Summary struct {
	Total    int
	Passed   int
	Errors   int
	Warnings int
	Info     int
}
type RuleResult struct {
	RuleID      string
	Passed      bool
	Severity    config.Severity
	Message     string
	FixHint     string
	DocURL      string
}
