package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/jojomi/bkp"
	script "github.com/jojomi/go-script"
	"github.com/spf13/cobra"
)

var (
	flagRootDryRun bool
	flagRootAll    bool
)

func makeRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: buildName,
		Run: cmdRoot,
	}
	rootCmd.Flags().BoolVarP(&flagRootDryRun, "dry-run", "d", false, "dry run only")
	rootCmd.Flags().BoolVarP(&flagRootAll, "all", "a", false, "execute all relevant backup jobs without asking")
	return rootCmd
}

func cmdRoot(cmd *cobra.Command, args []string) {
	err := bkp.CheckEnvironment(minResticVersion)
	if err != nil {
		sugar.Fatal(err)
	}

	// warn about nice (Linux, MacOS X) and ionice (Linux)
	sc := script.NewContext()
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		if !sc.CommandExists("nice") {
			sugar.Warn("\"nice\" command not found. Please make sure it is in your PATH to keep your system responsive while doing backups.")
		}
	}
	if runtime.GOOS == "linux" {
		if !sc.CommandExists("ionice") {
			sugar.Warn("\"ionice\" command not found. Please make sure it is in your PATH to keep your system responsive while doing backups.")
		}
	}

	sourceDirs := SourceDirs()
	jl := bkp.JobList{}
	jl.Load(sourceDirs)

	var (
		good = true
	)

	ctx := script.NewContext()
	relevantJobs := jl.Relevant()
	selectionMap := make(map[string]*bkp.Job, len(relevantJobs))
	options := make([]string, len(relevantJobs))
	for i, job := range relevantJobs {
		jobIdentifier := job.String()
		// TODO if already set there is jobs with equal name, generate an error message and abort run
		selectionMap[jobIdentifier] = job
		options[i] = jobIdentifier
	}

	selections, err := ctx.ChooseMultiStrings("Which backups should be executed? (Spacebar to select, Return to start backup)", options)
	if err != nil {
		log.Fatal(err)
	}
	selectedJobs := make([]*bkp.Job, len(selections))
	for i, selection := range selections {
		selectedJobs[i] = selectionMap[selection]
	}

	for _, job := range selectedJobs {
		err = job.Execute(bkp.JobExecuteOptions{
			DryRun: flagRootDryRun,
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
