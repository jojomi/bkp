package bkp

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	script "github.com/jojomi/go-script/v2"
	"github.com/rs/zerolog/log"
)

type ResticExecutor struct {
	DryRun bool

	context  *script.Context
	target   *Target
	cacheDir string

	hasNiceCommand   bool
	hasIONiceCommand bool
}

func NewResticExecutor() *ResticExecutor {
	ex := ResticExecutor{}
	ex.context = script.NewContext()

	ex.hasNiceCommand = ex.context.CommandExists("nice")
	ex.hasIONiceCommand = ex.context.CommandExists("ionice")

	return &ex
}

func (e *ResticExecutor) SetCacheDir(cacheDir string) {
	e.cacheDir = cacheDir
}

func (e *ResticExecutor) SetTarget(t *Target) {
	e.target = t
	if t.Password != "" {
		e.context.SetEnv("RESTIC_PASSWORD", t.Password)
	}
	e.context.SetEnv("RESTIC_REPOSITORY", t.Path)
}

func (e *ResticExecutor) Command(command string, args ...string) (*script.ProcessResult, error) {
	var fullArgs []string
	if e.context.CommandExists("nice") {
		var niceValue int
		switch runtime.GOOS {
		case "linux":
			niceValue = 19
		case "darwin":
			niceValue = 20
		default:
			niceValue = 19
		}
		niceArgs := []string{"nice", "-n", strconv.Itoa(niceValue)}
		log.Debug().Strs("args", niceArgs).Msg("nice is available, using it")
		fullArgs = append(fullArgs, niceArgs...)
	}
	if e.context.CommandExists("ionice") {
		ioniceArgs := []string{"ionice", "-c2", "-n7"}
		log.Debug().Strs("args", ioniceArgs).Msg("ionice is available, using it")
		fullArgs = append(fullArgs, ioniceArgs...)
	}
	resticBaseArgs := []string{"restic", command}
	log.Debug().Strs("args", resticBaseArgs).Msg("building basic restic command")
	fullArgs = append(fullArgs, resticBaseArgs...)

	if e.cacheDir != "" {
		cacheArgs := []string{"--cache-dir", e.cacheDir}
		log.Debug().Strs("args", cacheArgs).Msg("adding cache args")
		fullArgs = append(fullArgs, cacheArgs...)
	}

	log.Debug().Strs("args", args).Msg("adding restic command args")
	fullArgs = append(fullArgs, args...)

	localCommand := script.NewLocalCommand()
	localCommand.AddAll(fullArgs...)
	log.Info().
		Str("command", localCommand.String()).
		Strs("repository", e.context.GetCustomEnv()).
		Msg("full command")

	if e.DryRun {
		fmt.Println(strings.Join(fullArgs, " "))
		return nil, nil
	}
	return e.context.Execute(script.CommandConfig{
		RawStdout:    true,
		RawStderr:    true,
		ConnectStdin: false,
	}, localCommand)
}
