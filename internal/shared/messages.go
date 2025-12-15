package shared

import "fmt"

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
	messages := map[int]string{
		ExitSuccess: "Success",
		ExitFailed:  "Validation failed",
		ExitConfig:  "Configuration error",
		ExitFS:      "File system error",
		ExitArgs:    "Invalid arguments",
	}
	if msg, ok := messages[code]; ok {
		return msg
	}
	return "Unknown error"
}

// Error messages
const (
	ErrConfigNotFound    = "config_not_found"
	ErrInvalidConfig     = "invalid_config"
	ErrNoProject         = "no_project"
	ErrPermissionDenied  = "permission_denied"
	ErrUnknownRule       = "unknown_rule"
	ErrFixFailed         = "fix_failed"
	ErrDetectionFailed   = "detection_failed"
)

func ErrorMessage(errType string) string {
	messages := map[string]string{
		ErrConfigNotFound:   "Configuration not found. Run 'psx init' to create one",
		ErrInvalidConfig:    "Invalid configuration. Run 'psx config validate' to check",
		ErrNoProject:        "Not in a project directory",
		ErrPermissionDenied: "Permission denied",
		ErrUnknownRule:      "Unknown rule. Run 'psx rules' to list available rules",
		ErrFixFailed:        "Fix operation failed",
		ErrDetectionFailed:  "Could not detect project type",
	}
	if msg, ok := messages[errType]; ok {
		return msg
	}
	return "An error occurred"
}

// Help messages
func Help() {
	fmt.Print(`PSX - Project Structure Checker

Usage: psx [command] [flags]

Commands:
  check       Validate project structure
  fix         Fix structural issues
  init        Create configuration file
  rules       List validation rules
  config      Manage configuration
  version     Show version

Flags:
  -v, --verbose       Show detailed output
  -q, --quiet         Minimal output
      --config PATH   Config file path
      --no-color      Disable colors
  -h, --help          Show help

Examples:
  psx check
  psx fix --interactive
  psx init --template nodejs

Documentation: https://github.com/m-mdy-m/psx
`)
}

func Version(version string) {
	fmt.Printf("PSX version %s\n", version)
}

func VerboseVersion(version string) {
	fmt.Printf(`PSX - Project Structure Checker
Version: %s
License: MIT
Repository: https://github.com/m-mdy-m/psx
`, version)
}

// Success messages
func CheckSuccess(passed, total int) string {
	if passed == total {
		return fmt.Sprintf("All %d checks passed", total)
	}
	return fmt.Sprintf("%d of %d checks passed", passed, total)
}

func FixSuccess(fixed int) string {
	if fixed == 0 {
		return "No fixes needed"
	}
	if fixed == 1 {
		return "Fixed 1 issue"
	}
	return fmt.Sprintf("Fixed %d issues", fixed)
}

func InitSuccess(path string) string {
	return fmt.Sprintf("Created configuration: %s", path)
}

// Verbose messages
func VerboseCheckStart(path string) string {
	return fmt.Sprintf("Checking project: %s", path)
}

func VerboseDetected(projectType string) string {
	return fmt.Sprintf("Detected: %s", projectType)
}

func VerboseRulesLoaded(count int) string {
	return fmt.Sprintf("Loaded %d rules", count)
}

func VerboseRuleCheck(ruleID string) string {
	return fmt.Sprintf("Checking: %s", ruleID)
}

func VerboseFixApplied(ruleID string) string {
	return fmt.Sprintf("Applied fix: %s", ruleID)
}

func VerboseConfigLoaded(path string) string {
	return fmt.Sprintf("Config loaded: %s", path)
}

// Fix messages
func FixPrompt(count int) string {
	if count == 0 {
		return "No issues to fix"
	}
	if count == 1 {
		return "Found 1 fixable issue. Fix it?"
	}
	return fmt.Sprintf("Found %d fixable issues. Fix them?", count)
}

func FixDryRunHeader() string {
	return "Dry run - no changes will be made"
}

func FixInteractiveHeader() string {
	return "Interactive mode - confirm each fix"
}

