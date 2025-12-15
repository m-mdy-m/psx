package shared

import (
	"fmt"

	"github.com/m-mdy-m/psx/internal/resources"
)

// Exit codes
const (
	ExitSuccess = 0
	ExitFailed  = 1
	ExitConfig  = 2
	ExitFS      = 3
	ExitArgs    = 4
)

// Exit status messages
func ExitMessage(code int) string {
	codeMap := map[int]string{
		ExitSuccess: "success",
		ExitFailed:  "failed",
		ExitConfig:  "config",
		ExitFS:      "fs",
		ExitArgs:    "args",
	}
	if key, ok := codeMap[code]; ok {
		return resources.ExitMessage(key)
	}
	return "Unknown error"
}

// Error messages
const (
	ErrConfigNotFound   = "config_not_found"
	ErrInvalidConfig    = "invalid_config"
	ErrNoProject        = "no_project"
	ErrPermissionDenied = "permission_denied"
	ErrUnknownRule      = "unknown_rule"
	ErrFixFailed        = "fix_failed"
	ErrDetectionFailed  = "detection_failed"
)

func ErrorMessage(errType string) string {
	return resources.ErrorMessage(errType)
}

// Help messages
func Help() {
	fmt.Print(resources.HelpMain())
}

func Version(version string) {
	fmt.Printf("PSX version %s\n", version)
}

func VerboseVersion(version string) {
	fmt.Printf("PSX version %s\nLicense: MIT\nRepository: https://github.com/m-mdy-m/psx\n", version)
}

// Success messages
func CheckSuccess(passed, total int) string {
	return resources.CheckSuccess(passed, total)
}

func FixSuccess(fixed int) string {
	return resources.FixSuccess(fixed)
}

func InitSuccess(path string) string {
	return resources.InitSuccess(path)
}

// Verbose messages
func VerboseCheckStart(path string) string {
	return resources.CheckStart(path)
}

func VerboseDetected(projectType string) string {
	return resources.VerboseDetected(projectType)
}

func VerboseRulesLoaded(count int) string {
	return resources.VerboseRulesLoaded(count)
}

func VerboseRuleCheck(ruleID string) string {
	return resources.CheckRule(ruleID)
}

func VerboseFixApplied(ruleID string) string {
	return resources.FixApplied(ruleID)
}

func VerboseConfigLoaded(path string) string {
	return resources.VerboseConfigLoaded(path)
}

// Fix messages
func FixPrompt(count int) string {
	return resources.FixPrompt(count)
}

func FixDryRunHeader() string {
	return resources.FixDryRun()
}

func FixInteractiveHeader() string {
	return resources.FixInteractive()
}