package main

import (
	"github.com/spf13/cobra"
)

type EnvRoot struct {
	EnvGlobal

	AllJobs    bool
	Jobs       []string
	ConfigDirs []string

	DryRun  bool
	Verbose bool
}

func ParseEnvRoot(cmd *cobra.Command, args []string) (EnvRoot, error) {
	var (
		f   = cmd.Flags()
		err error
	)
	env := EnvRoot{}
	env.EnvGlobal.Parse(cmd, args)

	env.AllJobs, err = f.GetBool("all-jobs")
	if err != nil {
		return env, err
	}

	env.Jobs, err = f.GetStringArray("jobs")
	if err != nil {
		return env, err
	}

	env.ConfigDirs, err = f.GetStringArray("config-dirs")
	if err != nil {
		return env, err
	}

	env.DryRun, err = f.GetBool("dry-run")
	if err != nil {
		return env, err
	}

	env.Verbose, err = f.GetBool("verbose")
	if err != nil {
		return env, err
	}

	return env, nil
}
