package report

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/m-mdy-m/psx/internal/rules"
)

type jsonRenderer struct{ opts Options }

func (r *jsonRenderer) Render(w io.Writer, res *rules.ExecutionResult) error {
	doc := toDoc(res, r.opts.Version)
	if !r.opts.ShowPassed && !r.opts.Verbose {
		doc.Results = filterOutPassed(doc.Results)
	}
	return writeJSON(w, doc)
}

func filterOutPassed(list []jsonRule) []jsonRule {
	out := make([]jsonRule, 0, len(list))
	for _, r := range list {
		if r.Passed {
			continue
		}
		out = append(out, r)
	}
	return out
}

type ndjsonRenderer struct{ opts Options }

func (r *ndjsonRenderer) Render(w io.Writer, res *rules.ExecutionResult) error {
	doc := toDoc(res, r.opts.Version)

	summary, err := json.Marshal(map[string]any{
		"type":    "summary",
		"status":  doc.Status,
		"summary": doc.Summary,
		"context": doc.Context,
	})
	if err != nil {
		return fmt.Errorf("encode summary: %w", err)
	}
	if _, err := fmt.Fprintf(w, "%s\n", summary); err != nil {
		return err
	}

	for _, rule := range doc.Results {
		if rule.Passed && !r.opts.ShowPassed && !r.opts.Verbose {
			continue
		}
		line, err := json.Marshal(map[string]any{
			"type":     "rule",
			"rule_id":  rule.RuleID,
			"status":   rule.Status,
			"severity": rule.Severity,
			"message":  rule.Message,
		})
		if err != nil {
			return fmt.Errorf("encode rule: %w", err)
		}
		if _, err := fmt.Fprintf(w, "%s\n", line); err != nil {
			return err
		}
	}
	return nil
}

type sarifRenderer struct{ opts Options }

func (r *sarifRenderer) Render(w io.Writer, res *rules.ExecutionResult) error {
	var results []sarifResult

	for _, x := range res.Results {
		if x.Status != rules.StatusFailed {
			continue
		}
		level := "note"
		switch x.Severity {
		case "error":
			level = "error"
		case "warning":
			level = "warning"
		}

		results = append(results, sarifResult{
			RuleID:              x.RuleID,
			Level:               level,
			Message:             sarifMessage{Text: x.Message},
			Help:                sarifMessage{Text: x.FixHint},
			PartialFingerPrints: map[string]string{"psxRuleId/v1": x.RuleID},
		})
	}

	doc := sarifDoc{
		Schema:  "https://json.schemastore.org/sarif-2.1.0.json",
		Version: "2.1.0",
		Runs: []sarifRun{{
			Tool: sarifTool{Driver: sarifDriver{
				Name:           "psx",
				InformationURI: "https://github.com/m-mdy-m/psx",
				Version:        orUnknown(r.opts.Version),
				Rules:          sarifRules(res),
			}},
			Results: results,
		}},
	}

	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return fmt.Errorf("encode sarif: %w", err)
	}
	_, err = fmt.Fprintln(w, string(data))
	return err
}

type sarifDoc struct {
	Schema  string     `json:"$schema"`
	Version string     `json:"version"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool    sarifTool     `json:"tool"`
	Results []sarifResult `json:"results"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name           string      `json:"name"`
	InformationURI string      `json:"informationUri"`
	Version        string      `json:"version,omitempty"`
	Rules          []sarifRule `json:"rules"`
}

type sarifRule struct {
	ID               string            `json:"id"`
	Name             string            `json:"name,omitempty"`
	ShortDescription sarifMessage      `json:"shortDescription"`
	HelpURI          string            `json:"helpUri,omitempty"`
	Properties       map[string]string `json:"properties,omitempty"`
}

type sarifResult struct {
	RuleID              string            `json:"ruleId"`
	Level               string            `json:"level"`
	Message             sarifMessage      `json:"message"`
	Help                sarifMessage      `json:"help,omitempty"`
	PartialFingerPrints map[string]string `json:"partialFingerprints,omitempty"`
	Properties          []string          `json:"-"`
}

func sarifRules(res *rules.ExecutionResult) []sarifRule {
	seen := map[string]bool{}
	out := make([]sarifRule, 0, len(res.Results))

	for _, x := range res.Results {
		if x.Status != rules.StatusFailed || seen[x.RuleID] {
			continue
		}
		seen[x.RuleID] = true
		out = append(out, sarifRule{
			ID:               x.RuleID,
			ShortDescription: sarifMessage{Text: x.Message},
			HelpURI:          x.DocURL,
		})
	}
	return out
}

type sarifMessage struct {
	Text string `json:"text"`
}

type githubRenderer struct{ opts Options }

func (r *githubRenderer) Render(w io.Writer, res *rules.ExecutionResult) error {
	for _, x := range res.Results {
		if x.Status != rules.StatusFailed {
			continue
		}
		level := "notice"
		switch x.Severity {
		case "error":
			level = "error"
		case "warning":
			level = "warning"
		}

		props := []string{fmt.Sprintf("title=%s", x.RuleID)}
		if x.FixHint != "" {
			props = append(props, fmt.Sprintf("hint=%s", x.FixHint))
		}
		fmt.Fprintf(w, "::%s file=%s::%s\n",
			level, escapeProp(r.opts.ProjectPath), escapeData(x.Message))
		_ = props
	}

	if !r.opts.Quiet {
		label, _ := statusStyle(res.Status)
		fmt.Fprintf(w, "%s: %d errors, %d warnings\n",
			label, res.Summary.Errors, res.Summary.Warnings)
	}
	return nil
}

func escapeProp(s string) string {
	s = strings.ReplaceAll(s, "%", "%25")
	s = strings.ReplaceAll(s, "\r", "%0D")
	s = strings.ReplaceAll(s, "\n", "%0A")
	return s
}

func escapeData(s string) string {
	s = strings.ReplaceAll(s, "%", "%25")
	s = strings.ReplaceAll(s, "\r", "%0D")
	s = strings.ReplaceAll(s, "\n", "%0A")
	return s
}

type junitRenderer struct{ opts Options }

func (r *junitRenderer) Render(w io.Writer, res *rules.ExecutionResult) error {
	var failures []string
	for _, x := range res.Results {
		if x.Status != rules.StatusFailed {
			continue
		}
		failures = append(failures, fmt.Sprintf(
			"    <testcase classname=\"psx.%s\" name=\"%s\">\n"+
				"      <failure message=\"%s\">%s</failure>\n"+
				"    </testcase>",
			xmlEscape(string(x.Severity)), xmlEscape(x.RuleID),
			xmlEscape(x.Message), xmlEscape(x.FixHint)))
	}

	total := res.Summary.Passed + res.Summary.Failed
	fmt.Fprintln(w, `<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Fprintf(w, "<testsuites tests=\"%d\" failures=\"%d\">\n", total, len(failures))
	fmt.Fprintln(w, `  <testsuite name="psx" tests="`+
		fmt.Sprint(total)+`" failures="`+fmt.Sprint(len(failures))+`">`)
	for _, f := range failures {
		fmt.Fprintln(w, f)
	}
	fmt.Fprintln(w, "  </testsuite>")
	fmt.Fprintln(w, "</testsuites>")
	return nil
}

func xmlEscape(s string) string {
	r := strings.NewReplacer(
		"&", "&amp;", "<", "&lt;", ">", "&gt;",
		`"`, "&quot;", "'", "&apos;",
	)
	return r.Replace(s)
}

func orUnknown(v string) string {
	if strings.TrimSpace(v) == "" {
		return "dev"
	}
	return v
}
