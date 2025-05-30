package bkp

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jojomi/go-script/v2/print"
)

type Job struct {
	Name     string `json:"name"`
	Filename string
	Weight   int    `json:"weight"` // lower values mean the job will be executed earlier
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	Source   string `json:"source"`
	CacheDir string `json:"cache-dir,omitempty"`
	Backup   struct {
		Args []string `json:"args"`
	}
	TargetName string `json:"target"`
	Target     *Target
	Forget     struct {
		Keep struct {
			Hourly  *int
			Daily   *int
			Weekly  *int
			Monthly *int
			Yearly  *int
		}
	}
}

func (j *Job) IsRelevant() bool {
	if j.Hostname != "" {
		hostname, err := os.Hostname()
		if err == nil && hostname != j.Hostname {
			return false
		}
	}

	return true
}

func (j *Job) Execute(opts JobExecuteOptions) error {
	print.Boldln(fmt.Sprintf("Backup %s...", j.Name))

	ex := NewResticExecutor()
	ex.SetTarget(j.Target)
	ex.SetCacheDir(j.CacheDir)
	ex.DryRun = opts.DryRun

	if opts.DoUnlock {
		print.Boldln("restic unlock...")
		_, err := ex.Command("unlock")
		if err != nil {
			return err
		}
		fmt.Println()
	}

	args := mergeStringSlices([]string{j.Source, "--verbose"}, j.Backup.Args)
	_, err := ex.Command("backup", args...)
	if err != nil {
		return err
	}

	snapshotArgs := make([]string, 0)
	if j.Hostname != "" {
		snapshotArgs = append(snapshotArgs, "--host", j.Hostname)
	}
	if j.Source != "" {
		snapshotArgs = append(snapshotArgs, "--path", strings.TrimRight(j.Source, `/`))
	}
	_, err = ex.Command("snapshots", snapshotArgs...)
	if err != nil {
		return err
	}

	if opts.DoForget {
		forgetArgs := []string{
			"--verbose",
		}
		if j.Hostname != "" {
			forgetArgs = append(forgetArgs, "--host", j.Hostname)
		}
		if j.Source != "" {
			forgetArgs = append(forgetArgs, "--path", j.Source)
		}
		if j.Forget.Keep.Hourly != nil {
			forgetArgs = append(forgetArgs, "--keep-hourly", strconv.Itoa(*j.Forget.Keep.Hourly))
		}
		if j.Forget.Keep.Daily != nil {
			forgetArgs = append(forgetArgs, "--keep-daily", strconv.Itoa(*j.Forget.Keep.Daily))
		}
		if j.Forget.Keep.Weekly != nil {
			forgetArgs = append(forgetArgs, "--keep-weekly", strconv.Itoa(*j.Forget.Keep.Weekly))
		}
		if j.Forget.Keep.Monthly != nil {
			forgetArgs = append(forgetArgs, "--keep-monthly", strconv.Itoa(*j.Forget.Keep.Monthly))
		}
		if j.Forget.Keep.Yearly != nil {
			forgetArgs = append(forgetArgs, "--keep-yearly", strconv.Itoa(*j.Forget.Keep.Yearly))
		}
		_, err := ex.Command("forget", forgetArgs...)
		if err != nil {
			return err
		}
	}

	if opts.DoMaintenance {
		print.Boldln("restic prune...")
		_, err := ex.Command("prune")
		if err != nil {
			return err
		}
		fmt.Println()

		print.Boldln("restic rebuild-index...")
		_, err = ex.Command("rebuild-index")
		if err != nil {
			return err
		}
		fmt.Println()

		print.Boldln("restic clean-cache...")
		_, err = ex.Command("restic", "cache", "--cleanup")
		if err != nil {
			return err
		}
		fmt.Println()

		print.Boldln("restic check...")
		_, err = ex.Command("check")
		if err != nil {
			return err
		}
		fmt.Println()
	}

	/*if flagCheck {
		context.PrintlnBold("Konsistenz prüfen...")
		executeResticCommand(context, "check")

		// Im Moment wegen Bugs nötig, kann entfernt werden, wenn check stabil gut ist
		context.PrintlnBold("Index erneuern...")
		executeResticCommand(context, "rebuild-index")
	}

	context.PrintlnBold("Alte Snapshots nach Policy löschen...")
	executeResticCommand(context, "forget", "--keep-daily", "14", "--keep-weekly", "10", "--keep-monthly", "24", "--keep-yearly", "50")

	fmt.Println()
	context.PrintlnBold("Aktuelle Snapshots")
	executeResticCommand(context, "snapshots")

	fmt.Println()
	context.PrintlnBold("Speicherverbrauch")
	context.ExecuteDebug("du", "-sh", j.Target.Path)*/

	return nil
}

func (j *Job) String() string {
	return fmt.Sprintf("\"%s\" to \"%s\" (defined in %s)", j.Name, j.TargetName, j.Filename)
}

type JobExecuteOptions struct {
	DryRun        bool
	DoUnlock      bool
	DoForget      bool
	DoMaintenance bool
}
