package main

import (
	"fmt"
	"os"
	"path/filepath"

	script "github.com/jojomi/go-script"
	homedir "github.com/mitchellh/go-homedir"
)

func SourceDirs() []string {
	homePath, _ := homedir.Expand(fmt.Sprintf("~/.%s", buildName))
	workPath, _ := filepath.Abs(fmt.Sprintf(".%s", buildName))
	return []string{
		fmt.Sprintf("/etc/%s", buildName),
		homePath,
		workPath,
	}
}

func forceRoot() bool {
	context := script.NewContext()
	if context.IsUserRoot() {
		return false
	}
	context.ExecuteDebug("sudo", os.Args...)
	return true
}
