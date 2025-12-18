package rules

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func New(ctx *Context) *Checker {
	return &Checker{ctx: ctx}
}

func (c *Checker) CheckAll(rules []Rule) []CheckResult {
	results := make([]CheckResult, len(rules))
	var wg sync.WaitGroup

	for i, rule := range rules {
		wg.Add(1)
		go func(idx int, r Rule) {
			defer wg.Done()
			results[idx] = c.Check(r)
		}(i, rule)
	}

	wg.Wait()
	return results
}

func (c *Checker) Check(rule Rule) CheckResult {
	patterns := c.getPatterns(rule.Patterns)
	if len(patterns) == 0 {
		return c.failResult(rule, "No patterns defined for this project type")
	}
	for _, pattern := range patterns {
		fullPath := filepath.Join(c.ctx.ProjectPath, pattern)

		var exists bool
		switch rule.Type {
		case RuleTypeFile:
			exists = c.checkFile(fullPath)
		case RuleTypeFolder:
			exists = c.checkFolder(fullPath)
		}

		if exists {
			return c.successResult(rule, pattern)
		}
	}

	return c.failResult(rule, "Not found")
}

func (c *Checker) checkFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	if info.IsDir() {
		return false
	}

	return info.Size() > 0
}

func (c *Checker) checkFolder(path string) bool {
	if strings.Contains(path, "*") {
		matches, err := filepath.Glob(path)
		return err == nil && len(matches) > 0
	}

	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func (c *Checker) getPatterns(patterns map[string][]string) []string {
	if p, ok := patterns[c.ctx.ProjectType]; ok {
		return p
	}

	if p, ok := patterns["*"]; ok {
		return p
	}

	return nil
}
