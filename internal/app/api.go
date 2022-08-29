package app

import (
	"context"
	"errors"
	"log"
	"op-latency-mobile/internal/adb"
	"op-latency-mobile/internal/cmd"
	"op-latency-mobile/internal/core"
	"op-latency-mobile/internal/utils"
	"path"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Api struct
type Api struct {
	ctx       context.Context
	Cmd       cmd.Cmd
	VideoDir  string
	ImagesDir string
}

var Canceled = errors.New("context canceled")

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

func (a *Api) StartRecord(serial string) {
	log.Printf("start monitor")
	a.VideoDir, a.ImagesDir = utils.CreateWorkDir()
	log.Printf("workdir: %s", a.VideoDir)
	cmd, _ := cmd.StartScrcpyRecord(serial, a.VideoDir)
	a.Cmd = cmd
}

func (a *Api) StopRecord(serial string) {
	log.Printf("stop monitor")
	a.Cmd.Kill()
}

func (a *Api) StopProcessing() {
	log.Printf("stop monitor")
	a.Cmd.Kill()
}

func (a *Api) StartTransform() {
	log.Printf("prepare data")
	srcVideoPath := path.Join(a.VideoDir, "rec.mp4")
	runtime.EventsEmit(a.ctx, "latency:StartTransform")
	cmd, _ := cmd.StartVideoToImageTransform(srcVideoPath, a.ImagesDir)
	a.Cmd = cmd
}

func (a *Api) StartAnalyse() error {
	log.Printf("analyse data")
	log.Printf("workdir: %s", a.ImagesDir)
	runtime.EventsEmit(a.ctx, "latency:startAnalyse")
	responseTimes, _ := core.CalcTime(a.ImagesDir)
	done := "latency:done"
	runtime.EventsEmit(a.ctx, done, responseTimes)
	return nil
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
