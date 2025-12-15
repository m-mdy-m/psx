package rules

import (
	"fmt"
	reg	"github.com/m-mdy-m/psx/internal/registry"
	"github.com/m-mdy-m/psx/internal/checker"
)

var rulesReg = reg.New[Executor]("rules")
func GetExecutor(ruleID string) (Executor, error) {
	executor, ok := rulesReg.Get(ruleID)
	if !ok {
		return nil, fmt.Errorf("no executor registered for rule: %s", ruleID)
	}
	return executor, nil
}

func HasExecutor(ruleID string) bool {
	_, ok := rulesReg.Get(ruleID)
	return ok
}

func ListRegistered() []string {
	all :=	rulesReg.All()
	result := make([]string, 0, len(all))
	for id := range all {
		result = append(result, id)
	}
	return result
}

func init() {
	registerBuiltinExecutors()
}

func registerBuiltinExecutors() {
    // General
    rulesReg.Add("readme", ExecutorFunc(checker.CheckReadmeRule))
    rulesReg.Add("license", ExecutorFunc(checker.CheckLicenseRule))
    rulesReg.Add("gitignore", ExecutorFunc(checker.CheckGitignoreRule))
    rulesReg.Add("changelog", ExecutorFunc(checker.CheckChangelogRule))
    
    // Structure
    rulesReg.Add("src_folder", ExecutorFunc(checker.CheckSrcFolderRule))
    rulesReg.Add("tests_folder", ExecutorFunc(checker.CheckTestsFolderRule))
    rulesReg.Add("docs_folder", ExecutorFunc(checker.CheckDocsFolderRule))
    rulesReg.Add("scripts_folder", ExecutorFunc(checker.CheckScriptsFolderRule))
    rulesReg.Add("env_example", ExecutorFunc(checker.CheckEnvExampleRule))
    
    // Documentation
    rulesReg.Add("adr", ExecutorFunc(checker.CheckADRRule))
    rulesReg.Add("contributing", ExecutorFunc(checker.CheckContributingRule))
    rulesReg.Add("api_docs", ExecutorFunc(checker.CheckAPIDocsRule))
    rulesReg.Add("security", ExecutorFunc(checker.CheckSecurityRule))
    rulesReg.Add("code_of_conduct", ExecutorFunc(checker.CheckCodeOfConductRule))
    rulesReg.Add("pull_request_template", ExecutorFunc(checker.CheckPullRequestTemplateRule))
    rulesReg.Add("issue_templates", ExecutorFunc(checker.CheckIssueTemplatesRule))
    rulesReg.Add("funding", ExecutorFunc(checker.CheckFundingRule))
    rulesReg.Add("support", ExecutorFunc(checker.CheckSupportRule))
    rulesReg.Add("roadmap", ExecutorFunc(checker.CheckRoadmapRule))
    
    // CI/CD
    rulesReg.Add("ci_config", ExecutorFunc(checker.CheckCIConfigRule))
    rulesReg.Add("github_actions", ExecutorFunc(checker.CheckGitHubActionsRule))
    rulesReg.Add("renovate", ExecutorFunc(checker.CheckRenovateRule))
    rulesReg.Add("dependabot", ExecutorFunc(checker.CheckDependabotRule))
    
    // Quality
    rulesReg.Add("pre_commit", ExecutorFunc(checker.CheckPreCommitRule))
    rulesReg.Add("editorconfig", ExecutorFunc(checker.CheckEditorconfigRule))
    rulesReg.Add("code_owners", ExecutorFunc(checker.CheckCodeOwnersRule))
    rulesReg.Add("prettier", ExecutorFunc(checker.CheckPrettierRule))
    rulesReg.Add("prettierignore", ExecutorFunc(checker.CheckPrettierIgnoreRule))
    rulesReg.Add("eslint", ExecutorFunc(checker.CheckESLintRule))
    rulesReg.Add("commitlint", ExecutorFunc(checker.CheckCommitlintRule))
    rulesReg.Add("husky", ExecutorFunc(checker.CheckHuskyRule))
    rulesReg.Add("lint_staged", ExecutorFunc(checker.CheckLintStagedRule))
    rulesReg.Add("makefile", ExecutorFunc(checker.CheckMakefileRule))
    rulesReg.Add("gitattributes", ExecutorFunc(checker.CheckGitattributesRule))
    
    // DevOps
    rulesReg.Add("dockerfile", ExecutorFunc(checker.CheckDockerfileRule))
    rulesReg.Add("dockerignore", ExecutorFunc(checker.CheckDockerIgnoreRule))
    rulesReg.Add("docker_compose", ExecutorFunc(checker.CheckDockerComposeRule))
    rulesReg.Add("kubernetes", ExecutorFunc(checker.CheckKubernetesRule))
    rulesReg.Add("nginx_config", ExecutorFunc(checker.CheckNginxConfigRule))
    rulesReg.Add("infra_folder", ExecutorFunc(checker.CheckInfraFolderRule))
}