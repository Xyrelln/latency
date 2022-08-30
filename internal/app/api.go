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
	"time"

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

func (a *Api) StartRecord(serial string) error {
	log.Printf("start monitor")
	a.VideoDir, a.ImagesDir = utils.CreateWorkDir()
	log.Printf("workdir: %s", a.VideoDir)
	cmd, err := cmd.StartScrcpyRecord(serial, a.VideoDir)
	if err != nil {
		log.Fatal(err)
	}
	a.Cmd = cmd
	a.emitInfo(eventRecordStart)
	return nil
}

func (a *Api) Start(serial string, recordSecond int64) error {
	// timeout := time.After(time.Duration(recordSecond) * time.Second)
	// // err := a.StartRecord(serial)

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
	srcVideoPath := path.Join(a.VideoDir, "rec.mp4")
	a.emitInfo(eventTransformStart)
	cmd, err := cmd.StartVideoToImageTransform(srcVideoPath, a.ImagesDir)
	if err != nil {
		// log.Fatal(err)
		log.Print(err)
		return err
	}
	a.emitInfo(eventTransformFilish)
	a.Cmd = cmd
	return nil
}

func (a *Api) StartAnalyse() error {
	log.Printf("analyse data")
	log.Printf("workdir: %s", a.ImagesDir)
	a.emitInfo(eventAnalyseStart)
	responseTimes, _ := core.CalcTime(a.ImagesDir)
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
