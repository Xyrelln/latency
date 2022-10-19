package capture

import (
	"fmt"
	// "image"
	// _ "image/jpeg" // register
	"os"
	"os/exec"
	"path/filepath"
	// "runtime"
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
	startTime time.Time
}

// CaptureScreenshots ...
func (rsc *RustCapture) CaptureScreenshots(frames int) (imgs []ScreenshotWithTs, err error) {
	cmd := exec.Command(rsc.ExePath,
		"-f", fmt.Sprintf("%d", frames),
		"-o", rsc.OutputDir,
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	rsc.startTime = time.Now()
	err = cmd.Run()
	if err != nil {
		return
	}
	return rsc.loadScreenshotsResult()
}

// type readImgJob struct {
// 	path string
// 	t    time.Time
// }

// type readImgResult struct {
// 	screenshotWithTs *ScreenshotWithTs
// 	err              error
// }

func (rsc *RustCapture) loadScreenshotsResult() (imgs []ScreenshotWithTs, err error) {
	fentries, err := os.ReadDir(rsc.OutputDir)
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
		if !screenshotTime.After(rsc.startTime) {
			continue
		}

		imgFilePath := filepath.Join(rsc.OutputDir, f.Name())
		imgs = append(imgs, ScreenshotWithTs{FilePath: imgFilePath, Time: screenshotTime})
	}
	rsc.PrintFunc(fmt.Sprintf("截图数量: %d", len(imgs)))

	// resultCh := make(chan readImgResult, len(jobs))
	// jobCh := make(chan readImgJob, len(jobs))

	// wnum := runtime.NumCPU() - 1
	// for w := 0; w < wnum; w++ {
	// 	go func() {
	// 		for j := range jobCh {
	// 			imgf, err := os.Open(j.path)
	// 			if err != nil {
	// 				resultCh <- readImgResult{err: err}
	// 				return
	// 			}
	// 			defer imgf.Close()

	// 			img, _, err := image.Decode(imgf)
	// 			if err != nil {
	// 				resultCh <- readImgResult{err: err}
	// 				return
	// 			}
	// 			resultCh <- readImgResult{screenshotWithTs: &ScreenshotWithTs{Img: img, Time: j.t, FilePath: j.path}}
	// 		}
	// 	}()
	// }

	// for _, job := range jobs {
	// 	jobCh <- job
	// }
	// close(jobCh)

	// rsc.PrintFunc("正在加载截图...")
	// for i := 0; i < len(jobs); i++ {
	// 	r := <-resultCh
	// 	if r.screenshotWithTs != nil && r.err == nil {
	// 		imgs = append(imgs, *r.screenshotWithTs)
	// 	}
	// }

	sort.Slice(imgs, func(i, j int) bool { return imgs[i].Time.Before(imgs[j].Time) })
	return
}
