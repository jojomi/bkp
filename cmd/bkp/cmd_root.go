package main

import (
	"fmt"
	"os"

	"github.com/jojomi/bkp"
	"github.com/spf13/cobra"
)

func makeRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: buildName,
		Run: cmdRoot,
	}
	rootCmd.PersistentFlags().BoolVarP(&flagDryRun, "dry-run", "d", false, "dry run only")
	return rootCmd
}

func cmdRoot(cmd *cobra.Command, args []string) {
	sourceDirs := SourceDirs()
	jl := bkp.JobList{}
	jl.Load(sourceDirs)

	var (
		err  error
		good = true
	)

	for _, job := range jl.Relevant() {
		err = job.Execute(bkp.JobExecuteOptions{
			DryRun: flagDryRun,
		})
		if err != nil {
			fmt.Println("Backup error", err)
			good = false
		}
		fmt.Println()
	}

	if !good {
		os.Exit(1)
	}
}
