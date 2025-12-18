package rules

import (
	"github.com/m-mdy-m/psx/internal/utils"
)

func ValidateExists(ctx *Context, fullPath string, info any) (bool, string, error) {
	return true, "", nil
}

func ValidateNotEmpty(ctx *Context, fullPath string, info any) (bool, string, error) {
	fileInfo, ok := info.(interface {
		IsDir() bool
		Size() int64
	})
	if !ok {
		return false, "Cannot validate path", nil
	}

	if fileInfo.IsDir() {
		isEmpty, err := utils.IsDirEmpty(fullPath)
		if err != nil {
			return false, "Cannot check directory contents", err
		}
		if isEmpty {
			return false, "Folder exists but is empty", nil
		}
	} else {
		if fileInfo.Size() < 10 {
			return false, "File exists but appears to be empty", nil
		}
	}

	return true, "", nil
}
