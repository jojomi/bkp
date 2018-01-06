package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func cmdJobs(cmd *cobra.Command, args []string) {
	jobs := AllJobs()
	for _, job := range jobs {
		if flagJobsRelevant && !job.IsRelevant() {
			continue
		}
		fmt.Println(job)
	}

	fmt.Println()
	targets := AllTargets()
	for _, target := range targets {
		if flagJobsRelevant && !target.IsReady() {
			continue
		}
		fmt.Println(target)
	}
}
