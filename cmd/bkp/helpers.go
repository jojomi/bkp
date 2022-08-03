package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	script "github.com/jojomi/go-script/v2"
	homedir "github.com/mitchellh/go-homedir"
)

func SourceDirs(input []string) []string {
	if len(input) > 0 {
		var result []string
		for i, path := range input {
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
