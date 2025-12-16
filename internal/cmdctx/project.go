package cmdctx

import (
	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/detector"
	"github.com/m-mdy-m/psx/internal/flags"
	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/resources"
)

func LoadProject(args []string) (*ProjectContext, error) {
	pathCtx, err := ResolvePath(args)
	if err != nil {
		return nil, err
	}

	f := flags.GetFlags()

	logger.Verbose("Analyzing project...")
	logger.Verbosef("Path: %s", pathCtx.Abs)

	logger.Verbose("Loading configuration...")
	cfg, err := config.Load(f.GlobalFlags.ConfigFile, pathCtx.Abs)
	if err != nil {
		return nil, logger.Errorf("config load failed: %w", err)
	}

	logger.Verbose("Detecting project type...")
	var detection *detector.DetectionResult

	if cfg.Project.Type != "" {
		detection, err = detector.DetectWithHint(pathCtx.Abs, cfg.Project.Type)
		if err != nil {
			detection, err = detector.Detect(pathCtx.Abs)
			if err != nil {
				return nil, logger.Errorf("detection failed: %w", err)
			}
		}
	} else {
		detection, err = detector.Detect(pathCtx.Abs)
		if err != nil {
			return nil, logger.Errorf("detection failed: %w", err)
		}
	}

	cfg.Project.Type = detection.Type.Primary
	logger.Verbosef("Detected: %s", detection.Type.Primary)

	interactive := false
	if f.Fix.Interactive || f.Fix.All {
		interactive = true
	}

	projectInfo := resources.GetProjectInfo(pathCtx.Abs, interactive)

	if projectInfo == nil {
		logger.Warning("Could not get project info, using defaults")
		projectInfo = &resources.ProjectInfo{
			Name:       detection.Type.Primary,
			CurrentDir: pathCtx.Abs,
		}
	}

	return &ProjectContext{
		Path:        pathCtx,
		Config:      cfg,
		Detection:   detection,
		Flags:       f,
		ProjectInfo: projectInfo,
	}, nil
}
