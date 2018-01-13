package bkp

import (
	"github.com/blang/semver"
	script "github.com/jojomi/go-script"
)

func ResticIsInstalled() bool {
	sc := script.NewContext()
	return sc.CommandExists("restic")
}

func ResticVersion() (semver.Version, error) {
	sc := script.NewContext()
	pr, err := sc.ExecuteFullySilent("restic", "version")
	_ = pr
	if err != nil {
		v, _ := semver.Make("0.0.0")
		return v, err
	}
	return semver.Make("0.0.1")
}
