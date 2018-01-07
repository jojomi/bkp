package main

import (
	"fmt"

	script "github.com/jojomi/go-script"
)

func executeResticCommand(sc *script.Context, command string, args ...string) (*script.ProcessResult, error) {
	var fullArgs []string
	if sc.CommandExists("nice") && sc.CommandExists("ionice") {
		fullArgs = []string{"nice", "-n", "19", "ionice", "-c2", "-n7", "restic", command}
	} else {
		fullArgs = []string{"restic", command}
	}

	for _, a := range args {
		fullArgs = append(fullArgs, a)
	}

	fmt.Println("Executing", fullArgs)
	return nil, nil
	// return sc.ExecuteDebug(fullArgs[0], fullArgs[1:]...)
}
