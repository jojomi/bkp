package main

import (
	"fmt"
	"github.com/jojomi/bkp"
	"github.com/jojomi/go-script/v2/print"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func getJobsCmd() *cobra.Command {
	jobsCmd := &cobra.Command{
		Use:   "jobs",
		Short: "Lists all backup jobs available",
		Run:   handleJobs,
	}

	jobsCmd.PersistentFlags().BoolP("relevant", "r", false, "only show relevant jobs for the current environment")

	return jobsCmd
}

func handleJobs(cmd *cobra.Command, args []string) {
	env, err := ParseEnvJobs(cmd, args)
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse CLI env")
	}
	err = cmdJobs(env)
	if err != nil {
		log.Fatal().Err(err).Msg("could not execute bkp")
	}
}

func cmdJobs(env EnvJobs) error {
	env.HandleVerbosity()
	sourceDirs := env.SourceDirs()
	jl := bkp.JobList{}
	jl.Load(sourceDirs)

	print.Boldf("%d jobs evaluated\n", len(jl.All()))
	if env.RelevantOnly {
		fmt.Printf("%d jobs relevant\n", len(jl.Relevant()))
	}

	var jobs []*bkp.Job
	if env.RelevantOnly {
		jobs = jl.Relevant()
	} else {
		jobs = jl.All()
	}

	for _, job := range jobs {
		fmt.Println(job)
	}

	return nil
}
