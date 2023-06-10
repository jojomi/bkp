package main

import (
	"os"

	"github.com/jojomi/bkp"
	"github.com/jojomi/go-script/v2"
	"github.com/jojomi/go-script/v2/print"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func getMountCmd() *cobra.Command {
	mountCmd := &cobra.Command{
		Use:   "mount [target name]",
		Short: "Mounts a target for restore",
		Run:   handleMount,
	}
	return mountCmd
}

func handleMount(cmd *cobra.Command, args []string) {
	env, err := ParseEnvMount(cmd, args)
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse CLI env")
	}
	err = cmdMount(env)
	if err != nil {
		log.Fatal().Err(err).Msg("could not execute bkp")
	}
}

func cmdMount(env EnvMount) error {
	env.HandleVerbosity()
	if len(env.Targets) < 1 {
		// TODO add selection dialog instead
		log.Fatal().Msg("No target given")
	}

	targetName := env.Targets[0]
	target := bkp.TargetByName(targetName, env.SourceDirs())

	sc := script.NewContext()

	re := bkp.NewResticExecutor()
	re.SetTarget(target)
	err := sc.EnsureDirExists(target.RestoreDir, os.FileMode(0750))
	if err != nil {
		return err
	}
	print.Boldf("Mounting at %s\n", target.RestoreDir)
	lc := script.NewLocalCommand()
	lc.AddAll("xdg-open", target.RestoreDir)
	_, err = sc.ExecuteSilent(lc)
	if err != nil {
		return err
	}
	_, err = re.Command("mount", target.RestoreDir)
	if err != nil {
		return err
	}

	return nil
}
