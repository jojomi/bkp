package main

import (
	"github.com/jojomi/bkp"
	"github.com/spf13/cobra"
)

func getSnapshotsCmd() *cobra.Command {
	snapshotsCmd := &cobra.Command{
		Use:   "snapshots [target name] [extra restic args]",
		Short: "Lists snapshots for a given target",
		Run:   cmdSnapshots,
	}
	return snapshotsCmd
}

func cmdSnapshots(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		// TODO add selection dialog instead
		sugar.Fatal("No target given")
	}

	targetName := args[0]
	target := bkp.TargetByName(targetName, SourceDirs())

	re := bkp.NewResticExecutor()
	re.SetTarget(target)
	re.Command("snapshots", args[1:]...)
}
