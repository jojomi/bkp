package main

import (
	"log"
	"os"

	"github.com/jojomi/bkp"
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
