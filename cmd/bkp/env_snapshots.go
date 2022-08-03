package main

import (
	"github.com/spf13/cobra"
)

type EnvSnapshots struct {
	EnvGlobal

	Targets []string
	Args    []string

	Verbose bool
}

func ParseEnvSnapshots(cmd *cobra.Command, args []string) (EnvSnapshots, error) {
	env := EnvSnapshots{}

	env.Parse(cmd, args)

	if len(args) > 0 {
		env.Targets = args[0:1]
	}

	if len(args) > 1 {
		env.Args = args[1:]
	}

	return env, nil
}
