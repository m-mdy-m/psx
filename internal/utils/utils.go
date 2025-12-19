package utils

import (
	"github.com/m-mdy-m/psx/internal/logger"
)

const (
	ExitSuccess = 0
	ExitFailed  = 1
	ExitConfig  = 2
	ExitFS      = 3
	ExitArgs    = 4
)

func Version(version string) {
	logger.Infof("PSX version %s\n", version)
}
