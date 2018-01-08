package bkp

import (
	"errors"
)

func CheckEnvironment() error {
	// restic?
	if !ResticIsInstalled() {
		return errors.New("restic not installed. Please make sure it it ")
	}
	// warn about nice (Linux, MacOS X) and ionice (Linux)

	return nil
}
