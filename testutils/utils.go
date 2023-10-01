package testutils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// GetRootFolder determines the root folder of the project
func GetRootFolder() (string, error) {
	wd, _ := os.Getwd()
	return lookupRootFolder(wd, 0)
}

func lookupRootFolder(dir string, level int) (string, error) {
	if level > 5 {
		return "", errors.New("Failed to find root folder")
	}
	mainFile := fmt.Sprintf("%s/main.go", dir)
	if fileExists(mainFile) {
		return dir, nil
	}
	nextLevel := level + 1
	parentDir := filepath.Dir(dir)
	return lookupRootFolder(parentDir, nextLevel)
}

func fileExists(file string) bool {
	if stat, err := os.Stat(file); err == nil {
		return !stat.IsDir()
	}
	return false
}
