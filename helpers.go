package main

import (
	"fmt"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

func SourceDirs() []string {
	homePath, _ := homedir.Expand(fmt.Sprintf("~/.%s", buildName))
	workPath, _ := filepath.Abs(fmt.Sprintf(".%s", buildName))
	return []string{
		fmt.Sprintf("/etc/.%s", buildName),
		homePath,
		workPath,
	}
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func forceRoot() bool {
	if context.IsUserRoot() {
		return false
	}
	context.ExecuteDebug("sudo", os.Args...)
	return true
}
