package cmd

import (
	"bytes"
	"errors"
	log "github.com/sirupsen/logrus"
	"io"
	"os/exec"
	// "op-latency-mobile/third"
)

var ErrScrcpyNotFound = errors.New("scrcpy command not found on PATH")
var ErrFFmpegNotFound = errors.New("ffmpeg command not found on PATH")
var ErrTaskKillNotFound = errors.New("taskkill command not found on PATH")

var scrcpy string
var ffmpeg string

type Cmd struct {
	Path    string
	Args    []string
	execCmd *exec.Cmd
	Stdout  io.Writer
	Stderr  io.Writer
}

// Run starts the specified command and waits for it to complete.
// The returned error is nil if the command runs, has no problems copying
// stdout and stderr, and exits with a zero exit status.
func (c *Cmd) Run(name string) error {
	cmd := exec.Command(name, c.Args...)
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	c.execCmd = cmd
	log.Printf("cmd: %s", name)
	log.Printf("args: %v", c.Args)
	return cmd.Run()
}

// Start without wait
// func (c *Cmd) BackendRun(name string) error {
// 	cmd := exec.Command(name, c.Args...)
// 	cmd.Stdout = c.Stdout
// 	cmd.Stderr = c.Stderr
// 	c.execCmd = cmd
// 	return cmd.Start()
// }

// Call starts the specified command and waits for it to complete, returning the
// all stdout as a string.
// The returned error is nil if the command runs, has no problems copying
// stdout and stderr, and exits with a zero exit status.
func (c *Cmd) Call(name string) (string, error) {
	clone := *c
	stdout := &bytes.Buffer{}
	if clone.Stdout != nil {
		clone.Stdout = io.MultiWriter(clone.Stdout, stdout)
	} else {
		clone.Stdout = stdout
	}
	stderr := &bytes.Buffer{}
	if clone.Stdout != nil {
		clone.Stderr = io.MultiWriter(clone.Stdout, stderr)
	} else {
		clone.Stderr = stderr
	}
	err := clone.Run(name)
	return stdout.String(), err
}

func StartFFmpeg(srcVideoPath, destImagePath string) (*Cmd, error) {
	if ffmpeg == "" {
		return nil, ErrFFmpegNotFound
	}
	cmd := Cmd{
		Args: []string{
			"-i", srcVideoPath,
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
