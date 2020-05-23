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

	args := mergeStringSlices([]string{j.Source, "--verbose"}, j.Backup.Args)
	ex.Command("backup", args...)
	// TODO unlock if locked? ex.Command("unlock")

	snapshotArgs := make([]string, 0)
	if j.Hostname != "" {
		snapshotArgs = append(snapshotArgs, "--host", j.Hostname)
	}
	if j.Source != "" {
		snapshotArgs = append(snapshotArgs, "--path", strings.TrimRight(j.Source, `/`))
	}
	ex.Command("snapshots", snapshotArgs...)

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
		ex.Command("forget", forgetArgs...)
	}

	if opts.DoMaintenance {
		print.Boldln("restic prune...")
		ex.Command("prune")
		fmt.Println()

		print.Boldln("restic rebuild-index...")
		ex.Command("rebuild-index")
		fmt.Println()

		print.Boldln("restic clean-cache...")
		ex.Command("restic", "cache", "--cleanup")
		fmt.Println()

		print.Boldln("restic check...")
		ex.Command("check")
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
	DoForget      bool
	DoMaintenance bool
}
