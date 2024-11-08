package main

import (
	"github.com/spf13/cobra"
)

type EnvGlobal struct {
	ConfigDirs []string

	Verbose bool
}

func (x *EnvGlobal) Parse(cmd *cobra.Command, _ []string) error {
	var (
		f   = cmd.Flags()
		err error
	)

	x.ConfigDirs, err = f.GetStringArray("config-dirs")
	if err != nil {
		return err
	}

	x.Verbose, err = f.GetBool("verbose")
	if err != nil {
		return err
	}

	return nil
}

func (x *EnvGlobal) HandleVerbosity() {
	handleVerbosityFlag(x.Verbose)
}

func (x *EnvGlobal) SourceDirs() []string {
	return SourceDirs(x.ConfigDirs)
}
