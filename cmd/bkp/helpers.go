package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jojomi/go-script/v2"
	"github.com/mitchellh/go-homedir"
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

func restartAsRoot() int {
	lc := script.NewLocalCommand()
	lc.Add("sudo")
	lc.AddAll(os.Args...)

	context := script.NewContext()
	pr, err := context.ExecuteDebug(lc)
	if err != nil {
		log.Fatal().Err(err).Msg("execution failed")
	}
	exitCode, err := pr.ExitCode()
	if err != nil {
		log.Fatal().Err(err).Msg("could not determine exit code")
	}
	return exitCode
}
