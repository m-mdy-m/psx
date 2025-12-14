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
	rulesReg.Add("readme", ExecutorFunc(checker.CheckReadmeRule))
	rulesReg.Add("license", ExecutorFunc(checker.CheckLicenseRule))
	rulesReg.Add("gitignore", ExecutorFunc(checker.CheckGitignoreRule))
	rulesReg.Add("changelog", ExecutorFunc(checker.CheckChangelogRule))

	// Structure rules
	rulesReg.Add("src_folder", ExecutorFunc(checker.CheckSrcFolderRule))
	rulesReg.Add("tests_folder", ExecutorFunc(checker.CheckTestsFolderRule))
	rulesReg.Add("docs_folder", ExecutorFunc(checker.CheckDocsFolderRule))

	// Documentation rules
	rulesReg.Add("adr", ExecutorFunc(checker.CheckADRRule))
	rulesReg.Add("contributing", ExecutorFunc(checker.CheckContributingRule))
	rulesReg.Add("api_docs", ExecutorFunc(checker.CheckAPIDocsRule))

	// CI/CD rules
	rulesReg.Add("ci_config", ExecutorFunc(checker.CheckCIConfigRule))

	// Quality rules
	rulesReg.Add("pre_commit", ExecutorFunc(checker.CheckPreCommitRule))
	rulesReg.Add("editorconfig", ExecutorFunc(checker.CheckEditorconfigRule))
	rulesReg.Add("code_owners", ExecutorFunc(checker.CheckCodeOwnersRule))
}
