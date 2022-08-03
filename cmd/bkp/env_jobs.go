package main

import (
	"github.com/spf13/cobra"
)

type EnvJobs struct {
	EnvGlobal

	RelevantOnly bool

	Verbose bool
}

func ParseEnvJobs(cmd *cobra.Command, args []string) (EnvJobs, error) {
	var (
		f   = cmd.Flags()
		err error
	)
	env := EnvJobs{}

	env.Parse(cmd, args)

	env.RelevantOnly, err = f.GetBool("relevant")
	if err != nil {
		return env, err
	}

	return env, nil
}
