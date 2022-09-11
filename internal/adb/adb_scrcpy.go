package adb

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func (d *Device) KillScrcyServer() error {
	cmd := d.Command("ps -ef | grep scrcpy")

	out, err := cmd.Call()
	if err != nil {
		log.Error(err)
		return err
	}

	// fmt.Print(out)
	pids, err := parsePids(out)
	if err != nil {
		log.Error(err)
		return err

	}
	for _, pid := range pids {
		log.Infof("kill %d", pid)
		d.Command(fmt.Sprintf("kill %d", pid)).Run()
	}
	return nil
}

func parsePids(out string) ([]int, error) {
	var pids []int
	lines := strings.Split(out, "\n")
	if len(lines) <= 2 {
		return nil, errors.New("server not running")
	}
	for _, line := range lines {
		fields := strings.Fields(line)
		// fmt.Printf("line fields: %d", len(fields))
		if len(fields) > 1 {
			if pid, err := strconv.Atoi(fields[1]); err == nil {
				pids = append(pids, pid)
			} else {
				log.Error(err)
				return nil, errors.New("pid read wrong")
			}
		}
	}
	log.Infof("find scrcpy pids: %v", pids)
	return pids, nil
}
