package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/jojomi/bkp"
	script "github.com/jojomi/go-script"
	"github.com/jojomi/go-script/interview"
	"github.com/jojomi/go-script/print"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	flagRootVerbose    bool
	flagRootDryRun     bool
	flagRootAllJobs    bool
	flagRootJobs       string
	flagRootConfigDirs string
)

func makeRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: buildName,
		Run: cmdRoot,
	}
	rootCmd.PersistentFlags().BoolVarP(&flagRootVerbose, "verbose", "v", false, "verbose output (useful for debugging)")
	rootCmd.PersistentFlags().BoolVarP(&flagRootDryRun, "dry-run", "d", false, "dry run only")
	rootCmd.Flags().BoolVarP(&flagRootAllJobs, "all-jobs", "a", false, "execute all relevant backup jobs without asking")
	rootCmd.Flags().StringVarP(&flagRootJobs, "jobs", "j", "", "execute a backup jobs by name")
	rootCmd.Flags().StringVarP(&flagRootConfigDirs, "config-dirs", "c", "", "override default config dirs")
	return rootCmd
}

func handleVerbosityFlag(isVerbose bool) {
	if isVerbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		return
	}
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
}

func cmdRoot(cmd *cobra.Command, args []string) {
	handleVerbosityFlag(flagRootVerbose)
	print.Title("Check phase")

	err := check()
	if err != nil {
		log.Fatal().Err(err).Msg("System check failed.")
	}

	sourceDirs := SourceDirs()
	jl := bkp.JobList{}
	jl.Load(sourceDirs)

	var (
		good = true
	)

	relevantJobs := jl.Relevant()

	if len(relevantJobs) == 0 {
		log.Fatal().Msg("No relevant jobs found. Did you connect the backup targets?")
	}

	var (
		selections   []string
		selectedJobs []*bkp.Job
	)

	fmt.Println()
	print.Title("Preparation phase")

	if flagRootJobs == "" {
		selectionMap := make(map[string]*bkp.Job, len(relevantJobs))
		options := make([]string, len(relevantJobs))
		for i, job := range relevantJobs {
			jobIdentifier := job.String()
			// TODO if already set there is jobs with equal name, generate an error message and abort run
			selectionMap[jobIdentifier] = job
			options[i] = jobIdentifier
		}

		print.Subtitle("Backup selection")
		selections, err = interview.ChooseMultiStrings("Which backups should be executed? (Spacebar to select, Return to start backup)", options)
		if err != nil {
			log.Fatal().Err(err).Msg("No valid job selection")
		}
		selectedJobs = make([]*bkp.Job, len(selections))
		for i, selection := range selections {
			selectedJobs[i] = selectionMap[selection]
		}
	} else {
		// job list given on CLI
		jobNames := strings.Split(flagRootJobs, ",")
		selections = make([]string, len(jobNames))
		for i, jobName := range jobNames {
			selections[i] = strings.TrimSpace(jobName)
		}

		selectedJobs = make([]*bkp.Job, 0)
	outside:
		for _, selection := range selections {
			for _, relevantJob := range relevantJobs {
				if relevantJob.Name == selection {
					selectedJobs = append(selectedJobs, relevantJob)
					break outside
				}
			}
		}
	}

	print.Subtitle("Old Backups")
	doForget, err := interview.Confirm("Delete older backups as specified after finishing the new backup?", true)
	if err != nil {
		log.Fatal().Err(err).Msg("No valid answer.")
	}

	print.Subtitle("Maintenance")
	doMaintenance, err := interview.Confirm("Do maintenance operations (takes a lot of time)?", false)
	if err != nil {
		log.Fatal().Err(err).Msg("No valid answer.")
	}

	print.Subtitle("Shutdown")
	doShutdown, err := interview.Confirm("Shutdown after finishing?", false)
	if err != nil {
		log.Fatal().Err(err).Msg("No valid answer.")
	}

	fmt.Println()
	print.Title("Execution phase")
	for _, job := range selectedJobs {
		err = job.Execute(bkp.JobExecuteOptions{
			DryRun:        flagRootDryRun,
			DoForget:      doForget,
			DoMaintenance: doMaintenance,
		})
		if err != nil {
			fmt.Println("Backup error", err)
			good = false
		}
		fmt.Println()
	}

	if doShutdown {
		print.Title("Shutdown")
		if !flagRootDryRun {
			sc := script.NewContext()
			lc := script.NewLocalCommand()
			lc.AddAll("shutdown", "--poweroff", "+5")
			sc.ExecuteDebug(lc)
		}
	}

	if !good {
		os.Exit(1)
	}
}

func check() error {
	err := bkp.CheckEnvironment(minResticVersion)
	if err != nil {
		log.Fatal().Err(err).Msg("Environment check failed.")
	}

	// warn about nice (Linux, MacOS X) and ionice (Linux)
	sc := script.NewContext()
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		if !sc.CommandExists("nice") {
			log.Warn().Msg("\"nice\" command not found. Please make sure it is in your PATH to keep your system responsive while doing backups.")
		}
	}
	if runtime.GOOS == "linux" {
		if !sc.CommandExists("ionice") {
			log.Warn().Msg("\"ionice\" command not found. Please make sure it is in your PATH to keep your system responsive while doing backups.")
		}
	}
	return err
}
