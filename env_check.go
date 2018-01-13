package bkp

import (
	"errors"
	"fmt"

	"github.com/blang/semver"
)

func CheckEnvironment(minResticVersion semver.Version) error {
	// restic?
	if !ResticIsInstalled() {
		return errors.New("restic not installed. Please make sure it is in your PATH.")
	}
	resticVersion, err := ResticVersion()
	if err != nil {
		return fmt.Errorf("restic version not detectable.")
	}
	if resticVersion.LT(minResticVersion) {
		return fmt.Errorf("restic version should be >=%s, but version %s detected.", minResticVersion, resticVersion)
	}
	return nil
}
