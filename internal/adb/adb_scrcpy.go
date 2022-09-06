package adb

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func (d *Device) KillScrcyServer() error {
	cmd := d.Command("ps -ef | grep scrcpy")
	if out, err := cmd.Call(); err == nil {
		fmt.Print(out)
		pids, err := parsePids(out)
		if err == nil {
			for _, pid := range pids {
				fmt.Printf("kill %d", pid)
				d.Command(fmt.Sprintf("kill %d", pid)).Run()
			}
		} else {
			log.Print(err)
		}
	} else {
		log.Print(err)
		return err
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
				log.Print(err)
				return nil, errors.New("pid read wrong")
			}
		}
	}
	fmt.Printf("find scrcpy pids: %v", pids)
	return pids, nil
}
