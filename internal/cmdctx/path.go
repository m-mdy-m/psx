package cmdctx

import (
	"path/filepath"

	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/utils"
)

func ResolvePath(args []string) (*PathContext, error) {
	root := "."
	if len(args) > 0 {
		root = args[0]
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, logger.Errorf("Invalid path: %w", err)
	}
	exists, _ := utils.FileExists(abs)
	if !exists {
		return nil, logger.Errorf("path does not exists: %s", abs)
	}
	return &PathContext{
		Root: root,
		Abs:  abs,
	}, nil
}
