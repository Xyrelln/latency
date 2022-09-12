//go:build darwin
// +build darwin

package cmd

import (
	//"fmt"
	"os/exec"
	"path/filepath"
	"syscall"

	log "github.com/sirupsen/logrus"
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

func (c *Cmd) Run(name string) error {
	cmd := exec.Command(name, c.Args...)
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	c.execCmd = cmd
	log.Printf("cmd: %s", name)
	log.Printf("args: %v", c.Args)
	return cmd.Run()
}

func (c *Cmd) BackendRun(name string) error {
	cmd := exec.Command(name, c.Args...)
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	c.execCmd = cmd
	log.Infof("cmd: %s", name)
	log.Infof("args: %v", c.Args)
	return cmd.Start()
}

func (c *Cmd) Kill() error {
	if c.execCmd.Process != nil {
		log.Info("kill proces")
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
			"--max-fps", "60",
			"-Nr", recFile,
			// "-n", // no-control
			// "-w", // stay awake
		},
	}

	if err := cmd.BackendRun(scrcpy); err == nil {
		return &cmd, nil
	} else {
		return nil, err
	}
}

func StartFFmpeg(srcVideoPath, destImagePath string) (*Cmd, error) {
	if ffmpeg == "" {
		return nil, ErrFFmpegNotFound
	}
	cmd := Cmd{
		Args: []string{
			// "-r", "1",
			"-i", srcVideoPath,
			// "-r", "1",
			"-threads", "4",
			destImagePath,
		},
	}
	if err := cmd.Run(ffmpeg); err == nil {
		return &cmd, nil
	} else {
		return nil, err
	}
}
