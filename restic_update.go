package bkp

import (
	"io"
	"net/http"
	"os"

	"github.com/jojomi/go-script/v2"
)

func ResticUpdate() error {
	var packagedRestic = "restic-release/restic_0.9.5_linux_amd64"
	outPath := "/usr/local/bin/restic"
	sc := script.NewContext()
	err := sc.CopyFile(packagedRestic, outPath)
	if err != nil {
		return err
	}

	err = os.Chmod(sc.AbsPath(outPath), os.FileMode(int(0755)))
	if err != nil {
		return err
	}

	return nil
}

func DownloadFile(remote, local string) error {
	out, err := os.Create(local)
	defer out.Close()
	if err != nil {
		return err
	}
	resp, err := http.Get(remote)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
