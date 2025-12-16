package fixer

import (
	"github.com/m-mdy-m/psx/internal/resources"
	"github.com/m-mdy-m/psx/internal/shared"
)

func FixDockerfile(ctx *FixContext) (*FixResult, error) {
	content := resources.GetDockerfile(ctx.ProjectInfo, ctx.ProjectType)
	if content == "" {
		return &FixResult{RuleID: "dockerfile", Skipped: true}, nil
	}

	return FixSingleFile(ctx, FileFixSpec{
		RuleID:       "dockerfile",
		Path:         "Dockerfile",
		Description:  "Create Dockerfile",
		PromptText:   "Create Dockerfile?",
		FormatForDry: true,
		GetContent: func(ctx *FixContext) (string, error) {
			return content, nil
		},
	})
}

func FixDockerIgnore(ctx *FixContext) (*FixResult, error) {
	content := resources.GetDockerignore(ctx.ProjectType)
	if content == "" {
		return &FixResult{RuleID: "dockerignore", Skipped: true}, nil
	}

	return FixSingleFile(ctx, FileFixSpec{
		RuleID:       "dockerignore",
		Path:         ".dockerignore",
		Description:  "Create .dockerignore",
		PromptText:   "Create .dockerignore?",
		FormatForDry: true,
		GetContent: func(ctx *FixContext) (string, error) {
			return content, nil
		},
	})
}

func FixDockerCompose(ctx *FixContext) (*FixResult, error) {
	// Ask if user wants database
	withDB := false
	if ctx.Interactive && !ctx.DryRun {
		withDB = shared.Prompt("Include database service?")
	}

	return FixSingleFile(ctx, FileFixSpec{
		RuleID:       "docker_compose",
		Path:         "docker-compose.yml",
		Description:  "Create docker-compose.yml",
		PromptText:   "Create docker-compose.yml?",
		FormatForDry: true,
		GetContent: func(ctx *FixContext) (string, error) {
			return resources.GetDockerCompose(ctx.ProjectInfo, withDB), nil
		},
	})
}

func FixKubernetes(ctx *FixContext) (*FixResult, error) {
	return FixFolder(ctx, FolderFixSpec{
		RuleID:      "kubernetes",
		Path:        "k8s",
		Description: "Create k8s/ folder with deployment files",
		PromptText:  "Create Kubernetes configurations?",
		Files: []FolderFileSpec{
			{
				Name: "deployment.yml",
				GetContent: func(ctx *FixContext) (string, error) {
					return resources.GetKubernetesDeployment(ctx.ProjectInfo), nil
				},
				FormatForDry: true,
			},
			{
				Name: "service.yml",
				GetContent: func(ctx *FixContext) (string, error) {
					return resources.GetKubernetesService(ctx.ProjectInfo), nil
				},
				FormatForDry: true,
			},
			{
				Name: "ingress.yml",
				GetContent: func(ctx *FixContext) (string, error) {
					return resources.GetKubernetesIngress(ctx.ProjectInfo), nil
				},
				FormatForDry: true,
			},
		},
	})
}

func FixNginxConfig(ctx *FixContext) (*FixResult, error) {
	// Ask user which type of nginx config
	configType := "static"
	if ctx.Interactive && !ctx.DryRun {
		choice, _ := shared.PromptChoice(
			"Select Nginx configuration type:",
			[]string{"Static site", "Reverse proxy"},
		)
		if choice == "Reverse proxy" {
			configType = "reverse_proxy"
		}
	}

	return FixSingleFile(ctx, FileFixSpec{
		RuleID:       "nginx_config",
		Path:         "nginx.conf",
		Description:  "Create nginx.conf",
		PromptText:   "Create nginx.conf?",
		FormatForDry: true,
		GetContent: func(ctx *FixContext) (string, error) {
			return resources.GetNginxConfig(ctx.ProjectInfo, configType), nil
		},
	})
}

func FixInfraFolder(ctx *FixContext) (*FixResult, error) {
	return FixFolder(ctx, FolderFixSpec{
		RuleID:      "infra_folder",
		Path:        "infra",
		Description: "Create infra/ folder",
		PromptText:  "Create infrastructure folder?",
		Files: []FolderFileSpec{
			{
				Name: "README.md",
				GetContent: func(ctx *FixContext) (string, error) {
					return "# Infrastructure\n\nInfrastructure as Code configurations go here.\n", nil
				},
			},
		},
	})
}
