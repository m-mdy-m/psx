package command

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/m-mdy-m/psx/internal/flags"
	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/rules"
)

func newCheckCmd() *cobra.Command {
	opts := baseOptions()

	cmd := &cobra.Command{
		Use:   "check [path]",
		Short: "Validate project structure",
		Long: "Check a project against the active rules.\n\n" +
			"This command only reads: it never prompts and never modifies files.\n\n" +
			"Examples:\n" +
			"  psx check\n" +
			"  psx check ./sub-project\n" +
			"  psx check --only readme,license\n" +
			"  psx check --category documentation\n" +
			"  psx check -o json | jq .summary",
		Args:          cobra.MaximumNArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			applyGlobals(&opts)
			return runCheckCommand(cmd, args, &opts)
		},
	}

	cobraCheck(cmd, &opts)
	cobraBaseline(cmd, &opts)
	return cmd
}

func runCheckCommand(cmd *cobra.Command, args []string, opts *flags.Options) error {
	if err := opts.Validate(); err != nil {
		return commandError(exitArgs, "%v", err)
	}

	ctx, err := loadReadOnly(args, *opts)
	if err != nil {
		return err
	}
	if err := filterRules(ctx, *opts); err != nil {
		return err
	}

	res, err := runCheck(ctx)
	if err != nil {
		return err
	}

	if base := opts.Check.BaselineFile; base != "" {
		if err := applyBaseline(res, base); err != nil {
			return err
		}
	}

	if err := emit(res, *opts); err != nil {
		return err
	}

	if shouldFail(res, opts.Check.FailOn) {
		os.Exit(exitCode)
	}
	return nil
}

func shouldFail(res *rules.ExecutionResult, failOn string) bool {
	switch failOn {
	case "none":
		return false
	case "warning":
		return res.Summary.Errors > 0 || res.Summary.Warnings > 0
	default:
		return res.Summary.Errors > 0
	}
}

type baselineRule struct {
	RuleID   string `json:"rule_id"`
	Severity string `json:"severity"`
}

// applyBaseline bootstraps a missing baseline from the current run, then filters.
func applyBaseline(res *rules.ExecutionResult, path string) error {
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		if err := recordBaseline(res, path); err != nil {
			return commandError(exitConfig, "baseline: %v", err)
		}
	}
	return applyBaselineTo(res, path)
}

// applyBaselineTo marks recorded rules as passed. A missing file forgives nothing,
// so watch never invents a baseline as a side effect of running.
func applyBaselineTo(res *rules.ExecutionResult, path string) error {
	known, err := loadBaseline(path)
	if err != nil {
		return commandError(exitConfig, "baseline: %v", err)
	}

	for i := range res.Results {
		r := &res.Results[i]
		if r.Status != rules.StatusFailed {
			continue
		}
		if !known[r.RuleID] {
			continue
		}
		r.Status = rules.StatusPassed
		r.Message = "known issue, recorded in " + path
		r.FixHint = ""
		r.Severity = ""
	}

	rules.Recount(res)
	return nil
}

// recordBaseline writes the current failures to path, so pointing --baseline at a
// file that does not exist yet adopts today's state instead of erroring.
func recordBaseline(res *rules.ExecutionResult, path string) error {
	known := map[string]bool{}
	var lines []string

	for _, r := range res.Results {
		if r.Status == rules.StatusFailed {
			known[r.RuleID] = true
			lines = append(lines, r.RuleID)
		}
	}

	var buf strings.Builder
	buf.WriteString("# psx baseline. Rules listed here are reported as known issues and\n")
	buf.WriteString("# do not fail the build. Delete a line once you have fixed it.\n")
	for _, id := range lines {
		fmt.Fprintf(&buf, "%s\n", id)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(buf.String()), 0o644); err != nil {
		return err
	}
	logger.Verbosef("baseline recorded %d known issues in %s", len(known), path)
	return nil
}

func loadBaseline(path string) (map[string]bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	trimmed := strings.TrimSpace(string(data))
	out := map[string]bool{}

	if strings.HasPrefix(trimmed, "[") {
		var ids []string
		if err := jsonUnmarshal([]byte(trimmed), &ids); err == nil {
			for _, id := range ids {
				out[id] = true
			}
			return out, nil
		}
		var objs []baselineRule
		if err := jsonUnmarshal([]byte(trimmed), &objs); err != nil {
			return nil, err
		}
		for _, o := range objs {
			out[o.RuleID] = true
		}
		return out, nil
	}

	for _, line := range strings.Split(trimmed, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		out[line] = true
	}
	return out, nil
}
