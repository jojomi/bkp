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

	err := env.Parse(cmd, args)
	if err != nil {
		return EnvMount{}, err
	}

	env.Targets = args

	return env, nil
}
