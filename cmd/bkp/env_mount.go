package main

import (
	"github.com/spf13/cobra"
)

type EnvMount struct {
	EnvGlobal

	Targets []string

	Verbose bool
}

func ParseEnvMount(cmd *cobra.Command, args []string) (EnvMount, error) {
	env := EnvMount{}

	env.Parse(cmd, args)

	env.Targets = args

	return env, nil
}
