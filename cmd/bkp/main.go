package main

import (
	"log"
	"os"

	"github.com/jojomi/bkp"
	script "github.com/jojomi/go-script"
	"github.com/spf13/cobra"
)

var (
	context         = script.NewContext()
	backupTargetDir = ""

	flagDryRun       bool
	flagJobsRelevant bool
)

func main() {
	if forceRoot() {
		os.Exit(0)
	}

	err := bkp.CheckEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	rootCmd := &cobra.Command{
		Use: buildName,
		Run: cmdRoot,
	}
	rootCmd.PersistentFlags().BoolVarP(&flagDryRun, "dry-run", "d", false, "dry run only")

	jobsCmd := &cobra.Command{
		Use:   "jobs",
		Short: "Lists all backup jobs available",
		Run:   cmdJobs,
	}
	jobsCmd.PersistentFlags().BoolVarP(&flagJobsRelevant, "relevant", "r", false, "only show relevant jobs for the current environment")

	rootCmd.AddCommand(jobsCmd)
	rootCmd.Execute()
}
