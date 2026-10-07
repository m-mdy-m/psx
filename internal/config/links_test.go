package config

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Every relative markdown link must resolve. A newcomer following a broken link
// hits a dead end with no way to know they mistyped.
var linkPattern = regexp.MustCompile(`\[[^\]]*\]\(([^)]+)\)`)

func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for range 10 {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	t.Fatal("could not find the repository root")
	return ""
}

// markdownFiles lists every doc that ships with the project.
func markdownFiles(t *testing.T) []string {
	t.Helper()
	root := repoRoot(t)

	var out []string
	for _, dir := range []string{"docs", "examples"} {
		entries, err := os.ReadDir(filepath.Join(root, dir))
		if err != nil {
			t.Fatalf("reading %s: %v", dir, err)
		}
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".md") {
				out = append(out, filepath.Join(root, dir, e.Name()))
			}
		}
	}
	out = append(out, filepath.Join(root, "readme.md"))
	return out
}

func TestEveryRelativeLinkResolves(t *testing.T) {
	root := repoRoot(t)

	for _, file := range markdownFiles(t) {
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("reading %s: %v", file, err)
		}

		rel, _ := filepath.Rel(root, file)
		for _, m := range linkPattern.FindAllStringSubmatch(string(data), -1) {
			target := strings.TrimSpace(m[1])
			if idx := strings.IndexAny(target, "#"); idx >= 0 {
				target = target[:idx]
			}
			switch {
			case target == "", strings.HasPrefix(target, "http://"),
				strings.HasPrefix(target, "https://"), strings.HasPrefix(target, "mailto:"):
				continue
			}

			resolved := filepath.Join(filepath.Dir(file), filepath.FromSlash(target))
			if _, err := os.Stat(resolved); err != nil {
				t.Errorf("%s links to %q, which does not exist", rel, m[1])
			}
		}
	}
}

// The files named in the README's documentation index have to exist, because that
// list is the only map a newcomer gets of the whole doc set.
func TestReadmePointsAtEveryGuide(t *testing.T) {
	data, err := os.ReadFile(filepath.Join(repoRoot(t), "readme.md"))
	if err != nil {
		t.Fatal(err)
	}

	readme := string(data)
	for _, want := range []string{
		"docs/GETTING-STARTED.md",
		"docs/CONFIGURATION.md",
		"docs/RULES.md",
		"docs/WORKFLOWS.md",
		"examples/README.md",
		"docs/INSTALLATION.md",
	} {
		if !strings.Contains(readme, want) {
			t.Errorf("readme.md never mentions %s", want)
		}
	}
}
