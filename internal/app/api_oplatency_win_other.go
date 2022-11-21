//go:build !windows
// +build !windows

package app

import (
	"op-latency-mobile/internal/core"
	latencywin "op-latency-mobile/internal/latency_win"
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

// LatencyWindowsErrorEvent ...
type LatencyWindowsErrorEvent struct {
	Error string `json:"error"`
}

// StartWinOpLatency ...
func (a *Api) StartWinOpLatency(config latencywin.Config) (imageCount int, err error) {
	return
}

// WinOpLatencyResult ...
type WinOpLatencyResult struct {
	Latency       int `json:"latency"`
	ResponseIndex int `json:"responseIndex"`
	ResponseTime  int `json:"responseTime"`
}

// CalculateLatencyByImageDiff ...
func (a *Api) CalculateLatencyByImageDiff(imageRect core.ImageRectInfo, diffThreshold int) (result WinOpLatencyResult) {
	return
}

// CalculateLatencyByCurrentImage ...
func (a *Api) CalculateLatencyByCurrentImage(currenIndex int) (result WinOpLatencyResult) {
	return
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
	return GetImageResp{}
}

// OpenImageInExplorer ...
func (a *Api) OpenImageInExplorer(index int) {
}

// CaptureWindowInfo ...
type CaptureWindowInfo struct {
	Title  string `json:"title,omitempty"`
	Handle int    `json:"handle,omitempty"`
}

// ListCaptureWindows ...
func (a *Api) ListCaptureWindows() []CaptureWindowInfo {
	return nil
}
