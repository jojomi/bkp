package main

import (
	"os"
	"time"

	"github.com/blang/semver"
	"github.com/jojomi/go-script/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	minResticVersion semver.Version
)

func main() {
	context := script.NewContext()
	if !context.IsUserRoot() {
		os.Exit(restartAsRoot())
	}

	log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	// start handling commandline
	rootCmd := makeRootCmd()
	rootCmd.AddCommand(getEnvCmd(), getJobsCmd(), getSnapshotsCmd(), getMountCmd())
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal().Err(err).Msg("execution failed")
	}
}

func init() {
	minResticVersion, _ = semver.Make("0.15.0")

	// setup logging
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
