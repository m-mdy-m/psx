package command

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/m-mdy-m/psx/internal/cmdctx"
	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/detect"
	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/resources"
)

// newRulesCmd builds the rules command, which documents what psx checks.
func newRulesCmd() *cobra.Command {
	var (
		asJSON   bool
		category string
		fixable  bool
	)

	cmd := &cobra.Command{
		Use:   "rules",
		Short: "List and inspect the available rules",
		Long:  "Show every rule psx knows about, including its severity and whether it can be fixed automatically.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRules(asJSON, category, fixable)
		},
	}

	cmd.Flags().BoolVar(&asJSON, "json", false, "emit JSON")
	cmd.Flags().StringVar(&category, "category", "", "filter by category")
	cmd.Flags().BoolVar(&fixable, "fixable", false, "show only rules with an automatic fix")
	return cmd
}

func runRules(asJSON bool, category string, fixableOnly bool) error {
	all := config.GetRulesMetadata().Rules

	ids := make([]string, 0, len(all))
	for id, m := range all {
		if category != "" && m.Category != category {
			continue
		}
		if fixableOnly && !m.Fixable() {
			continue
		}
		ids = append(ids, id)
	}
	sort.Strings(ids)

	if asJSON {
		return printRuleJSON(ids, all)
	}

	if len(ids) == 0 {
		logger.Info("No rules match the filter")
		return nil
	}

	w := tabwriter.NewWriter(logger.Out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "RULE\tCATEGORY\tSEVERITY\tFIX\tDESCRIPTION")
	for _, id := range ids {
		m := all[id]
		fix := "no"
		if m.Fixable() {
			fix = "yes"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			id, m.Category, m.DefaultSeverity, fix, m.Description)
	}
	return w.Flush()
}

func printRuleJSON(ids []string, all map[string]config.RuleMetadata) error {
	out := make([]map[string]any, 0, len(ids))
	for _, id := range ids {
		m := all[id]
		out = append(out, map[string]any{
			"id":          id,
			"category":    m.Category,
			"severity":    m.DefaultSeverity,
			"description": m.Description,
			"fixable":     m.Fixable(),
			"fix":         fixSpecJSON(m),
			"patterns":    m.Patterns,
			"message":     m.Message,
			"doc_url":     m.DocURL,
		})
	}

	data, err := json.MarshalIndent(map[string]any{"rules": out}, "", "  ")
	if err != nil {
		return err
	}
	logger.Plain(string(data))
	return nil
}

// fixSpecJSON describes what a rule's fix would create.
func fixSpecJSON(m config.RuleMetadata) any {
	if !m.Fixable() {
		return nil
	}
	files := make([]string, 0, len(m.Fix.Files))
	for path := range m.Fix.Files {
		files = append(files, path)
	}
	sort.Strings(files)

	return map[string]any{
		"path":  m.Fix.Path,
		"files": files,
		"safe":  m.Fix.Safe,
	}
}

// newWorkflowsCmd builds the workflows command for discovering CI templates.
func newWorkflowsCmd() *cobra.Command {
	var (
		asJSON bool
		group  string
		show   bool
	)

	cmd := &cobra.Command{
		Use:   "workflows",
		Short: "List the GitHub Actions templates",
		Long: "List the GitHub Actions workflow templates bundled with psx.\n\n" +
			"Each template is a standalone workflow file. Use --show to print one,\n" +
			"or add a workflow with `psx fix --rule workflow_<name>`.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWorkflows(group, asJSON, show)
		},
	}

	cmd.Flags().BoolVar(&asJSON, "json", false, "emit JSON")
	cmd.Flags().StringVar(&group, "group", "", "filter by group: ci, release, docker, deploy, security")
	cmd.Flags().BoolVar(&show, "show", false, "print the full body of each listed template")
	return cmd
}

func runWorkflows(group string, asJSON, show bool) error {
	groups := resources.WorkflowGroups()
	names := resources.WorkflowNames()

	if group != "" {
		groups = filterWorkflowGroups(groups, group)
		if len(groups) == 0 {
			return commandError(exitArgs, "unknown workflow group %q", group)
		}
		names = nil
		for _, g := range groups {
			names = append(names, g.Templates...)
		}
		sort.Strings(names)
	}

	if asJSON {
		return printWorkflowsJSON(names, show)
	}

	w := tabwriter.NewWriter(logger.Out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TEMPLATE\tGROUP\tDESCRIPTION")
	for _, g := range groups {
		for _, name := range g.Templates {
			if !containsName(names, name) {
				continue
			}
			fmt.Fprintf(w, "%s\t%s\t%s\n", name, g.Group, g.Purpose)
		}
	}
	if err := w.Flush(); err != nil {
		return err
	}

	if !show {
		return nil
	}

	// Rendering needs project variables; use placeholders for a neutral preview.
	vars := resources.PreviewVars()
	for _, name := range names {
		body, err := resources.Workflow(name, "generic", vars)
		if err != nil {
			continue
		}
		fmt.Fprintf(logger.Out, "\n───── %s ─────\n%s", name, body)
	}
	return nil
}

// filterWorkflowGroups keeps the groups matching a name.
func filterWorkflowGroups(groups []resources.WorkflowGroup, name string) []resources.WorkflowGroup {
	var out []resources.WorkflowGroup
	for _, g := range groups {
		if g.Group == name {
			out = append(out, g)
		}
	}
	return out
}

func containsName(list []string, name string) bool {
	for _, v := range list {
		if v == name {
			return true
		}
	}
	return false
}

func printWorkflowsJSON(names []string, show bool) error {
	groupOf := map[string]string{}
	for _, g := range resources.WorkflowGroups() {
		for _, n := range g.Templates {
			groupOf[n] = g.Group
		}
	}

	out := make([]map[string]any, 0, len(names))
	for _, name := range names {
		entry := map[string]any{
			"name":  name,
			"group": groupOf[name],
			"path":  resources.WorkflowPath(name),
		}
		if show {
			if body, err := resources.Workflow(name, "generic", resources.PreviewVars()); err == nil {
				entry["content"] = body
			}
		}
		out = append(out, entry)
	}

	data, err := json.MarshalIndent(map[string]any{"workflows": out}, "", "  ")
	if err != nil {
		return err
	}
	logger.Plain(string(data))
	return nil
}

// newExplainCmd builds the explain command for a single rule.
func newExplainCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "explain <rule>",
		Short: "Explain a rule in detail",
		Long:  "Show what a rule checks, why it matters, and what `psx fix` would create.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExplain(args[0])
		},
	}
}

func runExplain(id string) error {
	meta, ok := config.GetRulesMetadata().Rules[id]
	if !ok {
		return commandError(exitArgs, "unknown rule %q (run `psx rules`)", id)
	}

	out := logger.Out
	fmt.Fprintf(out, "%s\n", id)
	fmt.Fprintf(out, "  category    %s\n", meta.Category)
	fmt.Fprintf(out, "  severity    %s\n", meta.DefaultSeverity)
	fmt.Fprintf(out, "  description %s\n\n", meta.Description)

	fmt.Fprintf(out, "  message: %s\n\n", meta.Message)

	if patterns := describePatterns(meta.Patterns); patterns != "" {
		fmt.Fprintf(out, "  checks for:\n%s\n", patterns)
	}
	if meta.Fixable() {
		fmt.Fprintf(out, "  fix creates:\n")
		if meta.Fix.Path != "" {
			fmt.Fprintf(out, "    %s\n", meta.Fix.Path)
		}
		paths := make([]string, 0, len(meta.Fix.Files))
		for p := range meta.Fix.Files {
			paths = append(paths, p)
		}
		sort.Strings(paths)
		for _, p := range paths {
			fmt.Fprintf(out, "    %s\n", p)
		}
		fmt.Fprintln(out)
	} else {
		fmt.Fprintf(out, "  fix: none (manual)\n\n")
	}

	if meta.FixHint != "" {
		fmt.Fprintf(out, "  hint: %s\n", meta.FixHint)
	}
	if meta.DocURL != "" {
		fmt.Fprintf(out, "  docs: %s\n", meta.DocURL)
	}
	return nil
}

func describePatterns(patterns any) string {
	var lines []string

	switch p := patterns.(type) {
	case []any:
		for _, item := range p {
			if s, ok := item.(string); ok {
				lines = append(lines, "    "+s)
			}
		}
	case map[string]any:
		keys := make([]string, 0, len(p))
		for k := range p {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			lines = append(lines, fmt.Sprintf("    [%s]", k))
			list, _ := p[k].([]any)
			for _, item := range list {
				if s, ok := item.(string); ok {
					lines = append(lines, "      "+s)
				}
			}
		}
	}
	return strings.Join(lines, "\n")
}

// newDetectCmd builds the detect command.
func newDetectCmd() *cobra.Command {
	var asJSON bool

	cmd := &cobra.Command{
		Use:   "detect [path]",
		Short: "Infer the project type from its files",
		Long: "Report the language and project archetype psx infers, along with the signals used.\n\n" +
			"A declared `project.type` in the configuration always wins.",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := "."
			if len(args) > 0 {
				root = args[0]
			}
			return runDetect(root, asJSON)
		},
	}

	cmd.Flags().BoolVar(&asJSON, "json", false, "emit JSON")
	return cmd
}

func runDetect(root string, asJSON bool) error {
	abs, err := cmdctx.ResolvePath(root)
	if err != nil {
		return commandError(exitArgs, "%v", err)
	}

	profile, err := detect.Detect(abs, nil)
	if err != nil {
		return commandError(exitConfig, "detect: %v", err)
	}

	if asJSON {
		data, err := json.MarshalIndent(profile, "", "  ")
		if err != nil {
			return err
		}
		logger.Plain(string(data))
		return nil
	}

	out := logger.Out
	fmt.Fprintf(out, "type       %s\n", profile.ProjectType)
	fmt.Fprintf(out, "kind       %s\n", profile.Kind)
	fmt.Fprintf(out, "confidence %.0f%%\n", profile.Confidence*100)
	if len(profile.Workspace) > 0 {
		fmt.Fprintf(out, "workspace  %s\n", strings.Join(profile.Workspace, ", "))
	}
	fmt.Fprintf(out, "signals    %s\n", strings.Join(profile.Signals, ", "))

	if len(profile.Managers) > 0 {
		fmt.Fprintln(out, "\npackage managers:")
		for _, m := range profile.Managers {
			fmt.Fprintf(out, "  %-8s %s (lock: %s)\n", m.Name, m.File, m.Lock)
		}
	}
	return nil
}

// writeFileIfAbsent is a small helper for the init command.
func writeFileIfAbsent(path, content string, force bool) error {
	if !force {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("%s already exists", path)
		}
	}
	return os.WriteFile(path, []byte(content), 0o644)
}
