package main

import (
	"fmt"
	"os"

	"github.com/jojomi/bkp"
	"github.com/spf13/cobra"
)

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
