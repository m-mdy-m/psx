// Package report renders an execution result in the requested format.
//
// Every renderer writes to the single writer it is given, so a caller decides
// whether output goes to stdout or a file. Renderers never print progress
// messages of their own.
package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/rules"
)

// Options controls what a renderer includes.
type Options struct {
	Format      string
	Level       string
	ShowPassed  bool
	ShowSkipped bool
	Verbose     bool
	Quiet       bool
	NoColor     bool
	ProjectPath string
	ProjectType string
	Version     string
}

// Renderer is implemented by each output format.
type Renderer interface {
	Render(w io.Writer, res *rules.ExecutionResult) error
}

func New(opts Options) (Renderer, error) {
	switch strings.ToLower(opts.Format) {
	case "", "table":
		return &tableRenderer{opts: opts}, nil
	case "compact":
		return &compactRenderer{opts: opts}, nil
	case "json":
		return &jsonRenderer{opts: opts}, nil
	case "ndjson":
		return &ndjsonRenderer{opts: opts}, nil
	case "sarif":
		return &sarifRenderer{opts: opts}, nil
	case "github":
		return &githubRenderer{opts: opts}, nil
	case "junit":
		return &junitRenderer{opts: opts}, nil
	case "markdown":
		return &markdownRenderer{opts: opts}, nil
	default:
		return nil, fmt.Errorf("unsupported output format %q", opts.Format)
	}
}

type section struct {
	title   string
	results []rules.RuleResult
}

func sections(res *rules.ExecutionResult, level string) []section {
	out := []section{
		{title: "Errors"},
		{title: "Warnings"},
		{title: "Info"},
	}

	for _, r := range res.Results {
		if r.Status != rules.StatusFailed {
			continue
		}
		if !includeLevel(level, r.Severity) {
			continue
		}
		switch r.Severity {
		case config.SeverityError:
			out[0].results = append(out[0].results, r)
		case config.SeverityWarning:
			out[1].results = append(out[1].results, r)
		default:
			out[2].results = append(out[2].results, r)
		}
	}
	return out
}

// visibleSummary recounts failures that the level filter actually displays, so the
// summary never advertises findings the report just hid.
func visibleSummary(res *rules.ExecutionResult, level string) rules.Summary {
	s := res.Summary
	s.Failed, s.Errors, s.Warnings, s.Info = 0, 0, 0, 0

	for _, r := range res.Results {
		if r.Status != rules.StatusFailed || !includeLevel(level, r.Severity) {
			continue
		}
		s.Failed++
		switch r.Severity {
		case config.SeverityError:
			s.Errors++
		case config.SeverityWarning:
			s.Warnings++
		default:
			s.Info++
		}
	}
	return s
}

func includeLevel(level string, sev config.Severity) bool {
	switch level {
	case "", "all":
		return true
	case "error":
		return sev == config.SeverityError
	case "warning":
		return sev == config.SeverityWarning
	case "info":
		return sev == config.SeverityInfo
	default:
		return true
	}
}

func levelFor(opts Options) string {
	if opts.Verbose {
		return "all"
	}
	if opts.Level == "" {
		return "all"
	}
	return opts.Level
}

// severityLabel renders a severity as a fixed-width tag.
func severityLabel(s config.Severity) string {
	switch s {
	case config.SeverityError:
		return "error"
	case config.SeverityWarning:
		return "warn"
	case config.SeverityInfo:
		return "info"
	default:
		return string(s)
	}
}

// jsonDoc is the stable machine-readable shape.
type jsonDoc struct {
	Version string      `json:"version"`
	Status  string      `json:"status"`
	Summary jsonSummary `json:"summary"`
	Context jsonContext `json:"context"`
	Results []jsonRule  `json:"results"`
}

type jsonSummary struct {
	Total    int `json:"total"`
	Passed   int `json:"passed"`
	Failed   int `json:"failed"`
	Skipped  int `json:"skipped"`
	Errors   int `json:"errors"`
	Warnings int `json:"warnings"`
	Info     int `json:"info"`
}

type jsonContext struct {
	ProjectPath string `json:"project_path"`
	ProjectType string `json:"project_type"`
}

type jsonRule struct {
	RuleID   string `json:"rule_id"`
	Category string `json:"category"`
	Status   string `json:"status"`
	Severity string `json:"severity"`
	Passed   bool   `json:"passed"`
	Fixable  bool   `json:"fixable"`
	Message  string `json:"message,omitempty"`
	Evidence string `json:"evidence,omitempty"`
	FixHint  string `json:"fix_hint,omitempty"`
	DocURL   string `json:"doc_url,omitempty"`
}

func toDoc(res *rules.ExecutionResult, version string) jsonDoc {
	meta := config.GetRulesMetadata()

	doc := jsonDoc{
		Version: version,
		Status:  string(res.Status),
		Summary: jsonSummary{
			Total:    res.Summary.Total,
			Passed:   res.Summary.Passed,
			Failed:   res.Summary.Failed,
			Skipped:  res.Summary.Skipped,
			Errors:   res.Summary.Errors,
			Warnings: res.Summary.Warnings,
			Info:     res.Summary.Info,
		},
		Context: jsonContext{
			ProjectPath: res.Context.ProjectPath,
			ProjectType: res.Context.ProjectType,
		},
		Results: make([]jsonRule, 0, len(res.Results)),
	}

	for _, r := range res.Results {
		rule := jsonRule{
			RuleID:   r.RuleID,
			Category: meta.Rules[r.RuleID].Category,
			Status:   string(r.Status),
			Severity: string(r.Severity),
			Passed:   r.Passed(),
			Fixable:  meta.Rules[r.RuleID].Fixable(),
			Message:  r.Message,
			Evidence: r.Evidence,
			FixHint:  r.FixHint,
			DocURL:   r.DocURL,
		}
		doc.Results = append(doc.Results, rule)
	}
	return doc
}

func writeJSON(w io.Writer, doc jsonDoc) error {
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return fmt.Errorf("encode report: %w", err)
	}
	_, err = fmt.Fprintln(w, string(data))
	return err
}

func sortBySeverity(list []rules.RuleResult) {
	rank := map[config.Severity]int{
		config.SeverityError:   0,
		config.SeverityWarning: 1,
		config.SeverityInfo:    2,
	}
	sort.SliceStable(list, func(i, j int) bool {
		ri, rj := rank[list[i].Severity], rank[list[j].Severity]
		if ri != rj {
			return ri < rj
		}
		return list[i].RuleID < list[j].RuleID
	})
}
