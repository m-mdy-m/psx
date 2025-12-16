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

	var detection *detector.DetectionResult
	cache, _ := resources.LoadProjectCache(pathCtx.Abs)

	usedCache := false
	if cache != nil && cache.Detection != nil && cache.Detection.ProjectType != "" {
		logger.Verbose("Using cached project detection")
		detection = &detector.DetectionResult{
			Type: detector.ProjectType{
				Primary:  cache.Detection.ProjectType,
				Version:  cache.Detection.Version,
				Features: cache.Detection.Features,
			},
			Files: cache.Detection.Files,
		}
		logger.Verbosef("Detected: %s (cached)", detection.Type.Primary)
		usedCache = true
	}

	if !usedCache {
		logger.Verbose("Detecting project type...")

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

		logger.Verbosef("Detected: %s", detection.Type.Primary)

		var projectInfo *resources.ProjectInfo
		if cache != nil && cache.ProjectInfo != nil {
			projectInfo = cache.ProjectInfo
		}

		err = resources.SaveProjectCache(
			pathCtx.Abs,
			projectInfo,
			detection.Type.Primary,
			detection.Type.Version,
			detection.Type.Features,
			detection.Files,
		)
		if err != nil {
			logger.Warning("Failed to save detection cache")
		} else {
			logger.Verbose("Detection cached successfully")
		}
	}

	cfg.Project.Type = detection.Type.Primary

	interactive := false
	if f.Fix.Interactive || f.Fix.All {
		interactive = true
	}

	projectInfo := resources.GetOrCreateProjectInfo(pathCtx.Abs, interactive)

	if projectInfo == nil {
		logger.Warning("Could not get project info, using defaults")
		projectInfo = &resources.ProjectInfo{
			Name:       detection.Type.Primary,
			CurrentDir: pathCtx.Abs,
		}
		projectInfo.SetDefaults()
		projectInfo.BuildDerivedFields()
	}

	return &ProjectContext{
		Path:        pathCtx,
		Config:      cfg,
		Detection:   detection,
		Flags:       f,
		ProjectInfo: projectInfo,
	}, nil
}
