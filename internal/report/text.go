package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"

	"github.com/m-mdy-m/psx/internal/rules"
)

type tableRenderer struct{ opts Options }

func (r *tableRenderer) Render(w io.Writer, res *rules.ExecutionResult) error {
	if !r.opts.Quiet {
		r.header(w, res)
	}

	secs := sections(res, levelFor(r.opts))
	for _, s := range secs {
		if len(s.results) == 0 {
			continue
		}
		sortBySeverity(s.results)
		if err := r.section(w, s); err != nil {
			return err
		}
	}

	if r.opts.ShowSkipped && res.Summary.Skipped > 0 {
		if err := r.skipped(w, res); err != nil {
			return err
		}
	}

	if !r.opts.Quiet {
		r.summary(w, res)
	}
	return nil
}

func (r *tableRenderer) header(w io.Writer, res *rules.ExecutionResult) {
	if r.opts.Verbose {
		fmt.Fprintf(w, "project  %s\n", res.Context.ProjectPath)
		fmt.Fprintf(w, "type     %s\n", res.Context.ProjectType)
		fmt.Fprintf(w, "rules    %d enabled, %d skipped\n\n",
			res.Summary.Total-res.Summary.Skipped, res.Summary.Skipped)
	}
}

func (r *tableRenderer) section(w io.Writer, s section) error {
	fmt.Fprintf(w, "%s (%d)\n\n", r.paint(s.title, s.title), len(s.results))

	for _, res := range s.results {
		tag := severityLabel(res.Severity)
		fmt.Fprintf(w, "  %s  %s\n",
			r.paint(pad(tag, 5), tag), r.paint(res.RuleID, "id"))

		if !r.opts.Quiet && res.Message != "" {
			fmt.Fprintf(w, "        %s\n", res.Message)
		}
		if res.FixHint != "" {
			fmt.Fprintf(w, "        fix: %s\n", r.dim(res.FixHint))
		}
		if r.opts.Verbose && res.Evidence != "" {
			fmt.Fprintf(w, "        checked: %s\n", r.dim(res.Evidence))
		}
		if r.opts.Verbose && res.DocURL != "" {
			fmt.Fprintf(w, "        docs: %s\n", r.dim(res.DocURL))
		}
	}
	fmt.Fprintln(w)
	return nil
}

func (r *tableRenderer) skipped(w io.Writer, res *rules.ExecutionResult) error {
	var list []rules.RuleResult
	for _, x := range res.Results {
		if x.Status == rules.StatusSkipped {
			list = append(list, x)
		}
	}
	if len(list) == 0 {
		return nil
	}

	fmt.Fprintf(w, "Skipped (%d)\n\n", len(list))
	for _, x := range list {
		fmt.Fprintf(w, "  -  %s  %s\n", pad(x.RuleID, 24), r.dim(x.Message))
	}
	fmt.Fprintln(w)
	return nil
}

func (r *tableRenderer) summary(w io.Writer, res *rules.ExecutionResult) {
	s := visibleSummary(res, levelFor(r.opts))
	fmt.Fprintf(w, "Result: %s\n", r.summaryLine(s))

	if r.opts.Verbose {
		fmt.Fprintf(w, "  checked  %d\n", s.Total)
		fmt.Fprintf(w, "  passed   %d\n", s.Passed)
		fmt.Fprintf(w, "  failed   %d\n", s.Failed)
		if s.Skipped > 0 {
			fmt.Fprintf(w, "  skipped  %d\n", s.Skipped)
		}
	}

	status, style := statusStyle(res.Status)
	fmt.Fprintf(w, "Status: %s\n", r.paint(status, style))
}

// severityCounts lists only the non-zero severities, so a total and its parts
// always agree.
func severityCounts(s rules.Summary) string {
	var parts []string
	if s.Errors > 0 {
		parts = append(parts, fmt.Sprintf("%d errors", s.Errors))
	}
	if s.Warnings > 0 {
		parts = append(parts, fmt.Sprintf("%d warnings", s.Warnings))
	}
	if s.Info > 0 {
		parts = append(parts, fmt.Sprintf("%d info", s.Info))
	}
	return strings.Join(parts, ", ")
}

// summaryLine builds the compact counts sentence.
func (r *tableRenderer) summaryLine(s rules.Summary) string {
	if s.Failed == 0 {
		if s.Skipped > 0 {
			return fmt.Sprintf("%d passed, %d skipped, no problems found",
				s.Passed, s.Skipped)
		}
		return fmt.Sprintf("all %d checks passed", s.Passed)
	}
	return severityCounts(s)
}

// statusStyle maps an outcome to its label and colour.
func statusStyle(o rules.Outcome) (string, string) {
	switch o {
	case rules.OutcomeFailed:
		return "FAILED", "error"
	case rules.OutcomeWarnings:
		return "PASSED (with warnings)", "warn"
	default:
		return "PASSED", "ok"
	}
}

func (r *tableRenderer) paint(text, style string) string {
	if r.opts.NoColor || color.NoColor {
		return text
	}
	c := colorFor(style)
	if c == nil {
		return text
	}
	return c.Sprint(text)
}

func (r *tableRenderer) dim(text string) string {
	if r.opts.NoColor || color.NoColor {
		return text
	}
	return color.New(color.Faint).Sprint(text)
}

func colorFor(style string) *color.Color {
	switch style {
	case "error":
		return color.New(color.FgRed)
	case "warn":
		return color.New(color.FgYellow)
	case "ok":
		return color.New(color.FgGreen)
	case "id":
		return color.New(color.Bold)
	case "Errors":
		return color.New(color.FgRed, color.Bold)
	case "Warnings":
		return color.New(color.FgYellow, color.Bold)
	case "Info":
		return color.New(color.FgCyan, color.Bold)
	default:
		return nil
	}
}

func pad(text string, width int) string {
	for len(text) < width {
		text += " "
	}
	return text
}

type compactRenderer struct{ opts Options }

func (r *compactRenderer) Render(w io.Writer, res *rules.ExecutionResult) error {
	for _, s := range sections(res, levelFor(r.opts)) {
		for _, x := range s.results {
			fmt.Fprintf(w, "%s %s: %s\n", severityLabel(x.Severity), x.RuleID, x.Message)
			if x.FixHint != "" {
				fmt.Fprintf(w, "  → %s\n", x.FixHint)
			}
		}
	}
	if !r.opts.Quiet {
		label, _ := statusStyle(res.Status)
		fmt.Fprintf(w, "%s: %s\n", label, r.compactSummary(visibleSummary(res, levelFor(r.opts))))
	}
	return nil
}

func (r *compactRenderer) compactSummary(s rules.Summary) string {
	if s.Failed == 0 {
		return fmt.Sprintf("%d passed, %d skipped", s.Passed, s.Skipped)
	}
	return fmt.Sprintf("%d failed (%s)", s.Failed, severityCounts(s))
}

type markdownRenderer struct{ opts Options }

func (r *markdownRenderer) Render(w io.Writer, res *rules.ExecutionResult) error {
	fmt.Fprintf(w, "## psx check\n\n")
	fmt.Fprintf(w, "- project: `%s`\n", res.Context.ProjectPath)
	fmt.Fprintf(w, "- type: `%s`\n", res.Context.ProjectType)
	fmt.Fprintf(w, "- status: **%s**\n\n", strings.ToUpper(string(res.Status)))

	secs := sections(res, levelFor(r.opts))
	any := false
	for _, s := range secs {
		if len(s.results) == 0 {
			continue
		}
		any = true
		sortBySeverity(s.results)
		fmt.Fprintf(w, "### %s\n\n", s.title)
		fmt.Fprintln(w, "| Rule | Message | Fix |")
		fmt.Fprintln(w, "| --- | --- | --- |")
		for _, x := range s.results {
			hint := x.FixHint
			if hint == "" {
				hint = "—"
			} else {
				hint = "`" + hint + "`"
			}
			fmt.Fprintf(w, "| `%s` | %s | %s |\n", x.RuleID, escapeCell(x.Message), hint)
		}
		fmt.Fprintln(w)
	}
	if !any {
		fmt.Fprintf(w, "No problems found across %d rules.\n", res.Summary.Passed)
	}
	return nil
}

func escapeCell(s string) string {
	s = strings.ReplaceAll(s, "|", `\|`)
	return strings.ReplaceAll(s, "\n", " ")
}
