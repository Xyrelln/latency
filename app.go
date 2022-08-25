package main

import (
	"context"
	"errors"
	"log"
	"op-latency-mobile/src/adb"
	"op-latency-mobile/src/cmd"
	"op-latency-mobile/src/core"
)

// App struct
type App struct {
	ctx context.Context
	Cmd cmd.Cmd
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

func SetPointerLocationOn(serial string) {
	device := adb.GetDevice(serial)
	device.SetPointerLocationOn()
}

func SetPointerLocationOff(serial string) {
	device := adb.GetDevice(serial)
	device.SetPointerLocationOn()
}

func (a *App) StartRecord(serial string) {
	cmd, _ := cmd.StartScrcpyRecord(serial)
	a.Cmd = cmd
}

func (a *App) StopRecord(serial string) {
	a.Cmd.Kill()
}

func (a *App) StopProcessing() {
	a.Cmd.Kill()
}

func (a *App) StartTransform() {
	cmd, _ := cmd.StartVideoToImageTransform()
	a.Cmd = cmd
}

func (a *App) StartAnalyse() {
	core.CalcTime()
	// core.ListImageFile("/Users/jason/Developer/epc/op-latency-mobile/out/image/167-png/")
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
