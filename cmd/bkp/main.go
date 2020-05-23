package main

import (
	"os"
	"time"

	"github.com/blang/semver"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	minResticVersion semver.Version
)

func main() {
	if forceRoot() {
		os.Exit(0)
	}

	log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	// start handling commandline
	rootCmd := makeRootCmd()
	rootCmd.AddCommand(getEnvCmd(), getJobsCmd(), getSnapshotsCmd(), getMountCmd())
	rootCmd.Execute()
}

func init() {
	minResticVersion, _ = semver.Make("0.9.5")
}
