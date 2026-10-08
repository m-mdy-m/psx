// Package detect infers a project's language and archetype from its layout.
package detect

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/m-mdy-m/psx/internal/resources"
	"github.com/m-mdy-m/psx/internal/tree"
)

// Profile is the outcome of a detection pass.
type Profile struct {
	ProjectType string                     `json:"project_type"`
	Kind        string                     `json:"kind"`
	Confidence  float64                    `json:"confidence"`
	Signals     []string                   `json:"signals"`
	Workspace   []string                   `json:"workspace,omitempty"`
	Managers    []resources.PackageManager `json:"package_managers,omitempty"`
}

// marker is a file whose presence indicates a language.
type marker struct {
	file  string
	lang  string
	score float64
}

// markers are ordered by descending score; the first match wins.
var markers = []marker{
	{"tsconfig.json", "nodejs", 0.98},
	{"go.mod", "go", 0.97},
	{"package.json", "nodejs", 0.95},
	{"Cargo.toml", "rust", 0.95},
	{"pyproject.toml", "python", 0.93},
	{"composer.json", "php", 0.9},
	{"Gemfile", "ruby", 0.9},
	{"pom.xml", "java", 0.9},
	{"build.gradle", "java", 0.85},
	{"deno.json", "nodejs", 0.8},
	{"CMakeLists.txt", "c", 0.8},
	{"setup.py", "python", 0.75},
	{"requirements.txt", "python", 0.7},
}

// workspaceMarkers indicate a multi-package repository.
var workspaceMarkers = []struct {
	file string
	kind string
}{
	{"pnpm-workspace.yaml", "monorepo"},
	{"lerna.json", "monorepo"},
	{"nx.json", "monorepo"},
	{"turbo.json", "monorepo"},
	{"go.work", "monorepo"},
	{"rush.json", "monorepo"},
}

// kindMarkers indicate a project archetype from its directory layout.
var kindMarkers = []struct {
	dir  string
	kind string
}{
	{"services", "microservice"},
	{"apps", "monorepo"},
	{"packages", "monorepo"},
	{"crates", "monorepo"},
	{"extensions", "plugin"},
	{"plugins", "plugin"},
	{"addons", "plugin"},
	{"cmd", "cli"},
}

// Detect inspects a project directory and reports its profile.
func Detect(root string, ignore []string) (*Profile, error) {
	snap, err := tree.Scan(root, ignore)
	if err != nil {
		return nil, err
	}

	// files and dirs are kept apart: a directory called go.mod is not a
	// module manifest, and treating it as one misreports the project.
	files := make(map[string]bool, 256)
	dirs := make(map[string]bool, 64)
	for _, p := range snap.Paths() {
		if snap.HasDir(p) {
			dirs[p] = true
			continue
		}
		files[p] = true
	}

	p := &Profile{ProjectType: "generic", Kind: "app", Confidence: 0.4}

	for _, m := range markers {
		if files[m.file] {
			p.ProjectType = m.lang
			p.Confidence = m.score
			p.Signals = append(p.Signals, m.file)
			break
		}
	}

	for _, w := range workspaceMarkers {
		if files[w.file] {
			p.Kind = w.kind
			p.Signals = append(p.Signals, w.file)
			break
		}
	}
	if p.Kind != "monorepo" {
		for _, k := range kindMarkers {
			if dirs[k.dir] {
				p.Kind = k.kind
				p.Signals = append(p.Signals, k.dir+"/")
				break
			}
		}
	}

	if p.Kind == "app" && (files["main.go"] || hasPrefixDir(dirs, "cmd/")) {
		p.Kind = "cli"
		p.Signals = append(p.Signals, "cmd/")
	}
	if p.Kind == "app" && (files["index.js"] || files["index.ts"]) {
		p.Kind = "library"
		p.Signals = append(p.Signals, "index entry point")
	}

	p.Workspace = workspaceGlobs(p)
	p.Managers = resources.ManagersFor(p.ProjectType, files)
	return p, nil
}

func hasPrefixDir(dirs map[string]bool, prefix string) bool {
	for d := range dirs {
		if strings.HasPrefix(d, prefix) {
			return true
		}
	}
	return false
}

func workspaceGlobs(p *Profile) []string {
	if p.Kind != "monorepo" && p.Kind != "microservice" && p.Kind != "plugin" {
		return nil
	}

	var globs []string
	switch p.ProjectType {
	case "go":
		globs = []string{"./..."}
	case "rust":
		globs = []string{"crates/*"}
	case "nodejs":
		globs = []string{"packages/*", "apps/*"}
	default:
		globs = []string{"packages/*"}
	}

	switch p.Kind {
	case "microservice":
		globs = append(globs, "services/*")
	case "plugin":
		globs = append(globs, "extensions/*", "plugins/*")
	}
	sort.Strings(globs)
	return globs
}

func Kinds() []string {
	return []string{"app", "library", "cli", "monorepo", "microservice", "plugin"}
}

func Known(kind string) bool {
	for _, k := range Kinds() {
		if k == kind {
			return true
		}
	}
	return false
}

// DetectFromPath is a convenience wrapper for callers holding a directory.
func DetectFromPath(root string) (*Profile, error) {
	if _, err := os.Stat(root); err != nil {
		return nil, err
	}
	return Detect(filepath.Clean(root), nil)
}
