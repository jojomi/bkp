package bkp

import (
	"errors"
)

func CheckEnvironment() error {
	// restic?
	if !ResticIsInstalled() {
		return errors.New("restic not installed. Please make sure it is in your PATH.")
	}
	return nil
}
