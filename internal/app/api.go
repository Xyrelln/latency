package app

import (
	"context"
	"log"
	"op-latency-mobile/internal/adb"
	"op-latency-mobile/internal/cmd"
	"op-latency-mobile/internal/core"
	"op-latency-mobile/internal/utils"
	"path/filepath"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	recordFile     = "rec.mp4"
	firstImageFile = "0001.png"
)

// Api struct
type Api struct {
	ctx       context.Context
	Cmd       *cmd.Cmd
	VideoDir  string
	ImagesDir string
	// FFmpegFile fs.FS
	// ScrcpyDir  fs.FS
	AppData string
}

// NewApp creates a new Api application struct
func NewApp() *Api {
	return &Api{}
}

// startup is called when the Api starts. The context is saved
// so we can call the runtime methods
func (a *Api) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *Api) ListDevices() ([]*adb.Device, error) {
	devices, err := adb.Devices()
	if err != nil {
		log.Fatalf("ListDevices failed: %v", err)
		return nil, err
	}
	return devices, nil
}

func (a *Api) SetPointerLocationOn(serial string) error {
	log.Printf("set pointer location on")
	device := adb.GetDevice(serial)
	return device.SetPointerLocationOn()
}

func (a *Api) SetPointerLocationOff(serial string) error {
	log.Printf("set pointer location off")
	device := adb.GetDevice(serial)
	return device.SetPointerLocationOff()
}

func (a *Api) StartRecord(serial string) error {
	log.Printf("start monitor")
	a.VideoDir, a.ImagesDir = utils.CreateWorkDir()
	log.Printf("workdir: %s", a.VideoDir)
	recFile := filepath.Join(a.VideoDir, recordFile)
	cmd, err := cmd.StartScrcpyRecord(serial, recFile)
	if err != nil {
		log.Fatal(err)
	}
	a.Cmd = cmd
	a.emitInfo(eventRecordStart)
	return nil
}

func (a *Api) Start(serial string, recordSecond int64) error {
	// timeout := time.After(time.Duration(recordSecond) * time.Second)
	// ch := make(chan string, 1)
	// go func() {
	// 	err := a.StartRecord(serial)
	// 	ch <- err.Error()
	// }()

	// flag := false
	// for {
	// 	select {
	// 	case out := <-ch:
	// 		fmt.Print(out)
	// 		flag = true
	// 	case <-timeout:
	// 		a.StopProcessing()
	// 		fmt.Println("timeout 1")
	// 		flag = true
	// 	}
	// 	if flag {
	// 		break
	// 	}
	// }
	a.StartRecord(serial)
	time.Sleep(time.Duration(recordSecond) * time.Second)
	a.StopProcessing()
	a.StartTransform()
	// a.StartAnalyse()

	return nil

}

func (a *Api) StopRecord(serial string) error {
	log.Printf("stop monitor")
	a.emitInfo(eventRecordFilish)
	return a.Cmd.Kill()

}

func (a *Api) StopProcessing() error {
	log.Printf("stop monitor")
	a.emitInfo(eventRecordFilish)
	return a.Cmd.Kill()
}

func (a *Api) StartTransform() error {
	log.Printf("prepare data")
	srcVideoPath := filepath.Join(a.VideoDir, recordFile)
	destImagePath := filepath.Join(a.ImagesDir, "%4d.png")
	a.emitInfo(eventTransformStart)
	cmd, err := cmd.StartFFmpeg(srcVideoPath, destImagePath)
	if err != nil {
		// log.Fatal(err)
		log.Print(err)
		return err
	}

	a.emitInfo(eventTransformFilish)
	a.Cmd = cmd
	return nil
}

func (a *Api) GetFirstImageInfo() (core.ImageInfo, error) {
	firstImage := filepath.Join(a.ImagesDir, firstImageFile)
	mInfo, err := core.GetImageInfo(firstImage)
	log.Printf("Get first image: %v", mInfo)
	if err != nil {
		return mInfo, err
	}
	if utils.IsWindowsDrivePath(mInfo.Path) {
		mInfo.Path = "/" + strings.ReplaceAll(mInfo.Path, "\\", "/")
	}
	return mInfo, nil
}

func (a *Api) StartAnalyse(imageRect core.ImageRectInfo) error {
	log.Printf("analyse data")
	log.Printf("current rect: %v", imageRect)
	log.Printf("workdir: %s", a.ImagesDir)
	a.emitInfo(eventAnalyseStart)
	responseTimes, _ := core.CalcTime(a.ImagesDir, imageRect)
	a.emitData(eventAnalyseFilish, responseTimes)
	return nil
}

func (a *Api) emitInfo(eventName string) {
	runtime.EventsEmit(a.ctx, eventName)
}

func (a *Api) emitData(eventName string, data interface{}) {
	runtime.EventsEmit(a.ctx, eventName, data)
}

func (a *Api) StopAnalyse() {
}

func (a *Api) StopTransform() {
	a.Cmd.Kill()
}

func (a *Api) ClearImages() {

}

func (a *Api) ClearVideos() {

}

func (a *Api) ClearCacheData() {
	utils.ClearCacheDir()
}
