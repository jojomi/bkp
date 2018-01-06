package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func cmdRoot(cmd *cobra.Command, args []string) {
	jobs := AllJobs()

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

		fmt.Println("Executing Job", job)
		err = job.Execute()
		if err != nil {
			fmt.Println("Backup error", err)
			good = false
		}
	}

	if !good {
		os.Exit(1)
	}
}
