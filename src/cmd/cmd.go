package cmd

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os/exec"
	"path"
	"path/filepath"
	"syscall"
)

var ErrScrcpyNotFound = errors.New("scrcpy command not found on PATH")
var ErrFfmpegNotFound = errors.New("ffmpeg command not found on PATH")

var scrcpy string
var ffmpeg string

var videoPath = "/Users/jason/Developer/epc/op-latency-mobile/out/video/1.mp4"

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

type Cmd struct {
	Path    string
	Args    []string
	execCmd *exec.Cmd
	Stdout  io.Writer
	Stderr  io.Writer
}

func (c *Cmd) Run() error {
	cmd := exec.Command(ffmpeg, c.Args...)
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	c.execCmd = cmd
	return cmd.Run()
}

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

func (c *Cmd) Start() error {
	cmd := exec.Command(scrcpy, c.Args...)
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}
	c.execCmd = cmd
	return nil
}

func (c *Cmd) FFmpegStart() error {
	cmd := exec.Command(ffmpeg, c.Args...)
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}
	c.execCmd = cmd
	// go cmd.Wait()
	return nil
}

func (c *Cmd) Kill() error {
	if c.execCmd.Process != nil {
		err := syscall.Kill(c.execCmd.Process.Pid, syscall.SIGINT)
		if err != nil {
			return err
		}
	}
	return nil
}

func StartScrcpyRecord(serial string) (Cmd, error) {
	path := "/Users/jason/Developer/epc/op-latency-mobile/out/video/1.mp4"
	cmd := Cmd{
		Args: []string{
			"-s", serial,
			"-r", path,
		},
	}

	err := cmd.Start()
	if err != nil {
		return cmd, err
	}
	return cmd, nil
}

func StartVideoToImageTransform() (Cmd, error) {
	log.Println("prepare ffmpge transform")
	outImgDir := "/Users/jason/Developer/epc/op-latency-mobile/out/image/5"
	outImgPath := path.Join(outImgDir, "%4d.png")
	cmd := Cmd{
		Args: []string{
			"-i", videoPath,
			"-threads", "8",
			outImgPath,
		},
	}
	msg, err := cmd.Call()
	if err != nil {
		log.Printf("start ffmpge failed: %v", err)
		return cmd, err
	}
	log.Println(msg)
	return cmd, nil
}
