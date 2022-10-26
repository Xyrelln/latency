package capture

import (
	"fmt"
	"image"
	_ "image/jpeg" // register format
	"os"
	"time"
)

// ScreenshotWithTs ...
type ScreenshotWithTs struct {
	FilePath string
	// Img      image.Image
	Time time.Time
}

// DecodeImg ...
func (s ScreenshotWithTs) DecodeImg() (image.Image, error) {
	imgf, err := os.Open(s.FilePath)
	if err != nil {
		return nil, err
	}
	defer imgf.Close()

	img, _, err := image.Decode(imgf)
	return img, err
}

// ScreenshotSeq ...
type ScreenshotSeq []ScreenshotWithTs

// FindImageHashResponseTime ...
func (seq ScreenshotSeq) FindImageHashResponseTime(diffArea image.Rectangle, diffThreshold int, timeAfter time.Time) (int, error) {
	if len(seq) <= 1 {
		return 0, fmt.Errorf("not enough screenshots")
	}

	var (
		curImg, prevImg image.Image
		err             error
	)

	for idx := 1; idx < len(seq); idx++ {
		// 忽略在操作时间之前的截图
		if seq[idx].Time.Before(timeAfter) {
			continue
		}

		curImg, err = seq[idx].DecodeImg()
		if err != nil {
			return 0, fmt.Errorf("decode %s error: %w", seq[idx].FilePath, err)
		}
		if prevImg == nil {
			prevImg, err = seq[idx-1].DecodeImg()
			if err != nil {
				return 0, fmt.Errorf("decode %s error: %w", seq[idx-1].FilePath, err)
			}
		}

		distance, _ := ImageHashDistance(prevImg, curImg, diffArea)
		// fmt.Printf("i1: %d, i2: %d, distance: %d\n", seq[idx-1].Time.UnixMilli(), seq[idx].Time.UnixMilli(), distance)
		if distance > diffThreshold {
			return idx, nil
		}
		prevImg = curImg
	}
	return 0, fmt.Errorf("find response time failed")
}
