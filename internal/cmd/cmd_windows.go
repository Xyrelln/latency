//go:build windows
// +build windows

package cmd

import (
	"strconv"
)

func (c *Cmd) TaskKill(pid string) error {
	if taskkill == "" {
		return ErrTaskKillNotFound
	}
	cmd := Cmd{
		Args: []string{
			"/pid", pid,
		},
	}
	if err := cmd.Run(taskkill); err == nil {
		return nil
	} else {
		return err
	}
}

func (c *Cmd) Kill() error {
	if c.execCmd.Process != nil {
		// https://github.com/golang/go/issues/46345  windows not implemented signal
		return c.TaskKill(strconv.Itoa(c.execCmd.Process.Pid))
	}
	return nil
}
