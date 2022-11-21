//go:build windows
// +build windows

package capture

import (
	"encoding/json"
	"fmt"

	// "image"
	// _ "image/jpeg" // register
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// RustCapture ...
type RustCapture struct {
	ExePath   string
	OutputDir string
	PrintFunc func(string)
}

// CaptureScreenshots ...
func (rsc *RustCapture) CaptureScreenshots(frames int, windowHandle int) (imgs []ScreenshotWithTs, err error) {
	var cmd *exec.Cmd
	cmd = exec.Command(rsc.ExePath,
		"-f", fmt.Sprintf("%d", frames),
		"-o", rsc.OutputDir,
		"--window-handle", fmt.Sprintf("%d", windowHandle),
	)
	// if windowName == "" {
	// 	cmd = exec.Command(rsc.ExePath,
	// 		"-f", fmt.Sprintf("%d", frames),
	// 		"-o", rsc.OutputDir,
	// 	)
	// } else {
	// 	cmd = exec.Command(rsc.ExePath,
	// 		"-f", fmt.Sprintf("%d", frames),
	// 		"-o", rsc.OutputDir,
	// 		"-w", windowName,
	// 	)
	// }

	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return
	}
	return rsc.loadScreenshotsResult(rsc.OutputDir)
}

// WindowInfo ...
type WindowInfo struct {
	Title  string `json:"title"`
	Handle int    `json:"handle"`
}

// ListCapturableWindows ...
func (rsc *RustCapture) ListCapturableWindows() (windowNames []WindowInfo, err error) {
	cmd := exec.Command(rsc.ExePath, "-l")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	b, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	out := strings.TrimSpace(string(b))
	lines := strings.Split(out, "\n")

	res := make([]WindowInfo, 0, len(lines))
	for _, line := range lines {
		var info WindowInfo
		if err := json.Unmarshal([]byte(line), &info); err == nil {
			res = append(res, info)
		} else {
			log.Errorf("unmarshal rscapture window info error: %v", err)
		}
	}
	return res, nil
}

func (rsc *RustCapture) loadScreenshotsResult(dir string) (imgs []ScreenshotWithTs, err error) {
	fentries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	// jobs := make([]readImgJob, 0, len(fentries))
	for _, f := range fentries {
		if filepath.Ext(f.Name()) != ".jpg" {
			continue
		}
		ts := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
		i, err := strconv.ParseInt(ts, 10, 64)
		if err != nil {
			continue
		}
		screenshotTime := time.UnixMilli(i)

		imgFilePath := filepath.Join(dir, f.Name())
		imgs = append(imgs, ScreenshotWithTs{FilePath: imgFilePath, Time: screenshotTime})
	}
	rsc.PrintFunc(fmt.Sprintf("已完成截图，截图数量: %d", len(imgs)))

	sort.Slice(imgs, func(i, j int) bool { return imgs[i].Time.Before(imgs[j].Time) })
	return
}
