//go:build windows
// +build windows

package app

import (
	// "fmt"
	"op-latency-mobile/internal/core"
	latencywin "op-latency-mobile/internal/latency_win"

	log "github.com/sirupsen/logrus"
)

// StartWinOpLatency ...
func (a *Api) StartWinOpLatency(config latencywin.Config) (imageCount int, err error) {
	a.latencyWinManager = &latencywin.OpLatencyWindowsManager{}
	err = a.latencyWinManager.Start(config)
	if err != nil {
		return
	}
	imageCount = len(a.latencyWinManager.ScreenshotSeq)
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
	ImageCount    int    `json:"length"`
	CurrentIndex  int    `json:"currentIndex"`
	ImageFilePath string `json:"imageFilePath"`
	ImageWidth    int    `json:"imageWidth"`
	ImageHeight   int    `json:"imageHeight"`
}

// GetImage ...
func (a *Api) GetImage(index int) GetImageResp {
	if a.latencyWinManager == nil {
		return GetImageResp{}
	}
	if index < 0 || index >= len(a.latencyWinManager.ScreenshotSeq) {
		return GetImageResp{}
	}

	screenshot := a.latencyWinManager.ScreenshotSeq[index]
	img, err := screenshot.DecodeImg()
	if err != nil {
		log.Errorf("decode %s error: %w", screenshot.FilePath, err)
		return GetImageResp{}
	}

	return GetImageResp{
		ImageCount:    len(a.latencyWinManager.ScreenshotSeq),
		CurrentIndex:  index,
		ImageFilePath: screenshot.FilePath,
		ImageWidth:    img.Bounds().Dx(),
		ImageHeight:   img.Bounds().Dy(),
	}
}
