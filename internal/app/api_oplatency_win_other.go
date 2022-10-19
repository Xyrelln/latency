//go:build !windows
// +build !windows

package app

import (
	"op-latency-mobile/internal/core"
	latencywin "op-latency-mobile/internal/latency_win"
)

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
func (a *Api) CalculateLatencyByImageDiff(imageRect core.ImageRectInfo) (result WinOpLatencyResult) {
	return
}

// CalculateLatencyByCurrentImage ...
func (a *Api) CalculateLatencyByCurrentImage(currenIndex int) (result WinOpLatencyResult) {
	return
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
	return GetImageResp{}
}