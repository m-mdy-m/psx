package config

import (
	"github.com/m-mdy-m/psx/internal/logger"
)

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

func ParseSeverity(val any, defaultSev Severity) (*Severity, error) {
	if b, ok := val.(bool); ok {
		if !b {
			return nil, nil
		}

		return nil, logger.Errorf("invalid value 'true' - use 'error','warning' or 'info'")
	}

	if s, ok := val.(string); ok {
		sev := Severity(s)
		if !sev.IsValid() {
			return nil, logger.Errorf("invalid severity '%s' - valid values: error,warning,info or false to disable", s)
		}
		return &sev, nil
	}

	if val == nil {
		return &defaultSev, nil
	}

	return nil, logger.Errorf("invalid type - must be sting ('error','warning','info') or false to disable")
}
func (s Severity) IsValid() bool {
	switch s {
	case SeverityError, SeverityWarning, SeverityInfo:
		return true
	default:
		return false
	}
}
