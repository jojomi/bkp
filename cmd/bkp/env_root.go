package main

import (
	"github.com/spf13/cobra"
)

type EnvRoot struct {
	EnvGlobal

	AllJobs    bool
	Jobs       []string
	ConfigDirs []string

	AutoUnlock  *bool
	Forget      *bool
	Maintenance *bool
	Shutdown    *bool

	DryRun  bool
	Verbose bool
}

func ParseEnvRoot(cmd *cobra.Command, args []string) (*EnvRoot, error) {
	var (
		f   = cmd.Flags()
		err error
	)
	env := &EnvRoot{}
	err = env.EnvGlobal.Parse(cmd, args)
	if err != nil {
		return env, err
	}

	env.AllJobs, err = f.GetBool("all-jobs")
	if err != nil {
		return env, err
	}

	env.Jobs, err = f.GetStringArray("job")
	if err != nil {
		return env, err
	}

	env.ConfigDirs, err = f.GetStringArray("config-dirs")
	if err != nil {
		return env, err
	}

	if f.Changed("auto-unlock") {
		v, err := f.GetBool("auto-unlock")
		if err != nil {
			return env, err
		}
		env.AutoUnlock = &v
	}

	if f.Changed("forget") {
		v, err := f.GetBool("forget")
		if err != nil {
			return env, err
		}
		env.Forget = &v
	}

	if f.Changed("maintenance") {
		v, err := f.GetBool("maintenance")
		if err != nil {
			return env, err
		}
		env.Maintenance = &v
	}

	if f.Changed("shutdown") {
		v, err := f.GetBool("shutdown")
		if err != nil {
			return env, err
		}
		env.Shutdown = &v
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
