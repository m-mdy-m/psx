package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/m-mdy-m/psx/internal/cmdctx"
	"github.com/m-mdy-m/psx/internal/flags"
	"github.com/m-mdy-m/psx/internal/logger"
)

// newInitCmd builds the init command, which writes a starter configuration.
func newInitCmd() *cobra.Command {
	opts := baseOptions()

	cmd := &cobra.Command{
		Use:   "init [path]",
		Short: "Create a psx configuration file",
		Long: "Write a psx.yml for the project.\n\n" +
			"The project type is inferred from the layout unless --type is given, and the\n" +
			"generated file enables the rules that make sense for that profile.",
		Args:          cobra.MaximumNArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			applyGlobals(&opts)
			root := "."
			if len(args) > 0 {
				root = args[0]
			}
			return runInit(root, &opts)
		},
	}

	f := cmd.Flags()
	f.StringVar(&opts.Init.ProjectType, "type", "", "project type: nodejs, go, rust, python, generic")
	f.StringVar(&opts.Init.Kind, "kind", "", "project kind: app, library, cli, monorepo, microservice, plugin")
	f.StringVar(&opts.Init.License, "license", "MIT", "license to record")
	f.BoolVar(&opts.Init.Force, "force", false, "overwrite an existing configuration")
	f.BoolVar(&opts.Init.Minimal, "minimal", false, "write only the essential rules")

	return cmd
}

func runInit(root string, opts *flags.Options) error {
	path := "psx.yml"
	if existing := findConfigFile(root); existing != "" {
		path = existing
		if !opts.Init.Force {
			return commandError(exitArgs,
				"%s already exists; pass --force to overwrite", path)
		}
	}

	abs, err := filepath.Abs(filepath.Join(root, path))
	if err != nil {
		return commandError(exitConfig, "%v", err)
	}

	projectType := opts.Init.ProjectType
	kind := opts.Init.Kind
	if projectType == "" {
		det, err := detectProfile(root, opts.Global.ConfigFile)
		if err != nil {
			return err
		}
		projectType = det.ProjectType
		if kind == "" {
			kind = det.Kind
		}
		logger.Verbosef("inferred type=%s kind=%s", projectType, kind)
	}

	content := renderDefaultConfig(projectType, kind, opts.Init.License, opts.Init.Minimal)
	if err := writeFileIfAbsent(abs, content, opts.Init.Force); err != nil {
		return commandError(exitCode, "%v", err)
	}

	logger.Successf("Created %s", abs)
	logger.Info("Edit it to match your project, then run `psx check`")
	return nil
}

// findConfigFile locates an existing configuration in the project.
func findConfigFile(root string) string {
	for _, name := range []string{"psx.yml", ".psx.yml", "psx.yaml", ".psx.yaml"} {
		p := filepath.Join(root, name)
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			return p
		}
	}
	return ""
}

func detectProfile(root, configFile string) (*detectResult, error) {
	ctx, err := cmdctx.Load(root, configFile, cmdctx.ModeReadOnly)
	if err != nil {
		return nil, err
	}
	return &detectResult{
		ProjectType: ctx.ProjectType,
		Kind:        ctx.Kind,
	}, nil
}

// detectResult is the minimal profile information init needs.
type detectResult struct {
	ProjectType string
	Kind        string
}

// essentialRules are always written by init, since a project is not usable
// without them.
var essentialRules = []string{"readme", "license", "gitignore", "gitattributes"}

// libraryRules fit a package or library, where a container adds little.
var libraryRules = []string{"api_docs", "changelog", "tests_folder", "contributing", "editorconfig"}

// serviceRules fit something that runs as a service.
var serviceRules = []string{"ci_config", "dockerfile", "dockerignore", "env_example", "dependency_updates"}

// monorepoRules apply to every member of a workspace.
var monorepoRules = []string{"docs_folder", "adr", "security", "pull_request_template", "issue_templates"}

// cliRules fit a command line tool.
var cliRules = []string{"makefile", "scripts_folder", "release_workflow", "lockfile"}

// microserviceRules fit a deployed service.
var microserviceRules = []string{"docker_compose", "openapi", "architecture", "runbook"}

// pluginRules fit an extension host project.
var pluginRules = []string{"architecture", "adr", "runbook", "docs_folder"}

// defaultSeverity is the severity assigned to a generated rule entry.
func defaultSeverity(id string) string {
	if id == "readme" || id == "tests_folder" {
		return "error"
	}
	if id == "license" || id == "env_example" || id == "dependency_updates" {
		return "warning"
	}
	return "info"
}

// renderDefaultConfig produces a starter configuration for a profile.
func renderDefaultConfig(projectType, kind, license string, minimal bool) string {
	rules := append([]string{}, essentialRules...)
	if !minimal {
		rules = append(rules, pickRulesForKind(kind)...)
	}
	rules = dedupe(rules)

	var b strings.Builder
	b.WriteString("# psx configuration\n")
	b.WriteString("# Docs: https://github.com/m-mdy-m/psx\n\n")
	b.WriteString("version: 1\n\n")
	b.WriteString("project:\n")
	b.WriteString(fmt.Sprintf("  type: %s\n", projectType))
	if kind != "" {
		b.WriteString(fmt.Sprintf("  kind: %s\n", kind))
	}
	b.WriteString("\nrules:\n")
	for _, id := range rules {
		b.WriteString(fmt.Sprintf("  %s: %s\n", id, defaultSeverity(id)))
	}

	if !minimal {
		b.WriteString("\n# Paths excluded from every check and fix.\n")
		b.WriteString("ignore:\n")
		for _, p := range defaultIgnore {
			b.WriteString(fmt.Sprintf("  - %s\n", p))
		}
	}

	b.WriteString("\nfix:\n")
	b.WriteString("  interactive: true\n")
	b.WriteString("  backup: false\n")

	if license != "" {
		b.WriteString(fmt.Sprintf("\n# License: %s\n", license))
	}
	return b.String()
}

// pickRulesForKind selects the rule set for a project archetype.
func pickRulesForKind(kind string) []string {
	switch kind {
	case "library", "package":
		return append(libraryRules, monorepoRules...)
	case "cli", "tool":
		return append(cliRules, monorepoRules...)
	case "microservice", "service":
		return append(serviceRules, microserviceRules...)
	case "monorepo", "workspace":
		return append(monorepoRules, serviceRules...)
	case "plugin", "editor":
		return append(pluginRules, serviceRules...)
	default:
		return append(serviceRules, monorepoRules...)
	}
}

func dedupe(list []string) []string {
	seen := make(map[string]bool, len(list))
	out := make([]string, 0, len(list))
	for _, v := range list {
		if seen[v] {
			continue
		}
		seen[v] = true
		out = append(out, v)
	}
	return out
}

// defaultIgnore is the ignore list written into a generated configuration.
var defaultIgnore = []string{
	"node_modules/", "vendor/", ".git/", "dist/", "build/", "coverage/", ".psx/",
}
