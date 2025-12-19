package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/m-mdy-m/psx/internal/cmdctx"
	"github.com/m-mdy-m/psx/internal/flags"
	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/resources"
	"github.com/m-mdy-m/psx/internal/rules"
)

var FixCmd = &cobra.Command{
	Use:   "fix [path]",
	Short: "Fix structural issues",
	Long: `Automatically fix common structural issues in your project.

The fix command can:
- Create missing files (README, LICENSE, etc.)
- Create missing folders (src/, tests/, docs/)
- Generate configuration files

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
		logger.Info(resources.GetMessage("fix", "dry_run"))
	} else if f.Fix.Interactive {
		logger.Info(resources.GetMessage("fix", "interactive"))
	}
	fmt.Println()

	logger.Verbose(resources.FormatMessage("check", "start", ctx.Path.Abs))
	rulesCtx := &rules.Context{
		ProjectPath: ctx.Path.Abs,
		ProjectType: ctx.ProjectType,
		ProjectInfo: ctx.ProjectInfo,
		Config:      ctx.Config,
	}
	execResult, err := rules.Execute(ctx.Config, rulesCtx)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	failedRules := getFixableRules(execResult)

	if len(failedRules) == 0 {
		logger.Success(resources.GetMessage("fix", "success_none"))
		return nil
	}

	logger.Info(resources.FormatMessage("fix", "prompt_many", len(failedRules)))
	fmt.Println()
	if f.Fix.RuleID != "" {
		return fixSpecificRule(ctx, rulesCtx, f.Fix.RuleID)
	}
	fixCtx := &rules.FixContext{
		Context:       rulesCtx,
		Interactive:   f.Fix.Interactive && !f.Fix.All,
		DryRun:        f.Fix.DryRun,
		CreateBackups: f.Fix.CreateBackups,
	}

	results, err := rules.FixAll(ctx.Config, fixCtx, failedRules)
	if err != nil {
		return fmt.Errorf("fix failed: %w", err)
	}
	displayFixResults(results, f.Fix.DryRun)
	summary := generateSummary(results)
	displayFixSummary(summary, f.Fix.DryRun)

	if f.Fix.DryRun {
		fmt.Println()
		logger.Info("Run without --dry-run to apply changes")
		return nil
	}

	if summary.Fixed > 0 {
		fmt.Println()
		logger.Success(resources.FormatMessage("fix", "success_many", summary.Fixed))
		logger.Info("Run 'psx check' to verify")
	}

	return nil
}

func fixSpecificRule(ctx *cmdctx.ProjectContext, rulesCtx *rules.Context, ruleID string) error {
	f := flags.GetFlags()

	fixCtx := &rules.FixContext{
		Context:       rulesCtx,
		Interactive:   f.Fix.Interactive,
		DryRun:        f.Fix.DryRun,
		CreateBackups: f.Fix.CreateBackups,
	}

	result, err := rules.Fix(ctx.Config, fixCtx, ruleID)
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
			logger.Success(resources.GetMessage("fix", "success_one"))
		}
	}

	return nil
}

func getFixableRules(result *rules.ExecutionResult) []string {
	fixable := []string{}

	for _, r := range result.Results {
		if !r.Passed {
			fixable = append(fixable, r.RuleID)
		}
	}

	return fixable
}

func displayFixResults(results []*rules.FixResult, dryRun bool) {
	for _, fix := range results {
		if fix.Skipped {
			logger.Verbose(fmt.Sprintf("Skipped: %s", fix.RuleID))
			continue
		}

		if fix.Error != nil {
			logger.Error(fmt.Sprintf("%s: %v", fix.RuleID, fix.Error))
			continue
		}

		if fix.Fixed {
			logger.Verbose(resources.FormatMessage("fix", "applied", fix.RuleID))
			for _, change := range fix.Changes {
				printChange(change, dryRun)
			}
		}
	}
}

func printChange(change rules.Change, dryRun bool) {
	f := flags.GetFlags()

	if f.GlobalFlags.Quiet {
		return
	}

	prefix := "✓"
	if dryRun {
		prefix = "→"
	}

	fmt.Printf("%s %s\n", prefix, change.Description)

	if f.GlobalFlags.Verbose && change.Content != "" {
		fmt.Println(formatContent(change.Content, 5))
	}
}

type FixSummary struct {
	Total   int
	Fixed   int
	Skipped int
	Failed  int
	Changes int
}

func generateSummary(results []*rules.FixResult) FixSummary {
	summary := FixSummary{Total: len(results)}

	for _, fix := range results {
		if fix.Fixed {
			summary.Fixed++
			summary.Changes += len(fix.Changes)
		} else if fix.Skipped {
			summary.Skipped++
		} else if fix.Error != nil {
			summary.Failed++
		}
	}

	return summary
}

func displayFixSummary(summary FixSummary, dryRun bool) {
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

func formatContent(content string, maxLines int) string {
	lines := []string{}
	current := ""
	for _, char := range content {
		if char == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}

	if len(lines) <= maxLines {
		return content
	}

	displayed := ""
	for i := 0; i < maxLines; i++ {
		displayed += lines[i] + "\n"
	}
	remaining := len(lines) - maxLines
	return fmt.Sprintf("%s... (%d more lines)", displayed, remaining)
}
