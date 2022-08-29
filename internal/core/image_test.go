package core

import (
	"image"
	"image/png"
	"log"
	"os"
	"testing"
)

func TestCropCurserArea(t *testing.T) {
	fd, err := os.Open("/Users/jason/Developer/epc/op-latency-mobile/out/image/167-png/0001.png")
	if err != nil {
		log.Fatal(err)
		// return nil, err
	}
	img, _ := png.Decode(fd)
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	centerRect := image.Rect(width/4, height/4, width/4*3, height/4*3)

	cropImg, _ := CropImage(img, centerRect)
	// cropImage, err := CropCurserArea("/Users/jason/Developer/epc/op-latency-mobile/out/image/4/0001.png")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	grayImg := image.NewGray(cropImg.Bounds())
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			grayImg.Set(x, y, img.At(x, y))
		}
	}

	err = SaveImage("/Users/jason/Developer/epc/op-latency-mobile/out/image/crop/167-center-0001.png", grayImg)
	if err != nil {
		log.Fatal(err)
	}
}

func TestLoadImage(t *testing.T) {
	ListImageFile("/Users/jason/Developer/epc/op-latency-mobile/out/image/167-png/")

}

func TestCalcTime(t *testing.T) {
	CalcTime()
}
