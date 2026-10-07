package rules

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"

	"github.com/m-mdy-m/psx/internal/config"
	"github.com/m-mdy-m/psx/internal/resources"
)

// yamlExtensions are the file types that must parse as YAML.
var yamlExtensions = map[string]bool{".yml": true, ".yaml": true}

// TestGeneratedYAMLIsValid parses every generated YAML file.
//
// GitHub silently ignores a malformed workflow or funding file, so a broken
// template produces no visible error until the feature simply does not work.
// Parsing here turns that into a test failure.
func TestGeneratedYAMLIsValid(t *testing.T) {
	env := newFixEnv(t, map[string]string{}, "go", nil)

	meta := config.GetRulesMetadata().Rules
	var ids []string
	for id, m := range meta {
		if m.Fixable() {
			ids = append(ids, id)
		}
	}

	env.cfg = testConfig(t, ids, nil)
	env.ctx.Config = env.cfg
	env.run(t, ids, false, false)

	checked := 0
	err := filepath.WalkDir(env.root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if !yamlExtensions[strings.ToLower(filepath.Ext(path))] {
			return nil
		}

		data, rerr := os.ReadFile(path)
		if rerr != nil {
			return rerr
		}
		checked++

		// safe_load_all accepts the multi-document manifests that Kubernetes
		// requires, so use the streaming decoder.
		dec := yaml.NewDecoder(strings.NewReader(string(data)))
		for {
			var doc any
			derr := dec.Decode(&doc)
			if derr != nil {
				if strings.Contains(derr.Error(), "EOF") {
					return nil
				}
				rel, _ := filepath.Rel(env.root, path)
				t.Errorf("%s is not valid YAML: %v", rel, derr)
				return nil
			}
		}
	})

	if err != nil {
		t.Fatalf("walk: %v", err)
	}
	if checked == 0 {
		t.Fatal("no YAML files were generated; the test is not exercising the templates")
	}
}

// TestFundingTemplateUsesGitHubSchema pins the FUNDING.yml shape.
//
// GitHub reads this file as structured data keyed by platform. A markdown list
// parses as YAML but is ignored, which is why this needed its own check.
func TestFundingTemplateUsesGitHubSchema(t *testing.T) {
	body, err := resources.Template(resources.TmplFunding, "generic", map[string]string{
		"github_username": "ada",
		"repo_name":       "demo",
	})
	if err != nil {
		t.Fatal(err)
	}

	var doc struct {
		GitHub []string `yaml:"github"`
	}
	if err := yaml.Unmarshal([]byte(body), &doc); err != nil {
		t.Fatalf("funding template is not valid YAML: %v\n%s", err, body)
	}
	if len(doc.GitHub) == 0 {
		t.Errorf("funding template has no github key: %s", body)
	}
	if strings.Contains(body, "This project is community supported") {
		t.Error("funding template contains prose; GitHub only reads platform keys")
	}
}

// TestGeneratedFilesHaveTrailingNewline keeps generated files POSIX-clean, so a
// missing newline does not show up as a diff marker in review.
func TestGeneratedFilesHaveTrailingNewline(t *testing.T) {
	env := newFixEnv(t, map[string]string{}, "nodejs",
		[]string{"readme", "license", "changelog", "editorconfig", "env_example", "gitignore"})

	env.run(t, []string{"readme", "license", "changelog", "editorconfig", "env_example", "gitignore"}, false, false)

	err := filepath.WalkDir(env.root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		data, rerr := os.ReadFile(path)
		if rerr != nil {
			return rerr
		}
		if len(data) == 0 {
			return nil
		}
		if data[len(data)-1] != '\n' {
			rel, _ := filepath.Rel(env.root, path)
			t.Errorf("%s does not end with a newline", rel)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
