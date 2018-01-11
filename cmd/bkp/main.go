package main

import (
	"log"
	"os"

	"github.com/jojomi/bkp"
	script "github.com/jojomi/go-script"
)

var (
	context         = script.NewContext()
	backupTargetDir = ""

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

	rootCmd := makeRootCmd()
	rootCmd.AddCommand(getJobsCmd())
	rootCmd.Execute()
}
