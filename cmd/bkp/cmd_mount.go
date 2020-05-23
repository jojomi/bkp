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
		Run:   cmdMount,
	}
	return mountCmd
}

func cmdMount(cmd *cobra.Command, args []string) {
	handleVerbosityFlag(flagRootVerbose)
	if len(args) < 1 {
		// TODO add selection dialog instead
		log.Fatal().Msg("No target given")
	}

	targetName := args[0]
	target := bkp.TargetByName(targetName, SourceDirs())

	sc := script.NewContext()

	re := bkp.NewResticExecutor()
	re.SetTarget(target)
	err := sc.EnsureDirExists(target.RestoreDir, os.FileMode(int(0750)))
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	print.Boldf("Mounting at %s\n", target.RestoreDir)
	lc := script.NewLocalCommand()
	lc.AddAll("xdg-open", target.RestoreDir)
	sc.ExecuteSilent(lc)
	re.Command("mount", target.RestoreDir)
}
