//go:build darwin
// +build darwin

package cmd

import (
	"os/exec"
	"syscall"
	// "op-latency-mobile/third"
)

func (c *Cmd) BackendRun(name string) error {
	cmd := exec.Command(name, c.Args...)
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	c.execCmd = cmd
	return cmd.Start()
}

func (c *Cmd) Kill() error {
	if c.execCmd.Process != nil {
		if c.execCmd.Process != nil {
			return syscall.Kill(c.execCmd.Process.Pid, syscall.SIGINT)
		}
	}
	return nil
}
