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
	"github.com/m-mdy-m/psx/internal/shared"
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

	// Filter and group results
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

	// Sort by rule ID
	sort.Slice(errors, func(i, j int) bool { return errors[i].RuleID < errors[j].RuleID })
	sort.Slice(warnings, func(i, j int) bool { return warnings[i].RuleID < warnings[j].RuleID })
	sort.Slice(infos, func(i, j int) bool { return infos[i].RuleID < infos[j].RuleID })

	// Header
	if !f.GlobalFlags.Quiet {
		if f.GlobalFlags.Verbose {
			fmt.Printf("Project: %s\n", r.result.Context.ProjectPath)
			fmt.Printf("Type: %s\n", r.result.Context.ProjectType)
			fmt.Printf("Rules: %d\n", r.result.Summary.Total)
			fmt.Println()
		}
	}

	// Print sections
	if len(errors) > 0 {
		printSection("ERRORS", errors, config.SeverityError)
	}

	if len(warnings) > 0 {
		printSection("WARNINGS", warnings, config.SeverityWarning)
	}

	if len(infos) > 0 && f.GlobalFlags.Verbose {
		printSection("INFO", infos, config.SeverityInfo)
	}

	// Summary
	if !f.GlobalFlags.Quiet {
		fmt.Println()
		r.printSummary()
	}

	return nil
}

func (r *Reporter) reportJSON() error {
	output := map[string]any{
		"status": r.result.Status,
		"summary": map[string]int{
			"total":    r.result.Summary.Total,
			"passed":   r.result.Summary.Passed,
			"errors":   r.result.Summary.Errors,
			"warnings": r.result.Summary.Warnings,
			"info":     r.result.Summary.Info,
		},
		"context": map[string]string{
			"project_path": r.result.Context.ProjectPath,
			"project_type": r.result.Context.ProjectType,
		},
		"results": r.formatResultsForJSON(),
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	fmt.Println(string(data))
	return nil
}

func (r *Reporter) formatResultsForJSON() []map[string]any {
	results := []map[string]any{}

	for _, result := range r.result.Results {
		results = append(results, map[string]any{
			"rule_id":  result.RuleID,
			"passed":   result.Passed,
			"severity": result.Severity,
			"message":  result.Message,
			"fix_hint": result.FixHint,
			"doc_url":  result.DocURL,
		})
	}

	return results
}

func (r *Reporter) printSummary() {
	f := flags.GetFlags()

	total := r.result.Summary.Total
	passed := r.result.Summary.Passed
	errors := r.result.Summary.Errors
	warnings := r.result.Summary.Warnings

	// Compact summary
	if !f.GlobalFlags.Verbose {
		if errors > 0 || warnings > 0 {
			fmt.Printf("Result: %d errors, %d warnings\n", errors, warnings)
		} else {
			fmt.Println(shared.CheckSuccess(passed, total))
		}
	} else {
		// Detailed summary
		fmt.Printf("Checked: %d rules\n", total)
		fmt.Printf("Passed:  %d\n", passed)
		if errors > 0 {
			fmt.Printf("Errors:  %d\n", errors)
		}
		if warnings > 0 {
			fmt.Printf("Warnings: %d\n", warnings)
		}
	}

	// Status
	fmt.Print("Status: ")
	r.printStatus()
}

func (r *Reporter) printStatus() {
	f := flags.GetFlags()

	var statusText string
	var statusColor *color.Color

	switch r.result.Status {
	case checker.StatusPassed:
		statusText = "PASSED"
		statusColor = color.New(color.FgGreen)
	case checker.StatusWarnings:
		statusText = "PASSED (with warnings)"
		statusColor = color.New(color.FgYellow)
	case checker.StatusFailed:
		statusText = "FAILED"
		statusColor = color.New(color.FgRed)
	}

	if f.GlobalFlags.NoColor {
		fmt.Println(statusText)
	} else {
		statusColor.Println(statusText)
	}
}

func printSection(title string, results []checker.RuleResult, severity config.Severity) {
	f := flags.GetFlags()

	// Section header
	icon := getSeverityIcon(severity)

	if f.GlobalFlags.NoColor {
		fmt.Printf("%s %s (%d)\n", icon, title, len(results))
	} else {
		color := getSeverityColor(severity)
		color.Printf("%s %s (%d)\n", icon, title, len(results))
	}

	fmt.Println()

	// Results
	for _, result := range results {
		fmt.Printf("  %s %s\n", icon, result.RuleID)

		if !f.GlobalFlags.Quiet {
			fmt.Printf("      %s\n", result.Message)
		}

		if f.GlobalFlags.Verbose {
			if result.FixHint != "" {
				fmt.Printf("      Fix: %s\n", result.FixHint)
			}
			if result.DocURL != "" {
				fmt.Printf("      Docs: %s\n", result.DocURL)
			}
		}
	}

	fmt.Println()
}

func getSeverityIcon(severity config.Severity) string {
	icons := map[config.Severity]string{
		config.SeverityError:   "✗",
		config.SeverityWarning: "⚠",
		config.SeverityInfo:    "ℹ",
	}
	return icons[severity]
}

func getSeverityColor(severity config.Severity) *color.Color {
	colors := map[config.Severity]*color.Color{
		config.SeverityError:   color.New(color.FgRed),
		config.SeverityWarning: color.New(color.FgYellow),
		config.SeverityInfo:    color.New(color.FgCyan),
	}
	return colors[severity]
}

