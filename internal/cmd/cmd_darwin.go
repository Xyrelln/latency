//go:build darwin
// +build darwin

package cmd

import (
	"syscall"
	// "op-latency-mobile/third"
)

func (c *Cmd) Kill() error {
	if c.execCmd.Process != nil {
		if c.execCmd.Process != nil {
			return syscall.Kill(c.execCmd.Process.Pid, syscall.SIGINT)
		}
	}
	return nil
}
