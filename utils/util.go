package utils

import (
  "os"
  "path/filepath"
)

func IsGoonRepo(currDir, targetDir string) (string,bool) {
	// check if .goon already exists in any path of the existing director

	if currDir == "/" {
		return "",false
	}

	targetPath := filepath.Join(currDir, targetDir)

	//  check if this directory exists
	_, err := os.Stat(targetPath)
	if err == nil {

		return targetPath,true
	}

	parentDir := filepath.Dir(currDir)
	return IsGoonRepo(parentDir, targetDir)
}
