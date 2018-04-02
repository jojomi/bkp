package main

import (
	"os"

	"github.com/jojomi/bkp"
	"github.com/jojomi/go-script"
	"github.com/jojomi/go-script/print"
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

	sc := script.NewContext()

	re := bkp.NewResticExecutor()
	re.SetTarget(target)
	err := sc.EnsureDirExists(target.RestoreDir, os.FileMode(int(0750)))
	if err != nil {
		sugar.Fatal(err)
	}
	print.Boldf("Mounting at %s\n", target.RestoreDir)
	sc.ExecuteSilent("xdg-open", target.RestoreDir)
	re.Command("mount", target.RestoreDir)
}
