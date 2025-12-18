package rules

import "os"

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
		entries, err := os.ReadDir(fullPath)
		if err != nil {
			return false, "Cannot read directory", err
		}
		if len(entries) == 0 {
			return false, "Folder is empty", nil
		}
	} else {
		if fileInfo.Size() < 10 {
			return false, "File is too small or empty", nil
		}
	}

	return true, "", nil
}

func ValidateFileSize(minSize int64) ValidatorFunc {
	return func(ctx *Context, fullPath string, info any) (bool, string, error) {
		fileInfo, ok := info.(interface {
			IsDir() bool
			Size() int64
		})
		if !ok {
			return false, "Cannot validate path", nil
		}

		if fileInfo.IsDir() {
			return true, "", nil
		}

		if fileInfo.Size() < minSize {
			return false, "File is too small", nil
		}

		return true, "", nil
	}
}
