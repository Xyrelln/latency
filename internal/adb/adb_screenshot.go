package adb

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func (d *Device) GetScreenshot() (string, error) {
	innerTmpPath := "/sdcard/tmp001.png"
	cmd := d.Command(fmt.Sprintf("rm %s", innerTmpPath))
	err := cmd.Run()
	if err != nil {
		log.Error(err) // skip error for No such file or directory
		// return "", err
	}

	cmd = d.Command(fmt.Sprintf("screencap -p %s", innerTmpPath))
	err = cmd.Run()
	if err != nil {
		log.Error(err)
		return "", err
	}
	return innerTmpPath, nil
}
