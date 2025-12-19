package cmdctx

import (
	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/flags"
	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/resources"
)

type ProjectContext struct {
	Path        *PathContext
	Config      *config.Config
	ProjectType string
	Flags       *flags.Flags
	ProjectInfo *resources.ProjectInfo
}

type PathContext struct {
	Root string
	Abs  string
}

func LoadProject(args []string) (*ProjectContext, error) {
	pathCtx, err := ResolvePath(args)
	if err != nil {
		return nil, err
	}

	f := flags.GetFlags()

	logger.Verbose("Analyzing project...")
	logger.Verbosef("Path: %s", pathCtx.Abs)

	// Load configuration
	logger.Verbose("Loading configuration...")
	cfg, err := config.Load(f.GlobalFlags.ConfigFile, pathCtx.Abs)
	if err != nil {
		return nil, logger.Errorf("config load failed: %w", err)
	}

	// Project type is now determined from config or normalized
	projectType := cfg.Project.Type
	if projectType == "" {
		projectType = "generic"
	}

	logger.Verbosef("Project type: %s", projectType)
	logger.Verbosef("Active rules: %d", len(cfg.ActiveRules))

	// Get project info
	interactive := false
	if f.Fix.Interactive || f.Fix.All {
		interactive = true
	}

	projectInfo := resources.GetProjectInfo(pathCtx.Abs, interactive)
	if projectInfo == nil {
		logger.Warning("Could not get project info, using defaults")
		projectInfo = &resources.ProjectInfo{
			Name:   "project",
			Author: "Your Name",
			Email:  "you@example.com",
		}
	}

	return &ProjectContext{
		Path:        pathCtx,
		Config:      cfg,
		ProjectType: projectType,
		Flags:       f,
		ProjectInfo: projectInfo,
	}, nil
}
