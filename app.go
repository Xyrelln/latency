package main

import (
	"context"
	"errors"
	"log"
	"op-latency-mobile/src/adb"
	"op-latency-mobile/src/cmd"
	"syscall"
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

}

func SetPointerLocationOff(serial string) {

}

func (a *App) StartRecord(serial string) {
	// timeout := 10
	// a.Cmd.Ctx, a.Cmd.Cancel = context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	// a.Cmd.Ctx, a.Cmd.Cancel = context.WithCancel(context.Background())
	// defer a.Cmd.Cancel()
	cmd, _ := cmd.StartScrcpyRecord(serial, "/Users/jason/Developer/epc/op-latency-mobile/out/video/1.mp4")
	a.Cmd = cmd

}

func (a *App) StopRecord(serial string) {
	cmd.StopScrcpyRecord(serial)
}

func (a *App) StopProcessing() {
	// a.Cmd.Cancel()
	log.Printf("Interrupt")
	// signal.NotifyContext(a.Cmd.Ctx, os.Kill, syscall.SIGTERM)
	log.Printf("Kill")
	log.Printf("pid: %d", a.Cmd.Pid)
	_ = syscall.Kill(a.Cmd.ExecCmd.Process.Pid, syscall.SIGINT)
	// cmd.CancelProcess()
}

func StartTransform() {

}

func StartAnalyse() {

}

func ClearImages() {

}

func ClearVideos() {

}
