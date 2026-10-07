package rules

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/m-mdy-m/psx/internal/config"
)

// CustomFileSpec is a user-declared file to create.
type CustomFileSpec = config.CustomFile

// CustomFolderSpec is a user-declared folder to create.
type CustomFolderSpec = config.CustomFolder

// CustomHandler applies the user-declared files and folders from a config.
//
// A configuration file is untrusted input, so every target path is resolved
// against the project root and rejected if it escapes. Without that check a
// cloned repository could make `psx fix --yes` write anywhere on the machine.
type CustomHandler struct {
	root string
}

func (h CustomHandler) Apply(cfg *config.CustomConfig, fixCtx *FixContext) []*FixResult {
	if cfg == nil || h.root == "" {
		return nil
	}

	var out []*FixResult
	out = append(out, h.applyFiles(cfg.Files, fixCtx)...)
	out = append(out, h.applyFolders(cfg.Folders, fixCtx)...)
	return out
}

func (h CustomHandler) applyFiles(files []CustomFileSpec, fixCtx *FixContext) []*FixResult {
	out := make([]*FixResult, 0, len(files))

	for _, f := range files {
		res := &FixResult{RuleID: "custom:file:" + f.Path}

		full, err := h.resolve(f.Path)
		if err != nil {
			res.Error = err
			out = append(out, res)
			continue
		}

		if exists, info := statPath(full); exists && info.Size() > 0 {
			res.Skipped = true
			res.Reason = "file already has content"
			out = append(out, res)
			continue
		}

		res.Changes = []Change{{
			Type:        ChangeCreateFile,
			Path:        full,
			Description: "created " + f.Path,
			Mode:        0o644,
		}}

		if fixCtx != nil && fixCtx.DryRun {
			res.Fixed = true
			res.Changes[0].Description = "create " + f.Path
			res.Changes[0].Content = Preview(f.Content, 10)
			out = append(out, res)
			continue
		}

		if err := writeFile(full, f.Content, 0o644); err != nil {
			res.Error = err
			out = append(out, res)
			continue
		}
		res.Fixed = true
		out = append(out, res)
	}
	return out
}

func (h CustomHandler) applyFolders(folders []CustomFolderSpec, fixCtx *FixContext) []*FixResult {
	out := make([]*FixResult, 0, len(folders))

	for _, f := range folders {
		res := &FixResult{RuleID: "custom:folder:" + f.Path}

		if _, err := h.resolve(f.Path); err != nil {
			res.Error = err
			out = append(out, res)
			continue
		}

		paths := flattenStructure(f.Path, f.Structure)
		res.Changes = make([]Change, 0, len(paths))

		for _, rel := range paths {
			child, err := h.resolve(rel)
			if err != nil {
				res.Error = err
				break
			}
			res.Changes = append(res.Changes, Change{
				Type:        ChangeCreateFolder,
				Path:        child,
				Description: "created " + rel,
			})
		}

		if res.Error != nil {
			out = append(out, res)
			continue
		}

		if fixCtx == nil || !fixCtx.DryRun {
			for _, rel := range paths {
				if err := os.MkdirAll(mustResolve(h.root, rel), 0o755); err != nil {
					res.Error = err
					break
				}
			}
		}
		res.Fixed = res.Error == nil
		out = append(out, res)
	}
	return out
}

func (h CustomHandler) resolve(rel string) (string, error) {
	if strings.TrimSpace(rel) == "" {
		return "", fmt.Errorf("empty path")
	}
	if filepath.IsAbs(rel) || strings.HasPrefix(rel, "/") || strings.HasPrefix(rel, `\`) {
		return "", fmt.Errorf("path %q must be relative to the project", rel)
	}

	full := filepath.Join(h.root, filepath.FromSlash(rel))
	inside, err := insideRoot(h.root, full)
	if err != nil {
		return "", err
	}
	if !inside {
		return "", fmt.Errorf("path %q escapes the project directory", rel)
	}
	return full, nil
}

// insideRoot reports whether target lies within root after symlink resolution.
func insideRoot(root, target string) (bool, error) {
	realRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		realRoot = root
	}

	// Resolve the deepest existing ancestor, since the target may not exist yet.
	probe := target
	for {
		if _, err := os.Lstat(probe); err == nil {
			break
		}
		parent := filepath.Dir(probe)
		if parent == probe {
			break
		}
		probe = parent
	}
	realProbe, err := filepath.EvalSymlinks(probe)
	if err != nil {
		realProbe = probe
	}

	rel, err := filepath.Rel(realRoot, realProbe)
	if err != nil {
		return false, err
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)), nil
}

// mustResolve is resolve for paths already validated by the caller.
func mustResolve(root, rel string) string {
	return filepath.Join(root, filepath.FromSlash(rel))
}

// flattenStructure expands a nested structure map into a sorted path list, so
// creation order is deterministic across runs.
func flattenStructure(base string, structure map[string]any) []string {
	var out []string
	var walk func(prefix string, node map[string]any)
	walk = func(prefix string, node map[string]any) {
		keys := make([]string, 0, len(node))
		for k := range node {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			path := k
			if prefix != "" {
				path = prefix + "/" + k
			}
			out = append(out, path)
			if child, ok := node[k].(map[string]any); ok {
				walk(path, child)
			}
		}
	}
	walk(base, structure)
	return out
}

func statPath(path string) (bool, os.FileInfo) {
	info, err := os.Stat(path)
	if err != nil {
		return false, nil
	}
	return true, info
}

func writeFile(path, content string, mode uint32) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create parent of %s: %w", path, err)
	}
	if err := os.WriteFile(path, []byte(content), os.FileMode(mode)); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
