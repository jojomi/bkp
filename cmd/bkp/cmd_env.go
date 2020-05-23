package main

import (
	"fmt"

	"github.com/jojomi/bkp"
	"github.com/jojomi/go-script/v2/interview"
	"github.com/jojomi/go-script/v2/print"
	"github.com/rs/zerolog/log"
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
	print.Boldln("# restic")
	fmt.Print("restic ")
	if bkp.ResticIsInstalled() {
		print.Success("installed")
		fmt.Printf(" in %s\n", bkp.ResticPath())

		print.Bold("restic version ")
		v, err := bkp.ResticVersion()
		if err != nil {
			fmt.Printf(" unknown\n")
		} else {
			if v.GE(minResticVersion) {
				print.Successf("%s\n", v)
			} else {
				print.Errorf("%s\n", v)
			}
		}
	} else {
		print.Errorln("not installed!")

		doInstall, err := interview.Confirm("Install restic?", true)
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}
		if doInstall {
			err := bkp.ResticUpdate()
			if err != nil {
				log.Fatal().Err(err).Msg("")
			}
			fmt.Printf("installed to %s.\n", bkp.ResticPath())
		}
	}
}
