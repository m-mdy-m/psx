package command

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/m-mdy-m/psx/internal/cmdctx"
	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/flags"
	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/rules"
	"github.com/m-mdy-m/psx/internal/watch"
)

func newWatchCmd() *cobra.Command {
	opts := baseOptions()
	wOpts := watch.Defaults()
	opts.Watch.Interval = wOpts.Interval
	opts.Watch.Debounce = wOpts.Debounce

	cmd := &cobra.Command{
		Use:   "watch [path]",
		Short: "Re-check the project as files change",
		Long: "Watch the project and report each change in rule status.\n\n" +
			"Only the differences are printed, so a long session stays readable.\n\n" +
			"Examples:\n" +
			"  psx watch\n" +
			"  psx watch --fix           # also create missing files\n" +
			"  psx watch --once-clean    # exit once the project is clean",
		Args:          cobra.MaximumNArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			applyGlobals(&opts)
			return runWatch(args, &opts)
		},
	}

	cobraCheck(cmd, &opts)
	cobraBaseline(cmd, &opts)
	f := cmd.Flags()
	f.BoolVar(&opts.Watch.Fix, "fix", false, "apply safe fixes as problems appear")
	f.BoolVar(&opts.Watch.Once, "once-clean", false, "exit as soon as no problems remain")
	f.DurationVar(&opts.Watch.Interval, "interval", opts.Watch.Interval, "polling interval")
	f.DurationVar(&opts.Watch.Debounce, "debounce", opts.Watch.Debounce, "settle time after a change")
	return cmd
}

func runWatch(args []string, opts *flags.Options) error {
	if opts.Watch.Interval < 100*time.Millisecond {
		return commandError(exitArgs, "--interval must be at least 100ms")
	}

	ctx, err := cmdctx.Load(watchRoot(args), opts.Global.ConfigFile, cmdctx.ModeReadOnly)
	if err != nil {
		return err
	}

	runCheck := func() (*rules.ExecutionResult, error) {
		// Rescan on every evaluation so fixes applied by a previous iteration
		// are visible to the next one.
		snap, err := watch.Snapshot(ctx.Path, ctx.Config.Ignore)
		if err != nil {
			return nil, err
		}
		res, err := rules.ExecuteSnapshot(ctx.Config, ctx.RuleContext(), snap)
		if err != nil {
			return nil, err
		}
		if base := opts.Check.BaselineFile; base != "" {
			// Missing is fine here: watch should not create a baseline as a
			// side effect of running, so only a readable one filters.
			if _, statErr := os.Stat(base); statErr == nil {
				if err := applyBaselineTo(res, base); err != nil {
					return nil, err
				}
			}
		}
		return res, nil
	}

	signals, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	wOpts := watch.Options{
		Root:       ctx.Path,
		ConfigFile: ctx.Config.ConfigFile,
		Ignore:     ctx.Config.Ignore,
		Interval:   opts.Watch.Interval,
		Debounce:   opts.Watch.Debounce,
		OnceClean:  opts.Watch.Once,
	}

	logger.Info("Watching for changes. Press Ctrl+C to stop.")
	err = watch.Run(signals, wOpts, runCheck, func(prev, cur *rules.ExecutionResult) {
		printWatchChange(prev, cur)
		if opts.Watch.Fix {
			applyWatchFixes(ctx, cur, opts)
		}
	})

	if err != nil && signals.Err() == nil {
		return commandError(exitCode, "watch: %v", err)
	}
	logger.Info("Stopped watching")
	return nil
}

func watchRoot(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return "."
}

func printWatchChange(prev, cur *rules.ExecutionResult) {
	if prev == nil {
		printSummaryLine(cur)
		return
	}

	diffs := watch.Compare(prev, cur)
	if len(diffs) == 0 {
		logger.Verbose("no rule status changes")
		return
	}

	for _, d := range diffs {
		switch d.To {
		case rules.StatusPassed:
			logger.Successf("%s: fixed", d.RuleID)
		case rules.StatusFailed:
			logger.Warningf("%s: %s", d.RuleID, d.Message)
		default:
			logger.Verbosef("%s: no longer applicable", d.RuleID)
		}
	}
	printSummaryLine(cur)
}

func printSummaryLine(res *rules.ExecutionResult) {
	s := res.Summary
	line := fmt.Sprintf("%d errors · %d warnings · %d skipped · watching",
		s.Errors, s.Warnings, s.Skipped)
	if s.Failed == 0 {
		logger.Success(line)
		return
	}
	logger.Warning(line)
}

func applyWatchFixes(ctx *cmdctx.ProjectContext, res *rules.ExecutionResult, opts *flags.Options) {
	var safe []string
	meta := config.GetRulesMetadata().Rules

	for _, r := range res.Results {
		if r.Status != rules.StatusFailed {
			continue
		}
		m, ok := meta[r.RuleID]
		if !ok || !m.Fixable() || (m.Fix != nil && !m.Fix.Safe) {
			continue
		}
		safe = append(safe, r.RuleID)
	}
	if len(safe) == 0 {
		return
	}

	results := rules.FixAll(ctx.Config, &rules.FixContext{
		Context: ctx.RuleContext(),
	}, safe, resourceOptions(ctx, *opts))

	for _, fr := range results {
		for _, c := range fr.Changes {
			logger.Step("created " + c.Description)
		}
	}
}
