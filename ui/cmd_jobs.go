package main

import (
	"fmt"

	"github.com/jojomi/bkp"
	"github.com/spf13/cobra"
)

func cmdJobs(cmd *cobra.Command, args []string) {
	sourceDirs := SourceDirs()
	jobs := bkp.AllJobs(sourceDirs)
	for _, job := range jobs {
		if flagJobsRelevant && !job.IsRelevant() {
			continue
		}
		if flagJobsRelevant && job.Target != nil && !job.Target.IsReady() {
			continue
		}
		fmt.Println(job)
	}

	fmt.Println()
	targets := bkp.AllTargets(sourceDirs)
	for _, target := range targets {
		if flagJobsRelevant && !target.IsReady() {
			continue
		}
		fmt.Println(target)
	}
}
