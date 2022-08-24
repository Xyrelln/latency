package main

import (
	"context"
	"fmt"
	"log"
	"op-latency-mobile/src/adb"
	"op-latency-mobile/src/cmd"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
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

}

func SetPointerLocationOff(serial string) {

}

func (a *App) StartRecord(serial string) {
	cmd.StartScrcpyRecord(serial, "/Users/jason/Developer/epc/op-latency-mobile/out/video/1.mp4")
}

func (a *App) StopRecord(serial string) {
	cmd.StartScrcpyRecord(serial, "/Users/jason/Developer/epc/op-latency-mobile/out/video/1.mp4")
}

func StartTransform() {

}

func StartAnalyse() {

}

func ClearImages() {

}

func ClearVideos() {

}
