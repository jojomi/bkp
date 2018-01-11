package bkp

import (
	"fmt"
	"os"

	script "github.com/jojomi/go-script"
)

type Job struct {
	Name       string `json:"name"`
	Filename   string
	Weight     int      `json:"weight"` // lower values mean the job will be executed earlier
	Hostname   string   `json:"hostname"`
	IP         string   `json:"ip"`
	Source     string   `json:"source"`
	Args       []string `json:"args"`
	TargetName string   `json:"target"`
	Target     *Target
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

	ex := NewResticExecutor()
	ex.SetTarget(j.Target)
	ex.DryRun = opts.DryRun

	context := script.NewContext()
	context.PrintlnBold(fmt.Sprintf("Backup %s...", j.Name))
	args := mergeStringSlices([]string{j.Source}, j.Args)
	ex.Command("backup", args...)

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
	DryRun bool
}
