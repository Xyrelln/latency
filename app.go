package main

import (
	"context"
	"errors"
	"log"
	"op-latency-mobile/src/adb"
	"op-latency-mobile/src/cmd"
	"op-latency-mobile/src/core"
	"op-latency-mobile/src/utils"
	"path"
)

// App struct
type App struct {
	ctx       context.Context
	Cmd       cmd.Cmd
	VideoDir  string
	ImagesDir string
}

var Canceled = errors.New("context canceled")

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) ListDevices() ([]*adb.Device, error) {
	devices, err := adb.Devices()
	if err != nil {
		log.Fatalf("ListDevices failed: %v", err)
		return nil, err
	}
	return devices, nil
}

func (a *App) SetPointerLocationOn(serial string) {
	log.Printf("set pointer location on")
	device := adb.GetDevice(serial)
	device.SetPointerLocationOn()
}

func (a *App) SetPointerLocationOff(serial string) {
	log.Printf("set pointer location off")
	device := adb.GetDevice(serial)
	device.SetPointerLocationOn()
}

func (a *App) StartRecord(serial string) {
	log.Printf("start monitor")
	a.VideoDir, a.ImagesDir = utils.CreateWorkDir()
	log.Printf("workdir: %s", a.VideoDir)
	cmd, _ := cmd.StartScrcpyRecord(serial, a.VideoDir)
	a.Cmd = cmd
}

func (a *App) StopRecord(serial string) {
	log.Printf("stop monitor")
	a.Cmd.Kill()
}

func (a *App) StopProcessing() {
	log.Printf("stop monitor")
	a.Cmd.Kill()
}

func (a *App) StartTransform() {
	log.Printf("prepare data")
	srcVideoPath := path.Join(a.VideoDir, "rec.mp4")
	cmd, _ := cmd.StartVideoToImageTransform(srcVideoPath, a.ImagesDir)
	a.Cmd = cmd
}

func (a *App) StartAnalyse() []int {
	log.Printf("analyse data")
	responseTimes, _ := core.CalcTime(a.ImagesDir)
	return responseTimes
}

func (a *App) StopAnalyse() {
}

func (a *App) StopTransform() {
	a.Cmd.Kill()
}

func (a *App) ClearImages() {

}

func (a *App) ClearVideos() {

}
