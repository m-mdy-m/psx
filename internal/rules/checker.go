package rules

import "github.com/m-mdy-m/psx/internal/tree"

type Checker struct {
	snap *tree.Snapshot
}

func NewChecker(snap *tree.Snapshot) *Checker {
	return &Checker{snap: snap}
}

func (c *Checker) CheckAny(patterns []string) (string, bool) {
	if c == nil || c.snap == nil {
		return "", false
	}
	for _, p := range patterns {
		if c.check(p) {
			return p, true
		}
	}
	return "", false
}

// check resolves one pattern: a trailing slash demands a directory, a glob
// matches recursively, and everything else must exist with content.
func (c *Checker) check(pattern string) bool {
	if pattern == "" {
		return false
	}
	if isDirPattern(pattern) {
		return c.snap.NonEmpty(pattern)
	}
	if hasGlobMeta(pattern) {
		return len(c.snap.Glob(pattern)) > 0
	}
	return c.snap.HasFile(pattern)
}
