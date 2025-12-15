package cmdctx

import (
	"os"
	"path/filepath"

	"github.com/m-mdy-m/psx/internal/detector"
	"github.com/m-mdy-m/psx/internal/flags"
	"github.com/m-mdy-m/psx/internal/logger"
)

func PrepareInit(args []string) (*InitContext, error) {
	pathCtx, err := ResolvePath(args)
	if err != nil {
		return nil, err
	}

	f := flags.GetFlags()
	configPath := filepath.Join(pathCtx.Abs, "psx.yml")

	if _, err := os.Stat(configPath); err == nil && !f.Init.Force {
		return nil, logger.Errorf("configuration already exists (use --force)")
	}

	projectType := f.Init.Template
	if projectType == "" {
		logger.Verbose("Detecting project type...")
		detection, err := detector.Detect(pathCtx.Abs)
		if err != nil {
			projectType = "generic"
		} else {
			projectType = detection.Type.Primary
		}
	}

	return &InitContext{
		Path:        pathCtx,
		ProjectType: projectType,
		Flags:       f,
		ConfigPath:  configPath,
	}, nil
}
