package reporter

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/flags"
	"github.com/m-mdy-m/psx/internal/rules"
	"github.com/m-mdy-m/psx/internal/checker"
)

type Reporter struct {
	format string
	result *rules.ExecutionResult
}

func New(format string, result *rules.ExecutionResult) *Reporter {
	return &Reporter{
		format: format,
		result: result,
	}
}

func (r *Reporter) Report() error {
	switch r.format {
	case "table":
		return r.reportTable()
	case "json":
		return r.reportJSON()
	default:
		return fmt.Errorf("unsupported format: %s", r.format)
	}
}

func (r *Reporter) reportTable() error {
	f := flags.GetFlags()

	errors := []checker.RuleResult{}
	warnings := []checker.RuleResult{}
	infos := []checker.RuleResult{}
	passed := []checker.RuleResult{}

	for _, result := range r.result.Results {
		if result.Passed {
			passed = append(passed, result)
		} else {
			switch result.Severity {
			case config.SeverityError:
				errors = append(errors, result)
			case config.SeverityWarning:
				warnings = append(warnings, result)
			case config.SeverityInfo:
				infos = append(infos, result)
			}
		}
	}

	sort.Slice(errors, func(i, j int) bool { return errors[i].RuleID < errors[j].RuleID })
	sort.Slice(warnings, func(i, j int) bool { return warnings[i].RuleID < warnings[j].RuleID })
	sort.Slice(infos, func(i, j int) bool { return infos[i].RuleID < infos[j].RuleID })

	// Print header
	fmt.Println()
	printSeparator("═")
	fmt.Println("PSX - Validation Results")
	printSeparator("═")
	fmt.Println()

	// Print errors
	if len(errors) > 0 {
		if !f.GlobalFlags.NoColor {
			color.Red("ERRORS (%d)", len(errors))
		} else {
			fmt.Printf("ERRORS (%d)\n", len(errors))
		}
		printSeparator("─")
		for _, result := range errors {
			printResult(result, "error")
		}
		fmt.Println()
	}

	// Print warnings
	if len(warnings) > 0 {
		if !f.GlobalFlags.NoColor {
			color.Yellow("WARNINGS (%d)", len(warnings))
		} else {
			fmt.Printf("WARNINGS (%d)\n", len(warnings))
		}
		printSeparator("─")
		for _, result := range warnings {
			printResult(result, "warning")
		}
		fmt.Println()
	}

	// Print info (only in verbose mode)
	if f.GlobalFlags.Verbose && len(infos) > 0 {
		if !f.GlobalFlags.NoColor {
			color.Cyan("INFO (%d)", len(infos))
		} else {
			fmt.Printf("INFO (%d)\n", len(infos))
		}
		printSeparator("─")
		for _, result := range infos {
			printResult(result, "info")
		}
		fmt.Println()
	}

	// Print passed (only in verbose mode)
	if f.GlobalFlags.Verbose && len(passed) > 0 {
		if !f.GlobalFlags.NoColor {
			color.Green("PASSED (%d)", len(passed))
		} else {
			fmt.Printf("PASSED (%d)\n", len(passed))
		}
		printSeparator("─")
		for _, result := range passed {
			printResult(result, "passed")
		}
		fmt.Println()
	}

	// Print summary
	printSeparator("═")
	printSummary(r.result)
	printSeparator("═")
	fmt.Println()

	return nil
}

// reportJSON outputs results in JSON format
func (r *Reporter) reportJSON() error {
	output := map[string]any{
		"status":  r.result.Status,
		"summary": r.result.Summary,
		"results": r.result.Results,
		"context": map[string]any{
			"project_path": r.result.Context.ProjectPath,
			"project_type": r.result.Context.ProjectType,
		},
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	fmt.Println(string(data))
	return nil
}

// Helper functions

func printSeparator(char string) {
	fmt.Println(strings.Repeat(char, 60))
}

func printResult(result checker.RuleResult, severity string) {
	f := flags.GetFlags()

	// Print icon and rule ID
	icon := getIcon(severity, f.GlobalFlags.NoColor)
	fmt.Printf("%s %s\n", icon, strings.ToUpper(result.RuleID))

	// Print message
	fmt.Printf("  Message: %s\n", result.Message)

	// Print fix hint if available
	if result.FixHint != "" {
		fmt.Printf("  Fix: %s\n", result.FixHint)
	}

	fmt.Println()
}

func printSummary(result *rules.ExecutionResult) {
	f := flags.GetFlags()

	fmt.Printf("Total Rules: %d\n", result.Summary.Total)
	fmt.Printf("Passed: %d\n", result.Summary.Passed)

	if result.Summary.Errors > 0 {
		if !f.GlobalFlags.NoColor {
			color.Red("Errors: %d", result.Summary.Errors)
		} else {
			fmt.Printf("Errors: %d\n", result.Summary.Errors)
		}
	}

	if result.Summary.Warnings > 0 {
		if !f.GlobalFlags.NoColor {
			color.Yellow("Warnings: %d", result.Summary.Warnings)
		} else {
			fmt.Printf("Warnings: %d\n", result.Summary.Warnings)
		}
	}

	if result.Summary.Info > 0 {
		if !f.GlobalFlags.NoColor {
			color.Cyan("Info: %d", result.Summary.Info)
		} else {
			fmt.Printf("Info: %d\n", result.Summary.Info)
		}
	}

	fmt.Println()

	// Print status
	switch result.Status {
	case checker.StatusPassed:
		if !f.GlobalFlags.NoColor {
			color.Green("Status: PASSED ✓")
		} else {
			fmt.Println("Status: PASSED")
		}
	case checker.StatusWarnings:
		if !f.GlobalFlags.NoColor {
			color.Yellow("Status: PASSED WITH WARNINGS ⚠")
		} else {
			fmt.Println("Status: PASSED WITH WARNINGS")
		}
	case checker.StatusFailed:
		if !f.GlobalFlags.NoColor {
			color.Red("Status: FAILED ✗")
		} else {
			fmt.Println("Status: FAILED")
		}
	}
}

func getIcon(severity string, noColor bool) string {
	if noColor {
		switch severity {
		case "error":
			return "✗"
		case "warning":
			return "⚠"
		case "info":
			return "ℹ"
		case "passed":
			return "✓"
		default:
			return "•"
		}
	}

	switch severity {
	case "error":
		return color.RedString("✗")
	case "warning":
		return color.YellowString("⚠")
	case "info":
		return color.CyanString("ℹ")
	case "passed":
		return color.GreenString("✓")
	default:
		return "•"
	}
}
