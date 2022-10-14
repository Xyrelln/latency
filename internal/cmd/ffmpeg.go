package cmd

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var ffmpeg string
var ErrFFmpegNotFound = errors.New("ffmpeg command not found on PATH")

func init() {
	// Fallback to searching on CurrentDirectory.
	if execPath, err := os.Executable(); err == nil {
		p := filepath.Join(filepath.Dir(execPath), "lib", "ffmpeg", ffmpegExecFile)
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			ffmpeg = p
			return
		}
	}

	// Fallback to searching on PATH.
	if p, err := exec.LookPath(ffmpegExecFile); err == nil {
		if p, err = filepath.Abs(p); err == nil {
			ffmpeg = p
			return
		}
	}

}

func StartFFmpeg(srcVideoPath, destImagePath string) (*Cmd, error) {
	cmd := Cmd{
		Args: []string{
			ffmpeg,
			"-r", "1",
			"-i", srcVideoPath,
			"-r", "1",
			"-threads", "4",
			destImagePath,
		},
	}
	if err := cmd.Run(); err == nil {
		return &cmd, nil
	} else {
		return nil, err
	}
}

func VideoToImage(ctx context.Context, srcVideoPath, destImagePath string) (data []byte, err error) {
	args := []string{
		"-r", "1",
		"-i", srcVideoPath,
		"-r", "1",
		"-threads", "4",
		destImagePath,
	}
	cmd := exec.CommandContext(ctx, ffmpeg, args...)
	cmd.SysProcAttr = procAttributes()

	return CmdRun(cmd)
}

func FFmpegStart(srcVideoPath, destImagePath string, callback CallbackFunc) (*CmdRunner, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)

	go func() {
		VideoToImage(ctx, srcVideoPath, destImagePath)
		callback()
	}()

	return &CmdRunner{Ctx: ctx, CancelFunc: cancelFn}, nil
}

func IsFFmpegReady() error {
	if ffmpeg == "" {
		return ErrFFmpegNotFound
	}
	return nil
}
