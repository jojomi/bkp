package main

import (
	"github.com/jojomi/bkp"
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
	if len(args) < 1 {
		// TODO add selection dialog instead
		sugar.Fatal("No target given")
	}

	targetName := args[0]
	target := bkp.TargetByName(targetName, SourceDirs())
	re := bkp.NewResticExecutor()
	re.SetTarget(target)
	/// re.DryRun = true
	// TODO create restore dir if it does not exist
	re.Command("mount", target.RestoreDir)
}
