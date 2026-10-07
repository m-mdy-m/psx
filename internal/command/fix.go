package command

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/m-mdy-m/psx/internal/cmdctx"
	"github.com/m-mdy-m/psx/internal/flags"
	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/rules"
	"github.com/m-mdy-m/psx/internal/ui"
)

func newFixCmd() *cobra.Command {
	opts := baseOptions()
	opts.Fix.Interactive = true

	cmd := &cobra.Command{
		Use:   "fix [path]",
		Short: "Create the files a project is missing",
		Long: "Create the files and directories that failed a check.\n\n" +
			"Nothing outside the project is touched, and an existing file with content is\n" +
			"never overwritten unless --force is given.\n\n" +
			"Examples:\n" +
			"  psx fix --dry-run          # preview without writing\n" +
			"  psx fix                    # confirm each change\n" +
			"  psx fix --yes              # apply everything unattended\n" +
			"  psx fix --rule readme      # fix a single rule\n" +
			"  psx fix --category cicd    # fix a whole category",
		Args:          cobra.MaximumNArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			applyGlobals(&opts)
			return runFixCommand(cmd, args, &opts)
		},
	}

	cobraScope(cmd, &opts)
	f := cmd.Flags()
	f.BoolVar(&opts.Fix.DryRun, "dry-run", false, "show what would be created without writing")
	f.BoolVar(&opts.Fix.Force, "force", false, "overwrite existing files that have content")
	f.StringVar(&opts.Fix.RuleID, "rule", "", "fix only this rule id")
	f.BoolVarP(&opts.Fix.Interactive, "interactive", "i", true, "confirm each change before applying")
	f.StringToStringVar(&opts.Fix.Answers, "answer", nil,
		"answer template questions, e.g. --answer ci_platform=github")

	return cmd
}

func runFixCommand(cmd *cobra.Command, args []string, opts *flags.Options) error {
	if err := opts.Validate(); err != nil {
		return commandError(exitArgs, "%v", err)
	}

	mode := cmdctx.ModeInteractive
	if opts.Fix.DryRun {
		mode = cmdctx.ModeDryRun
	}
	if opts.Global.Yes {
		opts.Fix.Interactive = false
	}

	ctx, err := cmdctx.Load(fixRoot(args), opts.Global.ConfigFile, mode)
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

	targets := fixableRules(res)
	if opts.Fix.RuleID != "" {
		targets = []string{opts.Fix.RuleID}
	}
	if len(targets) == 0 {
		logger.Success("Nothing to fix")
		return nil
	}

	if opts.Fix.DryRun {
		logger.Info("Dry run: no files will be written")
	}

	if opts.Fix.Interactive {
		targets = selectRules(targets)
		if len(targets) == 0 {
			logger.Info("No rules selected")
			return nil
		}
	}

	results := rules.FixAll(ctx.Config, &rules.FixContext{
		Context:       ctx.RuleContext(),
		Interactive:   opts.Fix.Interactive,
		DryRun:        opts.Fix.DryRun,
		CreateBackups: ctx.Config.Fix.Backup || opts.Fix.CreateBackups,
		Force:         opts.Fix.Force,
	}, targets, resourceOptions(ctx, *opts))

	reportFixResults(results, opts.Fix.DryRun, ctx.Path)
	summary := summarizeFix(results)

	if opts.Fix.DryRun {
		logger.Info("Run without --dry-run to apply")
		return nil
	}
	if summary.Failed > 0 {
		return commandError(exitCode, "%d rule(s) could not be fixed", summary.Failed)
	}
	if summary.Fixed > 0 {
		logger.Successf("Created %d file(s)", summary.Changes)
		logger.Info("Run `psx check` to verify")
	}
	return nil
}

func fixRoot(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return "."
}

// selectRules lets the user choose which failing rules to fix.
//
// The prompt is skipped entirely when there is no terminal, so an unattended run
// never blocks on stdin.
func selectRules(targets []string) []string {
	if !ui.IsInteractive() {
		return targets
	}

	present := make([]string, len(targets))
	copy(present, targets)
	sort.Strings(present)

	fmt.Fprintf(logger.Out, "Fix %d rule(s)?\n", len(present))
	choice := ui.Choose("Choose an action", []string{
		"all",
		"none",
		"select individually",
	})
	switch choice {
	case "none":
		return nil
	case "select individually":
		return selectIndividually(present)
	default:
		return present
	}
}

// selectIndividually prompts for each rule, remembering an "all" or "none" answer.
func selectIndividually(present []string) []string {
	var chosen []string
	applyToRest := false
	skipRest := false

	for _, id := range present {
		switch {
		case skipRest:
			continue
		case applyToRest:
			chosen = append(chosen, id)
			continue
		}

		switch ui.Prompt(fmt.Sprintf("Fix %s", id), []string{"yes", "no", "all", "quit"}) {
		case 0:
			chosen = append(chosen, id)
		case 2:
			chosen = append(chosen, id)
			applyToRest = true
		case 3:
			return chosen
		}
	}
	return chosen
}

// fixSummary aggregates fix outcomes.
type fixSummary struct {
	Fixed   int
	Skipped int
	Failed  int
	Changes int
}

func summarizeFix(results []*rules.FixResult) fixSummary {
	var s fixSummary
	for _, r := range results {
		switch {
		case r.Error != nil:
			s.Failed++
		case r.Fixed:
			s.Fixed++
			s.Changes += len(r.Changes)
		case r.Skipped:
			s.Skipped++
		}
	}
	return s
}

func reportFixResults(results []*rules.FixResult, dryRun bool, root string) {
	marker := "created"
	if dryRun {
		marker = "would create"
	}

	for _, r := range results {
		if r.Error != nil {
			logger.Errorf("%s: %v", r.RuleID, r.Error)
			continue
		}
		if r.Skipped {
			if r.Reason != "" {
				logger.Verbosef("skipped %s: %s", r.RuleID, r.Reason)
			}
			continue
		}
		for _, c := range r.Changes {
			logger.Step(fmt.Sprintf("%-13s %s", marker, c.Rel(root)))
			if globalOpts.verbose && c.Content != "" {
				logger.Plain(rules.Preview(c.Content, 8))
			}
		}
	}
}

// jsonUnmarshal is a small indirection so check.go need not import encoding/json.
func jsonUnmarshal(data []byte, v any) error { return json.Unmarshal(data, v) }
