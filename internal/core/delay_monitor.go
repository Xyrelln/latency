package core

import (
	"context"
	"fmt"
	"image"
	"op-latency-mobile/internal/ffprobe"
	"path"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	recordFile     = "rec.mp4"
	firstImageFile = "0001.png"
)

type DelayMonitor struct {
	VideoFolder         string          `json:"video_folder,omitempty"`
	ImagesFolder        string          `json:"images_folder,omitempty"`
	BlackWhiteThreshold int             `json:"black_white_threshold,omitempty"`
	PointerRect         image.Rectangle `json:"pointer_rect,omitempty"`
	PointerThreshold    float64         `json:"pointer_threshold,omitempty"`
	SceneRect           ImageRectInfo   `json:"scene_rect,omitempty"`
	SceneThreshold      int             `json:"scene_threshold,omitempty"`
}

func (dm *DelayMonitor) PTSPackets() (*ffprobe.PTSPackets, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	videoPath := path.Join(dm.VideoFolder, recordFile)
	data, err := ffprobe.ProbePTS(ctx, videoPath)
	if err != nil {
		log.Errorf("Error getting data: %v", err)
		return nil, err
	}
	return data, nil
}

func (dm *DelayMonitor) ProbeData(filePath string) (*ffprobe.ProbeData, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	data, err := ffprobe.ProbeURL(ctx, filePath)
	if err != nil {
		log.Errorf("Error getting data: %v", err)
		return nil, err
	}
	return data, nil
}

func (dm *DelayMonitor) FrameSpacingTime(ptsPackets *ffprobe.PTSPackets, startFrame, endFrame int) (*float64, error) {
	if endFrame < startFrame {
		return nil, fmt.Errorf("end frame must greater than start frame")
	}

	ptsSpaceing := ptsPackets.Packets[endFrame].PTS - ptsPackets.Packets[startFrame].PTS
	// timeBase := 1/90000
	t := float64(ptsSpaceing) * 1000.0 / 90000.0
	return &t, nil

}

func (dm *DelayMonitor) FrameSpacing() (startFrame, endFrame int, err error) {
	rect, err := GetCropRect(dm.SceneRect)
	if err != nil {
		log.Error("image with wrong scaling")
	}
	log.Printf("rect: %v", rect)
	imgs, err := ListImageFileWithCrop(dm.ImagesFolder, rect)
	if err != nil {
		log.Error("Specified directory with images inside does not exists or is corrupted")
	}

	var previousImg ImageFile
	var touched = false
	// var touchedIndex = 0
	// var spacing = 0

	/**
	1. 从 scoreT 评分判断触控操作
	2. 触控开始后检测中心内容区域评分
	3. 计算中间时间差
	**/
	for index, imageFile := range imgs {
		if index > 0 {
			if touched {
				diffCenter, _ := imageFile.ExtImgHashC.Distance(previousImg.ExtImgHashC)
				if diffCenter >= dm.SceneThreshold {
					log.Printf("find diffCenter: %d > threshold, index: %d", diffCenter, index)
					costTime := (float64(index - startFrame)) * (1000.0 / 60.0)
					// spacing = index - touchedIndex
					log.Printf("old cost time: %f", costTime)
					return startFrame, index, nil
					// return spacing, nil
				}
			} else {
				//diffTop, _ := imageFile.ExtImgHashT.Distance(previousImg.ExtImgHashT)
				_, diffTop, _ := CompareImages(imageFile.TouchAreaImg, previousImg.TouchAreaImg)
				if diffTop >= dm.PointerThreshold {
					log.Printf("find diffTop: %f > threshold, index: %d", diffTop, index)
					touched = true
					startFrame = index
				}
			}
		}
		previousImg = imageFile
	}
	return startFrame, endFrame, fmt.Errorf("failed to find start or end frame")
}

func NewDelayMonitor() *DelayMonitor {
	dm := DelayMonitor{}
	dm.BlackWhiteThreshold = 60
	dm.PointerThreshold = 6
	dm.SceneThreshold = 20
	return &dm
}

func (dm *DelayMonitor) Run() (*float64, error) {
	startFrame, endFrame, err := dm.FrameSpacing()
	if err != nil {
		return nil, fmt.Errorf("delay monitor run failed:%v", err)
	}

	ptsPackets, err := dm.PTSPackets()
	if err != nil {
		return nil, fmt.Errorf("delay monitor run failed:%v", err)
	}

	spacingTime, err := dm.FrameSpacingTime(ptsPackets, startFrame, endFrame)
	if err != nil {
		return nil, fmt.Errorf("delay monitor run failed:%v", err)
	}

	return spacingTime, nil
}
