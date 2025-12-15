package cmdctx

import (
	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/detector"
	"github.com/m-mdy-m/psx/internal/flags"
)

type PathContext struct {
	Root string
	Abs  string
}

type ProjectContext struct {
	Path      *PathContext
	Config    *config.Config
	Detection *detector.DetectionResult
	Flags     *flags.Flags
}


type InitContext struct {
	Path        *PathContext
	ProjectType string
	Flags       *flags.Flags
	ConfigPath  string
}

