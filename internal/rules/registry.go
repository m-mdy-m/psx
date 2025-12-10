package rules

import (
	"fmt"
	"sync"
)

// Global executor registry
var (
	executorRegistry = &Registry{
		executors: make(map[string]Executor),
	}
	registryOnce sync.Once
)

// Registry manages rule executors
type Registry struct {
	mu        sync.RWMutex
	executors map[string]Executor
}

// Register adds an executor for a rule
func Register(ruleID string, executor Executor) {
	executorRegistry.mu.Lock()
	defer executorRegistry.mu.Unlock()
	executorRegistry.executors[ruleID] = executor
}

// RegisterFunc adds a function executor for a rule
func RegisterFunc(ruleID string, fn ExecutorFunc) {
	Register(ruleID, fn)
}

// GetExecutor retrieves an executor by rule ID
func GetExecutor(ruleID string) (Executor, error) {
	executorRegistry.mu.RLock()
	defer executorRegistry.mu.RUnlock()

	executor, exists := executorRegistry.executors[ruleID]
	if !exists {
		return nil, fmt.Errorf("no executor registered for rule: %s", ruleID)
	}

	return executor, nil
}

// HasExecutor checks if an executor exists for a rule
func HasExecutor(ruleID string) bool {
	executorRegistry.mu.RLock()
	defer executorRegistry.mu.RUnlock()
	_, exists := executorRegistry.executors[ruleID]
	return exists
}

// ListRegistered returns all registered rule IDs
func ListRegistered() []string {
	executorRegistry.mu.RLock()
	defer executorRegistry.mu.RUnlock()

	result := make([]string, 0, len(executorRegistry.executors))
	for id := range executorRegistry.executors {
		result = append(result, id)
	}
	return result
}

// init registers all built-in executors
func init() {
	registryOnce.Do(func() {
		// Register all built-in rule executors
		registerBuiltinExecutors()
	})
}

// registerBuiltinExecutors registers all built-in rule executors
func registerBuiltinExecutors() {
	// General rules
	RegisterFunc("readme", executeReadmeRule)
	RegisterFunc("license", executeLicenseRule)
	RegisterFunc("gitignore", executeGitignoreRule)
	RegisterFunc("changelog", executeChangelogRule)

	// Structure rules
	RegisterFunc("src_folder", executeSrcFolderRule)
	RegisterFunc("tests_folder", executeTestsFolderRule)
	RegisterFunc("docs_folder", executeDocsFolderRule)

	// Documentation rules
	RegisterFunc("adr", executeADRRule)
	RegisterFunc("contributing", executeContributingRule)
	RegisterFunc("api_docs", executeAPIDocsRule)

	// CI/CD rules
	RegisterFunc("ci_config", executeCIConfigRule)

	// Quality rules
	RegisterFunc("pre_commit", executePreCommitRule)
	RegisterFunc("editorconfig", executeEditorconfigRule)
	RegisterFunc("code_owners", executeCodeOwnersRule)
}
