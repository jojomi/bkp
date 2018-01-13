package main

import (
	"fmt"

	"github.com/jojomi/bkp"
	script "github.com/jojomi/go-script"
	"github.com/spf13/cobra"
)

func getEnvCmd() *cobra.Command {
	envCmd := &cobra.Command{
		Use:   "env",
		Short: "prints relevant information about your environment",
		Run:   cmdEnv,
	}
	return envCmd
}

func cmdEnv(cmd *cobra.Command, args []string) {
	sc := script.NewContext()

	sc.PrintlnBold("# restic")
	sc.PrintBold("restic ")
	if bkp.ResticIsInstalled() {
		sc.PrintSuccess("installed")
		fmt.Printf(" in %s\n", bkp.ResticPath())

		sc.PrintBold("restic version ")
		v, err := bkp.ResticVersion()
		if err != nil {
			fmt.Printf(" unknown\n")
		} else {
			if v.GE(minResticVersion) {
				sc.PrintfSuccess("%s\n", v)
			} else {
				sc.PrintfError("%s\n", v)
			}
		}
	} else {
		sc.PrintError("not installed!")
	}
}
