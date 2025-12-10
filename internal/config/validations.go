package config

import (
	"fmt"
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

	if errs,warns := ValidateRules(c.Rules); len(errs)>0 || len(warns)>0{
		result.Errors   = append(result.Errors, errs...)
		result.Warnings = append(result.Warnings, warns...)
		if len(errs)>0{
			result.Valid = false
		}
	}

	if warns := validateIgnorePatterns(c.Ignore); len(warns) > 0 {
		result.Warnings = append(result.Warnings, warns...)
	}

	if warns := validateFixConfig(&c.Fix); len(warns) > 0 {
		result.Warnings = append(result.Warnings, warns...)
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
	supportedLang :=[]string{
		"javascript",
		"go",
	}
	found:=false
	for _,supported := range supportedLang{
		if pt == supported{
			found = true
			break
		}
	}
	if !found{
		warnings=append(warnings,fmt.Sprintf(
			"project.type '%s' is not a known type (supported: %s). Will use generic rules.",
			pt,
			strings.Join(supportedLang,", "),
		))
	}
	return warnings
}
func ValidateRules(rules map[string]RulesSeverity) ([]ValidationError, []string) {
	errors   := []ValidationError{}
	warnings := []string{}
	if rules == nil || len(rules)==0{
		warnings = append(warnings,"No rules configured")
		return errors, warnings
	}

	metadata := GetRulesMetadata()

	for id,severtity := range rules{
		ruleMeta,exists:= metadata.Rules[id]
		if !exists{
			warnings = append(warnings,fmt.Sprintf(
				"Unknown rule '%s' - will be ignored",id))
			continue
		}
		if err := validateRuleSeverity(id,severtity,ruleMeta.DefaultSeverity); err !=nil{
			errors = append(errors,*err)
		}
	}

	criticalRules := []string{"readme","license"}
	for _,id := range criticalRules{
		if sev,exists := rules[id]; exists{
			if b,ok := sev.(bool);ok && !b{
				warnings = append(warnings,fmt.Sprintf("Critical rule '%s' is disabled - this is not recommended",	id))
			}
		}
	}
	return errors,warnings
}

func validateRuleSeverity(id string, severity RulesSeverity,defaultSev Severity) *ValidationError {
	_,err :=ParseSeverity(severity,defaultSev)
	if err !=nil{
		return &ValidationError{
			Field:		fmt.Sprintf("rules.%s",id),
			Message:	err.Error(),
		}
	}
	return nil
}

func validateIgnorePatterns(patterns []string) []string {
	warnings := []string{}

	if len(patterns) == 0 {
		// No patterns is OK, just informational
		return warnings
	}

	for i, pattern := range patterns {
		if strings.TrimSpace(pattern) == "" {
			warnings = append(warnings, fmt.Sprintf(
				"ignore[%d] is empty",
				i,
			))
			continue
		}

		if pattern == "*" || pattern == "**" {
			warnings = append(warnings, fmt.Sprintf(
				"ignore[%d]: pattern '%s' will ignore everything - is this intentional?",
				i, pattern,
			))
		}

		// Warn about absolute paths
		if strings.HasPrefix(pattern, "/") {
			warnings = append(warnings, fmt.Sprintf(
				"ignore[%d]: absolute path '%s' - relative paths are recommended",
				i, pattern,
			))
		}
	}

	return warnings
}


func validateFixConfig(fix *FixConfig) []string {
	warnings := []string{}
	// Note: interactive and create_backups are just booleans, always valid
	return warnings
}

