//go:build windows
// +build windows

package app

import (
	"fmt"
	"op-latency-mobile/internal/core"
	latencywin "op-latency-mobile/internal/latency_win"

	log "github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// LatencyWindowsCompleteEvent ...
type LatencyWindowsCompleteEvent struct {
	ImageCount int    `json:"imageCount"`
	InputTime  string `json:"inputTime"`
}

// LatencyWindowsMessageEvent ...
type LatencyWindowsMessageEvent struct {
	Message string `json:"message"`
}

// StartWinOpLatency ...
func (a *Api) StartWinOpLatency(config latencywin.Config) {
	if a.latencyWinManager == nil {
		return
	}

	go func() {
		err := a.latencyWinManager.Start(config, func(s string) {
			runtime.EventsEmit(a.ctx, "latencyWindowsMessage", LatencyWindowsMessageEvent{Message: s})
		})
		if err != nil {
			return
		}
		imageCount := a.latencyWinManager.GetScreenshotCount()
		inputTime := a.latencyWinManager.GetInputTime()

		runtime.EventsEmit(a.ctx, "latencyWindowsComplete", LatencyWindowsCompleteEvent{ImageCount: imageCount, InputTime: fmt.Sprintf("%d", inputTime.UnixMilli())})

	}()
	return
}

// WinOpLatencyResult ...
type WinOpLatencyResult struct {
	Latency       int `json:"latency"`
	ResponseIndex int `json:"responseIndex"`
	ResponseTime  int `json:"responseTime"`
}

// CalculateLatencyByImageDiff ...
func (a *Api) CalculateLatencyByImageDiff(imageRect core.ImageRectInfo) (result WinOpLatencyResult) {
	if a.latencyWinManager == nil {
		return
	}

	respIndex, responseTime, latency, err := a.latencyWinManager.CalculateLatencyByImageDiff(imageRect)
	if err != nil {
		log.Errorf("calculate latency error: %v", err)
		return
	}
	return WinOpLatencyResult{
		Latency:       int(latency.Milliseconds()),
		ResponseIndex: respIndex,
		ResponseTime:  int(responseTime.UnixMilli()),
	}
}

// CalculateLatencyByCurrentImage ...
func (a *Api) CalculateLatencyByCurrentImage(currenIndex int) (result WinOpLatencyResult) {
	if a.latencyWinManager == nil {
		return
	}

	respIndex, responseTime, latency, err := a.latencyWinManager.CalculateLatencyByIndex(currenIndex)
	if err != nil {
		log.Errorf("calculate latency error: %v", err)
		return
	}
	return WinOpLatencyResult{
		Latency:       int(latency.Milliseconds()),
		ResponseIndex: respIndex,
		ResponseTime:  int(responseTime.UnixMilli()),
	}
}

// GetImageResp ...
type GetImageResp struct {
	ImageCount     int    `json:"length"`
	CurrentIndex   int    `json:"currentIndex"`
	ScreenshotTime string `json:"screenshotTime"`
	ImageFilePath  string `json:"imageFilePath"`
	ImageWidth     int    `json:"imageWidth"`
	ImageHeight    int    `json:"imageHeight"`
}

// GetImage ...
func (a *Api) GetImage(index int) GetImageResp {
	if a.latencyWinManager == nil {
		return GetImageResp{}
	}

	imageCount := a.latencyWinManager.GetScreenshotCount()
	if index < 0 || index >= imageCount {
		return GetImageResp{}
	}

	screenshot := a.latencyWinManager.GetScreenshotByIndex(index)
	img, err := screenshot.DecodeImg()
	if err != nil {
		log.Errorf("decode %s error: %w", screenshot.FilePath, err)
		return GetImageResp{}
	}

	return GetImageResp{
		ImageCount:     imageCount,
		CurrentIndex:   index,
		ImageFilePath:  screenshot.FilePath,
		ScreenshotTime: fmt.Sprintf("%d", screenshot.Time.UnixMilli()),
		ImageWidth:     img.Bounds().Dx(),
		ImageHeight:    img.Bounds().Dy(),
	}
}
