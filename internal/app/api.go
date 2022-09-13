package app

import (
	"context"
	"errors"
	"fmt"
	"op-latency-mobile/internal/adb"
	"op-latency-mobile/internal/cmd"
	"op-latency-mobile/internal/core"
	"op-latency-mobile/internal/utils"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	recordFile     = "rec.mp4"
	firstImageFile = "0001.png"
)

var autorun = true

// Api struct
type Api struct {
	ctx       context.Context `json:"ctx,omitempty"`
	Cmd       *cmd.Cmd        `json:"cmd,omitempty"`
	VideoDir  string          `json:"video_dir,omitempty"`
	ImagesDir string          `json:"images_dir,omitempty"`
	AppData   string          `json:"app_data,omitempty"`
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

// 获取设备列表
func (a *Api) ListDevices() ([]*adb.Device, error) {
	devices, err := adb.Devices()
	if err != nil {
		log.Errorf("ListDevices failed: %v", err)
		return nil, err
	}
	return devices, nil
}

// 检查 app 环境信息
func (a *Api) IsAppReady() error {
	ok := adb.IsAdbReady()
	if !ok {
		log.Error("adb path wrong")
		return errors.New("adb path wrong")
	}

	err := cmd.IsFFmpegReady()
	if err != nil {
		log.Error("ffmpeg path wrong")
		return err
	}

	err = cmd.IsScrcpyReady()
	if err != nil {
		log.Error("scrcpy path wrong")
		return err
	}
	return nil
}

// 获取屏幕信息
func (a *Api) GetDisplay(serial string) (*adb.Display, error) {
	device := adb.GetDevice(serial)
	display, err := device.DisplaySize()
	if err != nil {
		return nil, err
	}
	return display, nil
}

func (a *Api) SetAutoSwipeOn(sw adb.SwipeEvent, interval int) error {
	autorun = true
	devices, err := adb.Devices()
	if err != nil {
		log.Errorf("ListDevices failed: %v", err)
		return err
	}
	for {
		for _, device := range devices {
			go device.AutoSwipe(sw)
		}
		time.Sleep(time.Duration(interval) * time.Millisecond)
		if !autorun {
			break
		}
	}

	return nil
}

// 发送拖动事件
func (a *Api) InputSwipe(serial string, sw adb.SwipeEvent) error {
	device := adb.GetDevice(serial)
	return device.InputSwipe(sw)
}

func (a *Api) SetAutoSwipeOff() error {
	autorun = false
	return nil
}

// 开启指针位置显示
func (a *Api) SetPointerLocationOn(serial string) error {
	device := adb.GetDevice(serial)
	if err := device.SetPointerLocationOn(); err == nil {
		on, err := device.IsPointerLocationOn()
		if err != nil {
			return err
		}
		if !on {
			return errors.New("set pointer location failed")
		}
	} else {
		return err
	}

	return nil
}

// 关闭指针位置显示
func (a *Api) SetPointerLocationOff(serial string) error {
	log.Info("set pointer location off")
	device := adb.GetDevice(serial)
	return device.SetPointerLocationOff()
}

func (a *Api) StartRecord(serial string) error {
	log.Info("start monitor")
	log.Infof("workdir: %s", a.VideoDir)
	a.VideoDir, a.ImagesDir = utils.CreateWorkDir()
	recFile := filepath.Join(a.VideoDir, recordFile)
	cmd, err := cmd.StartScrcpy(serial, recFile)
	if err != nil {
		log.Error(err)
		return err
	}
	a.Cmd = cmd
	a.emitInfo(eventRecordStart)
	return nil
}

// 启动延迟测试
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

// 停止 scrcpy server
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

//func (a *Api) StopRecord(serial string) error {
//	log.Infof("stop monitor")
//	a.emitInfo(eventRecordFilish)
//	return a.Cmd.Kill()
//
//}

//func (a *Api) StopProcessing() error {
//	log.Printf("stop monitor")
//	a.emitInfo(eventRecordFilish)
//
//	return a.Cmd.Kill()
//}

func (a *Api) StartTransform() error {
	log.Infof("prepare data")
	srcVideoPath := filepath.Join(a.VideoDir, recordFile)
	destImagePath := filepath.Join(a.ImagesDir, "%4d.png")
	a.emitInfo(eventTransformStart)
	cmd, err := cmd.StartFFmpeg(srcVideoPath, destImagePath)
	if err != nil {
		log.Error(err)
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
	log.Infof("Get first image: %v", mInfo)
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
	log.Infof("current rect: %v", imageRect)
	log.Infof("workdir: %s", a.ImagesDir)
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

//func (a *Api) StopTransform() {
//	a.Cmd.Kill()
//}

func (a *Api) GetImageFiles() ([]string, error) {
	var imgs []string
	imgs, err := utils.GetImageFiles(a.ImagesDir, imgs)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return imgs, nil
}

func (a *Api) ClearCacheData() {
	utils.ClearCacheDir()
}
