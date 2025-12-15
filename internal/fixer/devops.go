package fixer

import (
	"path/filepath"

	"github.com/m-mdy-m/psx/internal/resources"
	"github.com/m-mdy-m/psx/internal/shared"
)

func FixDockerfile(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "dockerfile",
		Changes: []Change{},
	}

	dockerfilePath := filepath.Join(ctx.ProjectPath, "Dockerfile")

	exists, _ := shared.FileExists(dockerfilePath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	content := resources.GetDockerfile(ctx.ProjectInfo, ctx.ProjectType)
	if content == "" {
		result.Skipped = true
		return result, nil
	}

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        dockerfilePath,
			Description: "Create Dockerfile",
			Content:     FormatContent(content, 10),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create Dockerfile?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(dockerfilePath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        dockerfilePath,
		Description: "Created Dockerfile",
	})

	return result, nil
}

func FixDockerIgnore(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "dockerignore",
		Changes: []Change{},
	}

	dockerignorePath := filepath.Join(ctx.ProjectPath, ".dockerignore")

	exists, _ := shared.FileExists(dockerignorePath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	content := resources.GetDockerignore(ctx.ProjectType)
	if content == "" {
		result.Skipped = true
		return result, nil
	}

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        dockerignorePath,
			Description: "Create .dockerignore",
			Content:     FormatContent(content, 10),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create .dockerignore?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(dockerignorePath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        dockerignorePath,
		Description: "Created .dockerignore",
	})

	return result, nil
}

func FixDockerCompose(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "docker_compose",
		Changes: []Change{},
	}

	composePath := filepath.Join(ctx.ProjectPath, "docker-compose.yml")

	exists, _ := shared.FileExists(composePath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	// Ask if user wants database
	withDB := false
	if ctx.Interactive {
		withDB = shared.Prompt("Include database service?")
	}

	content := resources.GetDockerCompose(ctx.ProjectInfo, withDB)

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        composePath,
			Description: "Create docker-compose.yml",
			Content:     FormatContent(content, 10),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create docker-compose.yml?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(composePath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        composePath,
		Description: "Created docker-compose.yml",
	})

	return result, nil
}

func FixKubernetes(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "kubernetes",
		Changes: []Change{},
	}

	k8sPath := filepath.Join(ctx.ProjectPath, "k8s")

	exists, _ := shared.FileExists(k8sPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFolder,
			Path:        k8sPath,
			Description: "Create k8s/ folder with deployment files",
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create Kubernetes configurations?") {
			result.Skipped = true
			return result, nil
		}
	}

	// Create k8s folder
	if err := shared.CreateDir(k8sPath); err != nil {
		result.Error = err
		return result, err
	}

	// Create deployment.yml
	deploymentPath := filepath.Join(k8sPath, "deployment.yml")
	deploymentContent := resources.GetKubernetesDeployment(ctx.ProjectInfo)
	shared.CreateFile(deploymentPath, deploymentContent)

	// Create service.yml
	servicePath := filepath.Join(k8sPath, "service.yml")
	serviceContent := resources.GetKubernetesService(ctx.ProjectInfo)
	shared.CreateFile(servicePath, serviceContent)

	// Create ingress.yml
	ingressPath := filepath.Join(k8sPath, "ingress.yml")
	ingressContent := resources.GetKubernetesIngress(ctx.ProjectInfo)
	shared.CreateFile(ingressPath, ingressContent)

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFolder,
		Path:        k8sPath,
		Description: "Created k8s/ folder with deployment, service, and ingress files",
	})

	return result, nil
}

func FixNginxConfig(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "nginx_config",
		Changes: []Change{},
	}

	nginxPath := filepath.Join(ctx.ProjectPath, "nginx.conf")

	exists, _ := shared.FileExists(nginxPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	// Ask user which type of nginx config
	configType := "static"
	if ctx.Interactive {
		choice, _ := shared.PromptChoice("Select Nginx configuration type:", []string{
			"Static site",
			"Reverse proxy",
		})
		if choice == "Reverse proxy" {
			configType = "reverse_proxy"
		}
	}

	content := resources.GetNginxConfig(ctx.ProjectInfo, configType)

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFile,
			Path:        nginxPath,
			Description: "Create nginx.conf",
			Content:     FormatContent(content, 10),
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create nginx.conf?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateFile(nginxPath, content); err != nil {
		result.Error = err
		return result, err
	}

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFile,
		Path:        nginxPath,
		Description: "Created nginx.conf",
	})

	return result, nil
}

func FixInfraFolder(ctx *FixContext) (*FixResult, error) {
	result := &FixResult{
		RuleID:  "infra_folder",
		Changes: []Change{},
	}

	infraPath := filepath.Join(ctx.ProjectPath, "infra")

	exists, _ := shared.FileExists(infraPath)
	if exists {
		result.Skipped = true
		return result, nil
	}

	if ctx.DryRun {
		result.Changes = append(result.Changes, Change{
			Type:        ChangeCreateFolder,
			Path:        infraPath,
			Description: "Create infra/ folder",
		})
		result.Fixed = true
		return result, nil
	}

	if ctx.Interactive {
		if !shared.Prompt("Create infrastructure folder?") {
			result.Skipped = true
			return result, nil
		}
	}

	if err := shared.CreateDir(infraPath); err != nil {
		result.Error = err
		return result, err
	}

	// Create README
	readmePath := filepath.Join(infraPath, "README.md")
	readmeContent := "# Infrastructure\n\nInfrastructure as Code configurations go here.\n"
	shared.CreateFile(readmePath, readmeContent)

	result.Fixed = true
	result.Changes = append(result.Changes, Change{
		Type:        ChangeCreateFolder,
		Path:        infraPath,
		Description: "Created infra/ folder",
	})

	return result, nil
}