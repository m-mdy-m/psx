package config

import (
	"fmt"
	"slices"
	"strings"
)

func IsValid( r ValidationResult )bool{
	return r.Valid && len(r.Errors) ==0
}

func HasWarnings(r ValidationResult) bool{
	return len(r.Warnings)>0
}

func Validate(c *Config) ValidationResult{
	result :=ValidationResult{
		Valid:		true,
		Errors:		[]ValidationError{},
		Warnings:	[]string{},
	}
	if err := ValidateVersion(c.Version); err!=nil{
		result.Errors   = append(result.Errors,*err)
		result.Valid    = false
	}

	if warnings := ValidateProjectType(c.Project.Type); len(warnings)>0{
		result.Warnings = append(result.Warnings,warnings...)
	}

	if errs,warns := ValidateRules(&c.Rules); len(errs)>0 || len(warns)>0{
		result.Errors   = append(result.Errors, errs...)
		result.Warnings = append(result.Warnings, warns...)
		if len(errs)>0{
			result.Valid = false
		}
	}
	return result
}
func ValidateVersion(version int) *ValidationError{
	if version<0{
		return &ValidationError{
			Field:		"version",
			Message:	"version must be >=0",
		}
	}
	return nil
}

func ValidateProjectType(pt string)[]string{
	warnings :=[]string{}

	if pt == ""{
		return warnings
	}
	supported :=[]string{
		"javascript",
	}
	found:=false
	found = slices.Contains(supported, pt)
	if !found{
		warnings=append(warnings,fmt.Sprintf(
			"project.type '%s' is not a known type (supported: %s). Will use generic rules.",
			pt,
			strings.Join(supported,", "),
		))
	}
	return warnings
}
func ValidateRules(rules *RulesType) ([]ValidationError, []string) {
	errors := []ValidationError{}
	warnings := []string{}

	if errs, warns := validateGeneralRules(&rules.General); len(errs) > 0 || len(warns) > 0 {
		errors   = append(errors, errs...)
		warnings = append(warnings, warns...)
	}
	if errs, warns := validateStructureRules(&rules.Structure); len(errs) > 0 || len(warns) > 0 {
		errors   = append(errors, errs...)
		warnings = append(warnings, warns...)
	}

	if errs, warns := validateDocumentationRules(&rules.Documentation); len(errs) > 0 || len(warns) > 0 {
		errors   = append(errors, errs...)
		warnings = append(warnings, warns...)
	}

	if errs, warns := validateCiCdRules(&rules.CiCd); len(errs) > 0 || len(warns) > 0 {
		errors   = append(errors, errs...)
		warnings = append(warnings, warns...)
	}

	if errs, warns := validateQualityRules(&rules.Quality); len(errs) > 0 || len(warns) > 0 {
		errors   = append(errors, errs...)
		warnings = append(warnings, warns...)
	}

	return errors, warnings
}
func validateGeneralRules(general *GeneralType) ([]ValidationError, []string) {
	errors := []ValidationError{}
	warnings := []string{}

	// README
	if errs, warns := validateRuleOptions("general.readme", &general.Readme); len(errs) > 0 || len(warns) > 0 {
		errors = append(errors, errs...)
		warnings = append(warnings, warns...)
	}

	// License
	if errs, warns := validateRuleOptions("general.license", &general.License); len(errs) > 0 || len(warns) > 0 {
		errors = append(errors, errs...)
		warnings = append(warnings, warns...)
	}

	// Gitignore
	if errs, warns := validateRuleOptions("general.gitignore", &general.Gitignore); len(errs) > 0 || len(warns) > 0 {
		errors = append(errors, errs...)
		warnings = append(warnings, warns...)
	}

	// Changelog
	if errs, warns := validateRuleOptions("general.changelog", &general.Changelog); len(errs) > 0 || len(warns) > 0 {
		errors = append(errors, errs...)
		warnings = append(warnings, warns...)
	}

	return errors, warnings
}

func validateStructureRules(structure *StructureType) ([]ValidationError, []string) {
	errors := []ValidationError{}
	warnings := []string{}

	// Src folder
	if errs, warns := validateRuleOptions("structure.src", &structure.Src); len(errs) > 0 || len(warns) > 0 {
		errors = append(errors, errs...)
		warnings = append(warnings, warns...)
	}

	// Tests folder
	if errs, warns := validateRuleOptions("structure.tests", &structure.Tests); len(errs) > 0 || len(warns) > 0 {
		errors = append(errors, errs...)
		warnings = append(warnings, warns...)
	}

	// Docs folder
	if errs, warns := validateRuleOptions("structure.docs", &structure.Docs); len(errs) > 0 || len(warns) > 0 {
		errors = append(errors, errs...)
		warnings = append(warnings, warns...)
	}
	return errors, warnings
}

func validateDocumentationRules(documentation *DocumentationType) ([]ValidationError, []string) {
	errors := []ValidationError{}
	warnings := []string{}

	// ADR
	if errs, warns := validateRuleOptions("documentation.adr", &documentation.ADR); len(errs) > 0 || len(warns) > 0 {
		errors = append(errors, errs...)
		warnings = append(warnings, warns...)
	}

	// Contributing
	if errs, warns := validateRuleOptions("documentation.contributing", &documentation.Contributing); len(errs) > 0 || len(warns) > 0 {
		errors = append(errors, errs...)
		warnings = append(warnings, warns...)
	}

	return errors, warnings
}

func validateCiCdRules(cicd *CiCdType) ([]ValidationError, []string) {
	errors := []ValidationError{}
	warnings := []string{}

	// GitHub Actions
	if errs, warns := validateRuleOptions("cicd.github_actions", &cicd.GithubActions); len(errs) > 0 || len(warns) > 0 {
		errors = append(errors, errs...)
		warnings = append(warnings, warns...)
	}

	// GitLab CI
	if errs, warns := validateRuleOptions("cicd.gitlab_ci", &cicd.GitlabCi); len(errs) > 0 || len(warns) > 0 {
		errors = append(errors, errs...)
		warnings = append(warnings, warns...)
	}

	return errors, warnings
}

func validateQualityRules(quality *QualituType) ([]ValidationError, []string) {
	errors := []ValidationError{}
	warnings := []string{}

	// Pre-commit
	if errs, warns := validateRuleOptions("quality.pre_commit", &quality.PreCommit); len(errs) > 0 || len(warns) > 0 {
		errors = append(errors, errs...)
		warnings = append(warnings, warns...)
	}

	// EditorConfig
	if errs, warns := validateRuleOptions("quality.editorconfig", &quality.EditorConfig); len(errs) > 0 || len(warns) > 0 {
		errors = append(errors, errs...)
		warnings = append(warnings, warns...)
	}

	return errors, warnings
}

func validateRuleOptions(rulePath string, rule *RulesOptions) ([]ValidationError, []string) {
	errors := []ValidationError{}
	warnings := []string{}
	if !rule.Enabled {
		return errors, warnings
	}

	validSeverities := []string{"error", "warning", "info"}
	if !slices.Contains(validSeverities, rule.Severity) {
		errors = append(errors, ValidationError{
			Field: rulePath + ".severity",
			Message: fmt.Sprintf(
				"invalid severity '%s' (valid: %s)",
				rule.Severity,
				strings.Join(validSeverities, ", "),
			),
		})
	}

	if len(rule.Patterns) == 0 {
		warnings = append(warnings, fmt.Sprintf(
			"%s.patterns is empty - rule may not work correctly",
			rulePath,
		))
	}

	if strings.TrimSpace(rule.Message) == "" {
		warnings = append(warnings, fmt.Sprintf(
			"%s.message is empty - users won't know what's wrong",
			rulePath,
		))
	}

	for i, pattern := range rule.Patterns {
		if strings.TrimSpace(pattern) == "" {
			warnings = append(warnings, fmt.Sprintf(
				"%s.patterns[%d] is empty",
				rulePath, i,
			))
		}
	}
	return errors, warnings
}

