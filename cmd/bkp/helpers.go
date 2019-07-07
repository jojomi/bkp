package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	script "github.com/jojomi/go-script"
	homedir "github.com/mitchellh/go-homedir"
)

func SourceDirs() []string {
	if flagRootConfigDirs != "" {
		paths := strings.Split(flagRootConfigDirs, ",")
		result := make([]string, len(paths))
		for i, path := range paths {
			result[i] = strings.TrimSpace(path)
		}
		return result
	}

	// default values
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
	lc := script.NewLocalCommand()
	lc.Add("sudo")
	lc.AddAll(os.Args...)
	context.ExecuteDebug(lc)
	return true
}
