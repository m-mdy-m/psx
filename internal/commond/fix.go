package commond

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/m-mdy-m/psx/internal/cmdctx"
	"github.com/m-mdy-m/psx/internal/fixer"
	"github.com/m-mdy-m/psx/internal/flags"
	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/rules"
	"github.com/m-mdy-m/psx/internal/shared"
)

var FixCmd = &cobra.Command{
	Use:   "fix [path]",
	Short: "Fix structural issues",
	Long: `Automatically fix common structural issues in your project.

The fix command can:
- Create missing files (README, LICENSE, etc.)
- Create missing folders (src/, tests/, docs/)
- Update existing files with missing content

Examples:
  psx fix                       # Interactive mode (asks before each fix)
  psx fix --dry-run             # Preview changes without applying
  psx fix --rule readme         # Fix only README
  psx fix --all                 # Fix all issues without prompts
  psx fix --create-backups      # Create backups before modifying files`,
	Args: cobra.MaximumNArgs(1),
	RunE: runFixCommand,
}

func init() {
	f := flags.GetFlags()
	df := flags.DefaultValues.Fix

	FixCmd.Flags().BoolVarP(&f.Fix.Interactive, "interactive", "i", df.Interactive,
		"ask before each fix")

	FixCmd.Flags().BoolVar(&f.Fix.DryRun, "dry-run", df.DryRun,
		"show what would be fixed without applying changes")

	FixCmd.Flags().StringVar(&f.Fix.RuleID, "rule", df.RuleID,
		"fix specific rule only")

	FixCmd.Flags().BoolVar(&f.Fix.All, "all", df.All,
		"fix all issues without prompting")

	FixCmd.Flags().BoolVar(&f.Fix.CreateBackups, "create-backups", df.CreateBackups,
		"create backup files before modifying")
}

func runFixCommand(cmd *cobra.Command, args []string) error {
	ctx, err := cmdctx.LoadProject(args)
	if err != nil {
		return err
	}

	f := flags.GetFlags()

	if f.Fix.DryRun {
		logger.Info(shared.FixDryRunHeader())
	} else if f.Fix.Interactive {
		logger.Info(shared.FixInteractiveHeader())
	}
	fmt.Println()

	logger.Verbose(shared.VerboseCheckStart(ctx.Path.Abs))

	engine := rules.NewEngine(ctx.Config, ctx.Detection)
	execResult, err := engine.Execute()
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Get fixable failed rules
	fixableRules := fixer.GetFixableFails(execResult, ctx.Config)

	if len(fixableRules) == 0 {
		logger.Success(shared.FixSuccess(0))
		return nil
	}

	logger.Info(shared.FixPrompt(len(fixableRules)))
	fmt.Println()

	// Handle specific rule fix
	if f.Fix.RuleID != "" {
		return fixSpecificRule(ctx, f.Fix.RuleID)
	}

	// Create fixer context
	fixerCtx := &fixer.FixContext{
		ProjectPath:   ctx.Path.Abs,
		ProjectType:   ctx.Detection.Type.Primary,
		Config:        ctx.Config,
		Interactive:   f.Fix.Interactive && !f.Fix.All,
		DryRun:        f.Fix.DryRun,
		CreateBackups: f.Fix.CreateBackups,
	}

	// Create fixer engine
	fixEngine := fixer.NewEngine(fixerCtx)

	// Apply fixes
	plan, err := fixEngine.FixAll(fixableRules)
	if err != nil {
		return fmt.Errorf("fix failed: %w", err)
	}

	// Display results
	displayFixResults(plan, f.Fix.DryRun)

	// Generate summary
	summary := fixer.GenerateSummary(plan)
	displayFixSummary(summary, f.Fix.DryRun)

	if f.Fix.DryRun {
		fmt.Println()
		logger.Info("Run without --dry-run to apply changes")
		return nil
	}

	if summary.Fixed > 0 {
		fmt.Println()
		logger.Success(shared.FixSuccess(summary.Fixed))
		logger.Info("Run 'psx check' to verify")
	}

	return nil
}

func fixSpecificRule(ctx *cmdctx.ProjectContext, ruleID string) error {
	f := flags.GetFlags()

	fixerCtx := &fixer.FixContext{
		ProjectPath:   ctx.Path.Abs,
		ProjectType:   ctx.Detection.Type.Primary,
		Config:        ctx.Config,
		Interactive:   f.Fix.Interactive,
		DryRun:        f.Fix.DryRun,
		CreateBackups: f.Fix.CreateBackups,
	}

	fixEngine := fixer.NewEngine(fixerCtx)

	if !fixEngine.CanFix(ruleID) {
		return logger.Errorf("no fix available for rule: %s", ruleID)
	}

	result, err := fixEngine.Fix(ruleID)
	if err != nil {
		return fmt.Errorf("fix failed: %w", err)
	}

	if result.Skipped {
		logger.Info("Fix skipped")
		return nil
	}

	if result.Fixed {
		for _, change := range result.Changes {
			printChange(change, f.Fix.DryRun)
		}

		if f.Fix.DryRun {
			fmt.Println()
			logger.Info("Run without --dry-run to apply")
		} else {
			fmt.Println()
			logger.Success(shared.FixSuccess(1))
		}
	}

	return nil
}

func displayFixResults(plan *fixer.FixPlan, dryRun bool) {
	for _, fix := range plan.Fixes {
		if fix.Skipped {
			logger.Verbose(fmt.Sprintf("Skipped: %s", fix.RuleID))
			continue
		}

		if fix.Error != nil {
			logger.Error(fmt.Sprintf("%s: %v", fix.RuleID, fix.Error))
			continue
		}

		if fix.Fixed {
			logger.Verbose(shared.VerboseFixApplied(fix.RuleID))
			for _, change := range fix.Changes {
				printChange(change, dryRun)
			}
		}
	}
}

func printChange(change fixer.Change, dryRun bool) {
	f := flags.GetFlags()

	if f.GlobalFlags.Quiet {
		return
	}

	prefix := "✓"
	if dryRun {
		prefix = "→"
	}

	switch change.Type {
	case fixer.ChangeCreateFile:
		fmt.Printf("%s %s\n", prefix, change.Description)
		if f.GlobalFlags.Verbose && change.Content != "" {
			fmt.Println(fixer.FormatContent(change.Content, 5))
		}

	case fixer.ChangeCreateFolder:
		fmt.Printf("%s %s\n", prefix, change.Description)

	case fixer.ChangeModifyFile:
		fmt.Printf("%s %s\n", prefix, change.Description)
	}
}

func displayFixSummary(summary fixer.FixSummary, dryRun bool) {
	f := flags.GetFlags()

	if f.GlobalFlags.Quiet {
		return
	}

	fmt.Println()
	fmt.Println("Summary:")
	fmt.Printf("  Total:   %d\n", summary.Total)

	if dryRun {
		fmt.Printf("  Would fix: %d\n", summary.Fixed)
	} else {
		fmt.Printf("  Fixed:   %d\n", summary.Fixed)
	}

	if summary.Skipped > 0 {
		fmt.Printf("  Skipped: %d\n", summary.Skipped)
	}

	if summary.Failed > 0 {
		fmt.Printf("  Failed:  %d\n", summary.Failed)
	}

	if summary.Changes > 0 {
		fmt.Printf("  Changes: %d\n", summary.Changes)
	}
}

