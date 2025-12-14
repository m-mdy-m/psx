package reporter

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/fatih/color"

	"github.com/m-mdy-m/psx/internal/checker"
	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/flags"
	"github.com/m-mdy-m/psx/internal/rules"
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

	// Group by severity
	errors := []checker.RuleResult{}
	warnings := []checker.RuleResult{}
	infos := []checker.RuleResult{}

	for _, result := range r.result.Results {
		if result.Passed {
			continue
		}
		switch result.Severity {
		case config.SeverityError:
			errors = append(errors, result)
		case config.SeverityWarning:
			warnings = append(warnings, result)
		case config.SeverityInfo:
			infos = append(infos, result)
		}
	}

	// Sort
	sort.Slice(errors, func(i, j int) bool { return errors[i].RuleID < errors[j].RuleID })
	sort.Slice(warnings, func(i, j int) bool { return warnings[i].RuleID < warnings[j].RuleID })
	sort.Slice(infos, func(i, j int) bool { return infos[i].RuleID < infos[j].RuleID })

	// Print results
	if !f.GlobalFlags.Quiet {
		fmt.Println()

		// Verbose header
		if f.GlobalFlags.Verbose {
			fmt.Printf("Project: %s (%s)\n", r.result.Context.ProjectPath, r.result.Context.ProjectType)
			fmt.Printf("Rules: %d checked\n", r.result.Summary.Total)
			fmt.Println()
		}
	}

	// Errors
	if len(errors) > 0 {
		printSection("ERRORS", errors, "red")
	}

	// Warnings
	if len(warnings) > 0 {
		printSection("WARNINGS", warnings, "yellow")
	}

	// Info (only verbose)
	if len(infos) > 0 && f.GlobalFlags.Verbose {
		printSection("INFO", infos, "cyan")
	}

	// Summary
	if !f.GlobalFlags.Quiet {
		fmt.Println()
		printSummary(r.result)
	}

	return nil
}

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

func printSection(title string, results []checker.RuleResult, colorName string) {
	f := flags.GetFlags()

	// Header
	if !f.GlobalFlags.NoColor {
		switch colorName {
		case "red":
			color.Red("%s (%d)", title, len(results))
		case "yellow":
			color.Yellow("%s (%d)", title, len(results))
		case "cyan":
			color.Cyan("%s (%d)", title, len(results))
		}
	} else {
		fmt.Printf("%s (%d)\n", title, len(results))
	}

	// Results
	for _, result := range results {
		icon := getIcon(colorName)
		fmt.Printf("  %s %s\n", icon, result.RuleID)

		if f.GlobalFlags.Verbose {
			fmt.Printf("      %s\n", result.Message)
			if result.FixHint != "" {
				fmt.Printf("      Fix: %s\n", result.FixHint)
			}
		}
	}
	fmt.Println()
}

func printSummary(result *rules.ExecutionResult) {
	f := flags.GetFlags()

	total := result.Summary.Total
	passed := result.Summary.Passed
	errors := result.Summary.Errors
	warnings := result.Summary.Warnings

	if f.GlobalFlags.Verbose {
		fmt.Printf("Results: %d passed, %d errors, %d warnings (total: %d)\n",
			passed, errors, warnings, total)
	} else {
		if errors > 0 {
			fmt.Printf("%d errors, %d warnings\n", errors, warnings)
		} else if warnings > 0 {
			fmt.Printf("%d warnings\n", warnings)
		}
	}

	// Status
	fmt.Print("Status: ")
	switch result.Status {
	case checker.StatusPassed:
		if !f.GlobalFlags.NoColor {
			color.Green("PASSED")
		} else {
			fmt.Println("PASSED")
		}
	case checker.StatusWarnings:
		if !f.GlobalFlags.NoColor {
			color.Yellow("PASSED (with warnings)")
		} else {
			fmt.Println("PASSED (with warnings)")
		}
	case checker.StatusFailed:
		if !f.GlobalFlags.NoColor {
			color.Red("FAILED")
		} else {
			fmt.Println("FAILED")
		}
	}
}

func getIcon(colorName string) string {
	f := flags.GetFlags()

	icons := map[string]string{
		"red":    "✗",
		"yellow": "⚠",
		"cyan":   "ℹ",
		"green":  "✓",
	}

	icon := icons[colorName]
	if f.GlobalFlags.NoColor {
		return icon
	}

	switch colorName {
	case "red":
		return color.RedString(icon)
	case "yellow":
		return color.YellowString(icon)
	case "cyan":
		return color.CyanString(icon)
	case "green":
		return color.GreenString(icon)
	default:
		return icon
	}
}

