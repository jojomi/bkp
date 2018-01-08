package main

import (
	"fmt"
	"os"

	"github.com/jojomi/bkp"
	"github.com/spf13/cobra"
)

func cmdRoot(cmd *cobra.Command, args []string) {
	jobs := bkp.AllJobs(SourceDirs())

	var (
		err  error
		good = true
	)

	for _, job := range jobs {
		if !job.IsRelevant() {
			continue
		}
		if job.Target == nil || !job.Target.IsReady() {
			continue
		}

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
