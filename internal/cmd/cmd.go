package cmd

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"io"
	"os/exec"
	"strings"
)


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
func (c *Cmd) Run() error {
	cs := append(cmdStart, strings.Join(c.Args, " "))
	cmd := exec.Command(cs[0], cs[1:]...)
	cmd.SysProcAttr = procAttrs
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	c.execCmd = cmd
	log.Infof("args: %v", c.Args)
	return cmd.Run()
}

// BackendRun Start without wait
func (c *Cmd) BackendRun() error {
	cs := append(cmdStart, strings.Join(c.Args, " "))
	cmd := exec.Command(cs[0], cs[1:]...)
	cmd.SysProcAttr = procAttrs
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	c.execCmd = cmd
	log.Infof("args: %v", c.Args)
	return cmd.Start()
}

// Call starts the specified command and waits for it to complete, returning the
// all stdout as a string.
// The returned error is nil if the command runs, has no problems copying
// stdout and stderr, and exits with a zero exit status.
func (c *Cmd) Call() (string, error) {
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
	err := clone.Run()
	return stdout.String(), err
}
