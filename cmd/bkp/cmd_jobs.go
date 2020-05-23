package main

import (
	"fmt"

	"github.com/jojomi/bkp"
	"github.com/jojomi/go-script/print"
	"github.com/spf13/cobra"
)

var (
	flagJobsRelevant bool
)

func getJobsCmd() *cobra.Command {
	jobsCmd := &cobra.Command{
		Use:   "jobs",
		Short: "Lists all backup jobs available",
		Run:   cmdJobs,
	}
	jobsCmd.PersistentFlags().BoolVarP(&flagJobsRelevant, "relevant", "r", false, "only show relevant jobs for the current environment")
	return jobsCmd
}

func cmdJobs(cmd *cobra.Command, args []string) {
	handleVerbosityFlag(flagRootVerbose)
	sourceDirs := SourceDirs()
	jl := bkp.JobList{}
	jl.Load(sourceDirs)

	print.Boldf("%d jobs evaluated\n", len(jl.All()))
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
