package main

import (
	"fmt"

	"github.com/jojomi/bkp"
	script "github.com/jojomi/go-script"
	"github.com/spf13/cobra"
)

func cmdJobs(cmd *cobra.Command, args []string) {
	sourceDirs := SourceDirs()
	jl := bkp.JobList{}
	jl.Load(sourceDirs)

	c := script.NewContext()
	c.PrintfBold("%d jobs evaluated\n", len(jl.All()))
	if flagJobsRelevant {
		fmt.Printf("%d jobs relevant\n", len(jl.Relevant()))
	}

	var jobs []*bkp.Job
	if flagJobsRelevant {
		jobs = jl.Relevant()
	} else {
		jobs = jl.All()
	}

	for _, job := range jobs {
		fmt.Println(job)
	}
}
