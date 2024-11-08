package bkp

import (
	"regexp"

	"github.com/blang/semver"
	"github.com/jojomi/go-script/v2"
)

func ResticPath() string {
	sc := script.NewContext()
	return sc.CommandPath("restic")
}

func ResticIsInstalled() bool {
	return ResticPath() != ""
}

func ResticVersion() (semver.Version, error) {
	sc := script.NewContext()
	lc := script.LocalCommandFrom("restic version")
	pr, err := sc.ExecuteFullySilent(lc)
	rex := regexp.MustCompile(`[0-9]+(\.[0-9]+)?(\.[0-9]+)`)
	versionString := rex.FindString(pr.Output())
	if err != nil {
		v, _ := semver.Make("0.0.0")
		return v, err
	}
	return semver.Make(versionString)
}
