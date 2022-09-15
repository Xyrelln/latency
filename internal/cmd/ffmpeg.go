package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var ffmpeg string
var ErrFFmpegNotFound = errors.New("ffmpeg command not found on PATH")

func init() {
	// Fallback to searching on PATH.
	//if p, err := exec.LookPath(ffmpegExecFile); err == nil {
	//	if p, err = filepath.Abs(p); err == nil {
	//		ffmpeg = p
	//		return
	//	}
	//}

	// Fallback to searching on CurrentDirectory.
	if execPath, err := os.Executable(); err == nil {
		p := filepath.Join(filepath.Dir(execPath), "lib", "ffmpeg", ffmpegExecFile)
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			ffmpeg = p
			return
		} else {
			fmt.Errorf("ffmpeg path check failed: %s, reason: %v ", p, err)
			// log.Errorf("ffmpeg path check failed: %s, reason: v%", p, err)
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

func IsFFmpegReady() error {
	if ffmpeg == "" {
		return ErrFFmpegNotFound
	}
	return nil
}
