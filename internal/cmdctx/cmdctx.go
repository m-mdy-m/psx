// Package cmdctx assembles the shared context every command needs.
package cmdctx

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/detect"
	"github.com/m-mdy-m/psx/internal/logger"
	"github.com/m-mdy-m/psx/internal/resources"
	"github.com/m-mdy-m/psx/internal/rules"
	"github.com/m-mdy-m/psx/internal/tree"
)

// Mode describes how much a command is allowed to do.
type Mode int

const (
	// ModeReadOnly never prompts and never writes to the project.
	ModeReadOnly Mode = iota
	// ModeInteractive may prompt for project metadata and may write.
	ModeInteractive
	// ModeDryRun may prompt for metadata but never writes project files.
	ModeDryRun
)

// ProjectContext is the resolved state for one command invocation.
type ProjectContext struct {
	// Path is the absolute project root.
	Path string
	// Config is the loaded and validated configuration.
	Config *config.Config
	// ProjectType is the normalized language key.
	ProjectType string
	// Kind is the project archetype, such as "app" or "monorepo".
	Kind string
	// Info describes the project for template rendering.
	Info *resources.ProjectInfo
	// Snapshot is the scanned project tree.
	Snapshot *tree.Snapshot
	Detected bool
}

// configFile overrides discovery when non-empty; this is how --config is honoured.
func Load(root, configFile string, mode Mode) (*ProjectContext, error) {
	abs, err := ResolvePath(root)
	if err != nil {
		return nil, err
	}

	cfg, err := config.Load(configFile, abs)
	if err != nil {
		return nil, err
	}

	ctx := &ProjectContext{Path: abs, Config: cfg}

	// An explicitly declared type is normalized and trusted; otherwise the
	// layout decides. DeclaredType is the raw value, because Project.Type has
	// already been normalised to "generic" by the time it gets here.
	declared := cfg.DeclaredType
	if declared != "" && !isDetectPlaceholder(declared) {
		ctx.ProjectType = resources.NormalizeProjectType(declared)
	} else {
		profile, err := detect.Detect(abs, cfg.Ignore)
		if err != nil {
			return nil, fmt.Errorf("detect project: %w", err)
		}
		ctx.ProjectType = profile.ProjectType
		ctx.Kind = profile.Kind
		ctx.Detected = true

		for _, signal := range profile.Signals {
			logger.Verbose("detect: " + signal)
		}
		logger.Verbosef("detected type %s (kind %s, confidence %.0f%%)",
			profile.ProjectType, profile.Kind, profile.Confidence*100)
	}

	if ctx.Kind == "" {
		ctx.Kind = "app"
	}
	cfg.Project.Type = ctx.ProjectType
	cfg.ProjectType = ctx.ProjectType

	info, err := resources.ResolveInfo(abs, mode != ModeReadOnly)
	if err != nil {
		return nil, err
	}
	ctx.Info = info

	snap, err := tree.Scan(abs, cfg.Ignore)
	if err != nil {
		return nil, fmt.Errorf("scan project: %w", err)
	}
	ctx.Snapshot = snap

	logger.Verbosef("project %s type=%s kind=%s rules=%d",
		abs, ctx.ProjectType, ctx.Kind, len(cfg.ActiveRules))

	return ctx, nil
}

func isDetectPlaceholder(v string) bool {
	switch strings.ToLower(v) {
	case "", "auto", "detect", "auto-detect":
		return true
	default:
		return false
	}
}

func (c *ProjectContext) RuleContext() *rules.Context {
	return &rules.Context{
		ProjectPath: c.Path,
		ProjectType: c.ProjectType,
		ProjectInfo: c.Info,
		Config:      c.Config,
	}
}

func ResolvePath(root string) (string, error) {
	if strings.TrimSpace(root) == "" {
		root = "."
	}

	abs, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("resolve path %q: %w", root, err)
	}

	info, err := os.Stat(abs)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("path does not exist: %s", abs)
		}
		return "", fmt.Errorf("stat %s: %w", abs, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("not a directory: %s", abs)
	}
	return abs, nil
}
