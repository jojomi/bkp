package main

import (
	"github.com/jojomi/bkp"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func getSnapshotsCmd() *cobra.Command {
	snapshotsCmd := &cobra.Command{
		Use:   "snapshots [target name] [extra restic args]",
		Short: "Lists snapshots for a given target",
		Run:   handleSnapshots,
	}
	return snapshotsCmd
}

func handleSnapshots(cmd *cobra.Command, args []string) {
	env, err := ParseEnvSnapshots(cmd, args)
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse CLI env")
	}
	err = cmdSnapshots(env)
	if err != nil {
		log.Fatal().Err(err).Msg("could not execute bkp")
	}
}

func cmdSnapshots(env EnvSnapshots) error {
	env.HandleVerbosity()
	if len(env.Targets) < 1 {
		// TODO add selection dialog instead
		log.Fatal().Msg("No target given")
	}

	targetName := env.Targets[0]
	target := bkp.TargetByName(targetName, env.SourceDirs())

	re := bkp.NewResticExecutor()
	re.SetTarget(target)
	re.Command("snapshots", env.Args...)

	return nil
}
