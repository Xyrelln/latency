package cmd

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var ErrScrcpyNotFound = errors.New("scrcpy command not found on PATH")

var scrcpy string

func init() {
	// searching on CurrentDirectory.
	if execPath, err := os.Executable(); err == nil {
		p := filepath.Join(filepath.Dir(execPath), "lib", "scrcpy", scrcpyExecFile)
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			scrcpy = p
			return
		}
	}

	// Fallback to searching on PATH
	if p, err := exec.LookPath(scrcpyExecFile); err == nil {
		if p, err = filepath.Abs(p); err == nil {
			scrcpy = p
			return
		}
	}
}

// func StartScrcpy(serial, recFile string) (*Cmd, error) {
// 	cmd := Cmd{
// 		Args: []string{
// 			scrcpy,
// 			"--codec-options", "bitrate-mode=2", // https://developer.android.com/reference/android/media/MediaCodecInfo.EncoderCapabilities#BITRATE_MODE_CBR
// 			"--no-cleanup",
// 			"-s", serial,
// 			"--max-fps", "60",
// 			// "-n", // no-control
// 			"-w", // stay awake
// 			"-Nr", recFile,
// 		},
// 	}

// 	if err := cmd.BackendRun(); err == nil {
// 		return &cmd, nil
// 	} else {
// 		return nil, err
// 	}
// }

func RecordVideo(ctx context.Context, serial, recFile string) (data []byte, err error) {
	args := []string{
		"--codec-options", "bitrate-mode=2", // https://developer.android.com/reference/android/media/MediaCodecInfo.EncoderCapabilities#BITRATE_MODE_CBR
		"--no-cleanup",
		"-s", serial,
		"--max-fps", "60",
		// "-n", // no-control
		"-w", // stay awake
		"-Nr", recFile,
	}
	cmd := exec.CommandContext(ctx, scrcpy, args...)
	cmd.SysProcAttr = procAttributes()

	return CmdRun(cmd)
}

func ScrcpyStart(serial, recFile string, recordSecond int) (*CmdRunner, error) {
	var timeOut = 20
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Duration(recordSecond+timeOut)*time.Second)

	go func() {
		RecordVideo(ctx, serial, recFile)
	}()

	return &CmdRunner{Ctx: ctx, CancelFunc: cancelFn}, nil
}

func IsScrcpyReady() (string, error) {
	if scrcpy == "" {
		return "", ErrScrcpyNotFound
	}
	return scrcpy, nil
}
