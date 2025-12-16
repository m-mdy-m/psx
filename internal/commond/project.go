package commond

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/m-mdy-m/psx/internal/cmdctx"
	"github.com/m-mdy-m/psx/internal/flags"
	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/resources"
)

var ProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage project information",
	Long: `Manage project information cache.

The project information (name, author, email, etc.) is cached in .psx-project.yml
to avoid asking you the same questions every time.

Commands:
  psx project show     # Show current project info
  psx project reset    # Reset project info (will ask again next time)
  psx project edit     # Edit project info interactively`,
}

var projectShowCmd = &cobra.Command{
	Use:   "show [path]",
	Short: "Show project information",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runProjectShow,
}

var projectResetCmd = &cobra.Command{
	Use:   "reset [path]",
	Short: "Reset project information cache",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runProjectReset,
}

var projectEditCmd = &cobra.Command{
	Use:   "edit [path]",
	Short: "Edit project information",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runProjectEdit,
}

func init() {
	ProjectCmd.AddCommand(projectShowCmd)
	ProjectCmd.AddCommand(projectResetCmd)
	ProjectCmd.AddCommand(projectEditCmd)
}

func runProjectShow(cmd *cobra.Command, args []string) error {
	pathCtx, err := cmdctx.ResolvePath(args)
	if err != nil {
		return err
	}

	info, err := resources.LoadProjectInfo(pathCtx.Abs)
	if err != nil {
		return logger.Errorf("Failed to load project info: %w", err)
	}

	if info == nil {
		logger.Info("No project information cached yet.")
		logger.Info("Run 'psx fix' or 'psx check' to create cache.")
		return nil
	}

	f := flags.GetFlags()

	if !f.GlobalFlags.Quiet {
		fmt.Println("Project Information:")
		fmt.Println()
		fmt.Printf("  Name:        %s\n", info.Name)
		fmt.Printf("  Description: %s\n", info.Description)
		fmt.Printf("  Author:      %s\n", info.Author)
		fmt.Printf("  Email:       %s\n", info.Email)
		fmt.Printf("  License:     %s\n", info.License)
		fmt.Println()
		if info.GitHubUser != "" {
			fmt.Printf("  GitHub User: %s\n", info.GitHubUser)
			fmt.Printf("  Repo Name:   %s\n", info.RepoName)
			fmt.Printf("  Repo URL:    %s\n", info.RepoURL)
			fmt.Println()
		}
		fmt.Printf("  Cache file:  %s\n", filepath.Join(pathCtx.Abs, ".psx-project.yml"))
	}

	return nil
}

func runProjectReset(cmd *cobra.Command, args []string) error {
	pathCtx, err := cmdctx.ResolvePath(args)
	if err != nil {
		return err
	}

	cachePath := filepath.Join(pathCtx.Abs, ".psx-project.yml")

	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		logger.Info("No project cache to reset.")
		return nil
	}

	if err := os.Remove(cachePath); err != nil {
		return logger.Errorf("Failed to remove cache: %w", err)
	}

	logger.Success("Project information cache reset.")
	logger.Info("You will be asked for project information next time.")

	return nil
}

func runProjectEdit(cmd *cobra.Command, args []string) error {
	pathCtx, err := cmdctx.ResolvePath(args)
	if err != nil {
		return err
	}

	existing, _ := resources.LoadProjectInfo(pathCtx.Abs)

	var info *resources.ProjectInfo
	if existing != nil {
		logger.Info("Current information loaded. Press Enter to keep current value.")
		fmt.Println()
		info = existing
		info.CurrentDir = pathCtx.Abs
		info.PromptUser()
	} else {
		info = resources.GetProjectInfo(pathCtx.Abs, true)
	}

	// Save
	if err := resources.SaveProjectInfo(pathCtx.Abs, info); err != nil {
		return logger.Errorf("Failed to save project info: %w", err)
	}

	logger.Success("Project information updated.")

	return nil
}
