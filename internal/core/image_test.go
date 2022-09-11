package core

import (
	"image"
	"log"
	"os"
	"testing"
)

func TestCropCurserArea(t *testing.T) {
	var imageRect ImageRectInfo
	fd, err := os.Open("/Users/jason/Developer/epc/op-latency-mobile/build/bin/op-latency-mobile.app/Contents/MacOS/2022-09-01T15:33:13+08:00/images/0002.png")
	defer fd.Close()
	if err != nil {
		log.Fatal(err)
		// return nil, err
	}
	img, _, _ := image.Decode(fd)
	// width := img.Bounds().Dx()
	// height := img.Bounds().Dy()
	// centerRect := image.Rect(117, 500, 535, 535)
	imageRect.X = 44
	imageRect.Y = 187
	imageRect.W = 200
	imageRect.PreviewWidth = 717
	imageRect.PreviewHeight = 403
	imageRect.SourceHeight = 1920
	imageRect.SourceWidth = 1080

	rect, _ := GetCropRect(imageRect)
	log.Print(rect)
	cropImg, _ := CropImage(img, rect)

	// cropImage, err := CropCurserArea("/Users/jason/Developer/epc/op-latency-mobile/out/image/4/0001.png")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// grayImg := image.NewGray(cropImg.Bounds())
	// for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
	// 	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
	// 		grayImg.Set(x, y, img.At(x, y))
	// 	}
	// }

	err = SaveImage("/Users/jason/Developer/epc/op-latency-mobile/out/image/crop/1537-0001.png", cropImg)
	if err != nil {
		log.Fatal(err)
	}
}

func TestLoadImage(t *testing.T) {
	ListImageFile("/Users/jason/Developer/epc/op-latency-mobile/out/image/167-png/")

}

func TestCalcTime(t *testing.T) {
	imagePath := "/Users/jason/Developer/epc/op-latency-mobile/build/bin/op-latency-mobile.app/Contents/MacOS/cache/20220911175724.387/images"
	imgRect := ImageRectInfo{
		X:             20,
		Y:             26,
		W:             446,
		H:             70,
		PreviewWidth:  500,
		PreviewHeight: 281,
		SourceWidth:   1920,
		SourceHeight:  1080,
	}

	CalcTime(imagePath, imgRect, 20)
}
