package main

import (
	script "github.com/jojomi/go-script"
)

func executeResticCommand(sc *script.Context, command string, args ...string) (*script.ProcessResult, error) {
	fullArgs := []string{command}
	for _, a := range args {
		fullArgs = append(fullArgs, a)
	}
	// TODO check if nice and ionice are available (Windows...)
	return sc.ExecuteDebug("nice", "-n", "19", "ionice", "-c2", "-n7", "restic", fullArgs...)
}
