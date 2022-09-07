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
	log.Printf("workdir: %s", a.VideoDir)
	a.VideoDir, a.ImagesDir = utils.CreateWorkDir()
	recFile := filepath.Join(a.VideoDir, recordFile)
	cmd, err := cmd.StartScrcpyRecord(serial, recFile)
	if err != nil {
		log.Print(err)
		return err
	}
	a.Cmd = cmd
	a.emitInfo(eventRecordStart)
	return nil
}

func (a *Api) Start(serial string, recordSecond int64) error {
	err := a.StartRecord(serial)
	if err != nil {
		a.emitInfo(eventRecordStartError)
		return err
	}
	// wait for filish
	time.Sleep(time.Duration(recordSecond) * time.Second)
	err = a.StopScrcpyServer(serial)
	if err != nil {
		a.emitInfo(eventRecordStopError)
		return err
	}
	err = a.StartTransform()
	if err != nil {
		a.emitInfo(eventTransformStartError)
		return err
	}
	return nil
}

func (a *Api) StopScrcpyServer(serial string) error {
	device := adb.GetDevice(serial)
	err := device.KillScrcyServer()
	if err != nil {
		a.emitInfo(eventRecordStopError)
		return err
	}
	a.emitInfo(eventRecordFilish)
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
		log.Print(err)
		a.emitInfo(eventTransformStartError)
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
	// path to uri format for FileLoader on win
	if utils.IsWindowsDrivePath(mInfo.Path) {
		mInfo.Path = "/" + strings.ReplaceAll(mInfo.Path, "\\", "/")
	}
	return mInfo, nil
}

func (a *Api) StartAnalyse(imageRect core.ImageRectInfo, diffScore int) error {
	log.Printf("current rect: %v", imageRect)
	log.Printf("workdir: %s", a.ImagesDir)
	a.emitInfo(eventAnalyseStart)
	responseTimes, _ := core.CalcTime(a.ImagesDir, imageRect, diffScore)
	a.emitData(eventAnalyseFilish, responseTimes)
	return nil
}

func (a *Api) emitInfo(eventName string) {
	runtime.EventsEmit(a.ctx, eventName)
}

func (a *Api) emitData(eventName string, data interface{}) {
	runtime.EventsEmit(a.ctx, eventName, data)
}

func (a *Api) StopTransform() {
	a.Cmd.Kill()
}

func (a *Api) ClearCacheData() {
	utils.ClearCacheDir()
}
