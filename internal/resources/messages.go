package resources

import "fmt"

// Exit messages
func ExitMessage(code string) string {
	if msg, ok := messages.Exit[code]; ok {
		return msg
	}
	return "Unknown error"
}

// Error messages
func ErrorMessage(errType string) string {
	if msg, ok := messages.Errors[errType]; ok {
		return msg
	}
	return "An error occurred"
}

// Help message
func HelpMain() string {
	return messages.Help["main"]
}

// Check messages
func CheckSuccess(passed, total int) string {
	if passed == total {
		return fmt.Sprintf(messages.Check["success_all"], total)
	}
	return fmt.Sprintf(messages.Check["success_partial"], passed, total)
}

func CheckStart(path string) string {
	return fmt.Sprintf(messages.Check["start"], path)
}

func CheckRule(ruleID string) string {
	return fmt.Sprintf(messages.Check["rule_check"], ruleID)
}

// Fix messages
func FixSuccess(fixed int) string {
	if fixed == 0 {
		return messages.Fix["success_none"]
	}
	if fixed == 1 {
		return messages.Fix["success_one"]
	}
	return fmt.Sprintf(messages.Fix["success_many"], fixed)
}

func FixPrompt(count int) string {
	if count == 0 {
		return messages.Fix["prompt_none"]
	}
	if count == 1 {
		return messages.Fix["prompt_one"]
	}
	return fmt.Sprintf(messages.Fix["prompt_many"], count)
}

func FixDryRun() string {
	return messages.Fix["dry_run"]
}

func FixInteractive() string {
	return messages.Fix["interactive"]
}

func FixApplied(ruleID string) string {
	return fmt.Sprintf(messages.Fix["applied"], ruleID)
}

// Init messages
func InitSuccess(path string) string {
	return fmt.Sprintf(messages.Init["success"], path)
}

// Verbose messages
func VerboseDetected(lang string) string {
	return fmt.Sprintf(messages.Verbose["detected"], lang)
}

func VerboseRulesLoaded(count int) string {
	return fmt.Sprintf(messages.Verbose["rules_loaded"], count)
}

func VerboseConfigLoaded(path string) string {
	return fmt.Sprintf(messages.Verbose["config_loaded"], path)
}