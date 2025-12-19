package rules

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/m-mdy-m/psx/internal/utils"
)



func NewChecker(ctx *Context) *Checker {
	return &Checker{ctx: ctx}
}

func (c *Checker) CheckAny(patterns []string) bool {
	for _, pattern := range patterns {
		if c.checkPattern(pattern) {
			return true
		}
	}
	return false
}

func (c *Checker) checkPattern(pattern string) bool {
	fullPath := filepath.Join(c.ctx.ProjectPath, pattern)

	if strings.Contains(pattern, "*") {
		return c.checkGlob(fullPath)
	}
	exists, info := utils.FileExists(fullPath)
	if !exists {
		return false
	}
	return c.validateContent(fullPath, info)
}

func (c *Checker) checkGlob(pattern string) bool {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return false
	}

	if len(matches) == 0 {
		return false
	}

	for _, match := range matches {
		exists, info := utils.FileExists(match)
		if exists && c.validateContent(match, info) {
			return true
		}
	}

	return false
}

func (c *Checker) validateContent(path string, info os.FileInfo) bool {
	if info.IsDir() {
		isEmpty, err := utils.IsDirEmpty(path)
		return err == nil && !isEmpty
	}
	return info.Size() > 0
}
