package bkp

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	script "github.com/jojomi/go-script"
)

type ResticExecutor struct {
	DryRun bool

	context *script.Context
	target  *Target

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
		fullArgs = []string{"nice", "-n", strconv.Itoa(niceValue)}
	}
	if e.context.CommandExists("ionice") {
		fullArgs = mergeStringSlices(fullArgs, []string{"ionice", "-c2", "-n7"})
	}
	fullArgs = mergeStringSlices(fullArgs, []string{"restic", command})

	for _, a := range args {
		fullArgs = append(fullArgs, a)
	}

	if e.DryRun {
		fmt.Println(strings.Join(fullArgs, " "))
		return nil, nil
	}
	return e.context.ExecuteDebug(fullArgs[0], fullArgs[1:]...)
}
