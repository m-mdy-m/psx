package shared

import (
	"fmt"
)

func Help() {
	fmt.Println(`
╔═══════════════════════════════════════════════════════════════╗
║                PSX - Project Structure Checker                ║
║                    Version: development                        ║
╚═══════════════════════════════════════════════════════════════╝

PSX validates and standardizes project structures across different
programming languages and frameworks.

USAGE:
  psx [command] [flags]

CORE COMMANDS:
  check       Validate project structure
  fix         Automatically fix structural issues
  init        Initialize PSX configuration
  rules       List available validation rules
  config      Manage PSX configuration
  version     Show version information

EXAMPLES:
  # Basic usage
  psx check                    Check current directory
  psx fix --interactive        Fix issues with confirmation
  psx init                     Create configuration file

  # Advanced usage
  psx check --verbose          Detailed validation output
  psx check --output json      JSON output for CI/CD
  psx fix --dry-run            Preview fixes without applying
  psx rules --category general Show only general rules

GLOBAL FLAGS:
  -v, --verbose                Show detailed information
  -q, --quiet                  Minimal output
      --config <file>          Use specific config file
      --no-color               Disable colored output
  -h, --help                   Show help

Run 'psx [command] --help' for more information about a command.

DOCUMENTATION:
  Repository:  https://github.com/m-mdy-m/psx
  Issues:      https://github.com/m-mdy-m/psx/issues
  Email:       bitsgenix@gmail.com
`)
}

func ShortHelp() {
	fmt.Println(`
PSX - Project Structure Checker

Usage: psx [command] [flags]

Commands:
  check       Validate project structure
  fix         Fix structural issues
  init        Initialize configuration
  rules       List validation rules
  config      Manage configuration
  version     Show version

Run 'psx --help' for detailed information.
`)
}

func WelcomeBanner() {
	fmt.Println(`
╔═══════════════════════════════════════════════════════════════╗
║                PSX - Project Structure Checker                ║
╚═══════════════════════════════════════════════════════════════╝
`)
}

func GoodbyeMessage() {
	fmt.Println(`
Thank you for using PSX! 👋
Report issues: https://github.com/m-mdy-m/psx/issues
`)
}

func QuickStart() {
	fmt.Println(`
QUICK START:

1. Initialize Configuration:
   psx init

2. Check Your Project:
   psx check --verbose

3. Fix Issues:
   psx fix --interactive

4. Customize Rules:
   Edit psx.yml in your project root

Need help? Run 'psx --help' or visit:
https://github.com/m-mdy-m/psx
`)
}

func CommonErrors() map[string]string {
	return map[string]string{
		"config_not_found":   "No configuration file found. Run 'psx init' to create one.",
		"invalid_config":     "Configuration file is invalid. Run 'psx config validate' to check.",
		"no_project":         "Not in a project directory. Specify path with 'psx check ./path'",
		"permission_denied":  "Permission denied. Try running with appropriate permissions.",
		"unknown_rule":       "Unknown rule. Run 'psx rules' to see available rules.",
		"fix_failed":         "Fix operation failed. Review the error and try manual fixing.",
	}
}

func GetErrorHelp(errorType string) string {
	errors := CommonErrors()
	if help, exists := errors[errorType]; exists {
		return help
	}
	return "Run 'psx --help' for usage information."
}

func ExitMessages() map[int]string {
	return map[int]string{
		0: "Success",
		1: "Validation failed or errors found",
		2: "Configuration error",
		3: "File system error",
		4: "Invalid arguments",
	}
}

