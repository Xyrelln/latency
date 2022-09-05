//go:build darwin
// +build darwin

package cmd

import (
	"log"
	"os/exec"
	"path/filepath"
	"syscall"
	// "op-latency-mobile/third"
)

func init() {
	if p, err := exec.LookPath("scrcpy"); err == nil {
		if p, err = filepath.Abs(p); err == nil {
			scrcpy = p
		}
	}

	if p, err := exec.LookPath("ffmpeg"); err == nil {
		if p, err = filepath.Abs(p); err == nil {
			ffmpeg = p
		}
	}

}
func (c *Cmd) BackendRun(name string) error {
	cmd := exec.Command(name, c.Args...)
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	c.execCmd = cmd
	log.Printf("cmd: %s", name)
	log.Printf("args: %v", c.Args)
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

func StartScrcpyRecord(serial, recFile string) (*Cmd, error) {
	if scrcpy == "" {
		return nil, ErrScrcpyNotFound
	}
	cmd := Cmd{
		Args: []string{
			"-s", serial,
			"-Nr", recFile,
		},
	}

	if err := cmd.BackendRun(scrcpy); err == nil {
		return &cmd, nil
	} else {
		return nil, err
	}
}