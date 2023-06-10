package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/jojomi/bkp"
	"github.com/jojomi/go-script/v2"
	"github.com/jojomi/go-script/v2/interview"
	"github.com/jojomi/go-script/v2/print"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func makeRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: buildName,
		Run: handleRoot,
	}

	// persistent flags
	pFlags := rootCmd.PersistentFlags()
	pFlags.BoolP("verbose", "v", false, "more verbose output")
	pFlags.BoolP("dry-run", "n", false, "dry run only")
	pFlags.StringArrayP("config-dirs", "c", []string{}, "override default config dirs")

	// local flags
	lFlags := rootCmd.Flags()
	lFlags.BoolP("all-jobs", "a", false, "execute all relevant backup jobs without asking")
	lFlags.StringArrayP("job", "j", []string{}, "execute backup jobs by name")

	lFlags.Bool("auto-unlock", false, "")
	lFlags.Bool("forget", false, "")
	lFlags.Bool("maintenance", false, "")
	lFlags.Bool("shutdown", false, "")

	return rootCmd
}

func handleVerbosityFlag(isVerbose bool) {
	if isVerbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		return
	}
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
}

func handleRoot(cmd *cobra.Command, args []string) {
	env, err := ParseEnvRoot(cmd, args)
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse CLI env")
	}
	err = cmdRoot(env)
	if err != nil {
		log.Fatal().Err(err).Msg("could not execute bkp")
	}
}

func cmdRoot(env EnvRoot) error {
	env.HandleVerbosity()
	print.Title("Check phase")

	err := check()
	if err != nil {
		log.Fatal().Err(err).Msg("System check failed.")
	}

	sourceDirs := env.SourceDirs()
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

	if len(env.Jobs) == 0 {
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
		jobNames := env.Jobs
		selections = make([]string, len(jobNames))
		for i, jobName := range jobNames {
			selections[i] = strings.TrimSpace(jobName)
		}

		selectedJobs = make([]*bkp.Job, 0)
		for _, selection := range selections {
			foundRelevant := false
			for _, relevantJob := range relevantJobs {
				if relevantJob.Name == selection {
					selectedJobs = append(selectedJobs, relevantJob)
					foundRelevant = true
					break
				}
			}
			if !foundRelevant {
				log.Warn().Msgf("Job %s not executed, because it was not found in relevant jobs", selection)
			}
		}
	}

	var (
		doUnlock      bool
		doForget      bool
		doMaintenance bool
		doShutdown    bool
	)

	if env.AutoUnlock != nil {
		doUnlock = *env.AutoUnlock
	} else {
		print.Subtitle("Auto unlock")
		doUnlock, err = interview.Confirm("Unlock repository if necessary?", true)
		if err != nil {
			log.Fatal().Err(err).Msg("No valid answer.")
		}
	}

	if env.Forget != nil {
		doForget = *env.Forget
	} else {
		print.Subtitle("Old Backups")
		doForget, err = interview.Confirm("Delete older backups as specified after finishing the new backup?", true)
		if err != nil {
			log.Fatal().Err(err).Msg("No valid answer.")
		}
	}

	if env.Maintenance != nil {
		doMaintenance = *env.Maintenance
	} else {
		print.Subtitle("Maintenance")
		doMaintenance, err = interview.Confirm("Do maintenance operations (takes a lot of time)?", false)
		if err != nil {
			log.Fatal().Err(err).Msg("No valid answer.")
		}
	}

	if env.Shutdown != nil {
		doShutdown = *env.Shutdown
	} else {
		print.Subtitle("Shutdown")
		doShutdown, err = interview.Confirm("Shutdown after finishing?", false)
		if err != nil {
			log.Fatal().Err(err).Msg("No valid answer.")
		}
	}

	fmt.Println()
	print.Title("Execution phase")
	for _, job := range selectedJobs {
		log.Debug().Str("job", job.String()).Msg("Executing job")
		err = job.Execute(bkp.JobExecuteOptions{
			DryRun:        env.DryRun,
			DoUnlock:      doUnlock,
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
		if !env.DryRun {
			timeSpec := "+5"
			log.Debug().Str("time specification", timeSpec).Msg("Shutdown scheduled")
			sc := script.NewContext()
			lc := script.NewLocalCommand()
			lc.AddAll("shutdown", "--poweroff", timeSpec)
			_, err := sc.ExecuteDebug(lc)
			if err != nil {
				return err
			}
		}
	}

	if !good {
		os.Exit(1)
	}

	return nil
}

func check() error {
	err := bkp.CheckEnvironment(minResticVersion)
	if err != nil {
		log.Fatal().Err(err).Msg("Environment check failed.")
	}

	// warn about nice (Linux, MacOS X) and ionice (Linux)
	sc := script.NewContext()
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		fmt.Println("Checking for nice command...")
		if !sc.CommandExists("nice") {
			log.Warn().Msg("\"nice\" command not found. Please make sure it is in your PATH to keep your system responsive while doing backups.")
		}
	}
	if runtime.GOOS == "linux" {
		fmt.Println("Checking for ionice command...")
		if !sc.CommandExists("ionice") {
			log.Warn().Msg("\"ionice\" command not found. Please make sure it is in your PATH to keep your system responsive while doing backups.")
		}
	}
	return err
}
