package shared

import (
	"fmt"
)

// Exit codes
const (
	ExitSuccess = 0
	ExitFailed  = 1
	ExitConfig  = 2
	ExitFS      = 3
	ExitArgs    = 4
)

// Error messages
const (
	ErrConfigNotFound   = "config_not_found"
	ErrInvalidConfig    = "invalid_config"
	ErrNoProject        = "no_project"
	ErrPermissionDenied = "permission_denied"
	ErrUnknownRule      = "unknown_rule"
	ErrFixFailed        = "fix_failed"
	ErrDetectionFailed  = "detection_failed"
)

func Version(version string) {
	fmt.Printf("PSX version %s\n", version)
}

func VerboseVersion(version string) {
	fmt.Printf("PSX version %s\nLicense: MIT\nRepository: https://github.com/m-mdy-m/psx\n", version)
}
