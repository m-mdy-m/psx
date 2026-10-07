package command

import (
	"github.com/spf13/cobra"

	"github.com/m-mdy-m/psx/internal/cmdctx"
	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/flags"
	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/report"
	"github.com/m-mdy-m/psx/internal/resources"
	"github.com/m-mdy-m/psx/internal/rules"
)

// baseOptions returns a fresh option set with defaults applied.
//
// It deliberately does not read the global flags: it runs while the command tree is
// being built, before cobra has parsed anything, so copying the flag values here
// would capture the defaults and silently ignore whatever the user passed.
func baseOptions() flags.Options {
	opts := flags.Defaults()
	opts.Version = Version
	return opts
}

// applyGlobals copies the parsed persistent flags onto a command's options.
//
// This must run inside RunE, where the flag values are final.
func applyGlobals(opts *flags.Options) {
	opts.Global = flags.GlobalFlags{
		ConfigFile: globalOpts.configFile,
		Verbose:    globalOpts.verbose,
		Quiet:      globalOpts.quiet,
		NoColor:    globalOpts.noColor,
		Yes:        globalOpts.yes,
	}
}

func loadReadOnly(args []string, opts flags.Options) (*cmdctx.ProjectContext, error) {
	root := "."
	if len(args) > 0 {
		root = args[0]
	}
	return cmdctx.Load(root, opts.Global.ConfigFile, cmdctx.ModeReadOnly)
}

func runCheck(ctx *cmdctx.ProjectContext) (*rules.ExecutionResult, error) {
	logger.Verbosef("checking %s (%s, %s)", ctx.Path, ctx.ProjectType, ctx.Kind)

	res, err := rules.ExecuteSnapshot(ctx.Config, ctx.RuleContext(), ctx.Snapshot)
	if err != nil {
		return nil, commandError(exitConfig, "check failed: %v", err)
	}

	logger.Verbosef("%d passed, %d failed, %d skipped",
		res.Summary.Passed, res.Summary.Failed, res.Summary.Skipped)
	return res, nil
}

func emit(res *rules.ExecutionResult, opts flags.Options) error {
	renderer, err := report.New(report.Options{
		Format:      opts.Check.OutputFormat,
		Level:       opts.Check.Level,
		ShowPassed:  opts.Check.ShowPassed,
		ShowSkipped: opts.Check.ShowSkipped,
		Verbose:     opts.Verbose(),
		Quiet:       opts.Global.Quiet,
		NoColor:     opts.Global.NoColor,
		ProjectPath: res.Context.ProjectPath,
		ProjectType: res.Context.ProjectType,
		Version:     opts.Version,
	})
	if err != nil {
		return commandError(exitArgs, "%v", err)
	}
	return renderer.Render(logger.Out, res)
}

func filterRules(ctx *cmdctx.ProjectContext, opts flags.Options) error {
	if len(opts.Check.Only) == 0 && len(opts.Check.Category) == 0 {
		return nil
	}

	meta := config.GetRulesMetadata()
	keep := make(map[string]bool, len(opts.Check.Only))
	for _, id := range opts.Check.Only {
		if _, ok := ctx.Config.ActiveRules[id]; !ok {
			return commandError(exitArgs, "unknown rule %q (run `psx rules list`)", id)
		}
		keep[id] = true
	}
	for _, cat := range opts.Check.Category {
		for id, m := range meta.Rules {
			if m.Category == cat {
				keep[id] = true
			}
		}
	}

	filtered := make(map[string]*config.ActiveRule, len(keep))
	for id := range keep {
		if rule, ok := ctx.Config.ActiveRules[id]; ok {
			filtered[id] = rule
		}
	}
	if len(filtered) == 0 {
		return commandError(exitArgs, "no active rules match the requested filters")
	}
	ctx.Config.ActiveRules = filtered
	return nil
}

func fixableRules(res *rules.ExecutionResult) []string {
	var out []string
	for _, r := range res.Results {
		if r.Status != rules.StatusFailed {
			continue
		}
		out = append(out, r.RuleID)
	}
	return out
}

func resourceOptions(ctx *cmdctx.ProjectContext, opts flags.Options) resources.Options {
	return resources.Options{
		ProjectType: ctx.ProjectType,
		Answers:     opts.Fix.Answers,
		Force:       opts.Fix.Force,
	}
}

func cobraCheck(cmd *cobra.Command, opts *flags.Options) {
	f := cmd.Flags()
	f.StringVarP(&opts.Check.OutputFormat, "output", "o", opts.Check.OutputFormat,
		"output format: table, compact, json, ndjson, sarif, github, junit, markdown")
	f.StringVar(&opts.Check.Level, "level", opts.Check.Level,
		"minimum severity to report: error, warning, info, all")
	f.StringVar(&opts.Check.FailOn, "fail-on", opts.Check.FailOn,
		"exit non-zero on: error, warning, none")
	f.StringSliceVar(&opts.Check.Only, "only", nil, "check only these rule ids")
	f.StringSliceVar(&opts.Check.Category, "category", nil, "check only these categories")
	f.BoolVar(&opts.Check.ShowPassed, "all", false, "include passing rules in the output")
	f.BoolVar(&opts.Check.ShowSkipped, "show-skipped", false, "include rules that did not apply")
}

func cobraBaseline(cmd *cobra.Command, opts *flags.Options) {
	cmd.Flags().StringVar(&opts.Check.BaselineFile, "baseline", "",
		"only fail on rules not already recorded in this file")
}

func cobraScope(cmd *cobra.Command, opts *flags.Options) {
	f := cmd.Flags()
	f.StringSliceVar(&opts.Check.Only, "only", nil, "act only on these rule ids")
	f.StringSliceVar(&opts.Check.Category, "category", nil, "act only on these categories")
}
