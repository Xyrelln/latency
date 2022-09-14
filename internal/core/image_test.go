package core

import (
	"fmt"
	"image"
	"log"
	"os"
	"testing"

	"github.com/corona10/goimagehash"
)

func TestCropCurserArea(t *testing.T) {
	var imageRect ImageRectInfo
	fd, err := os.Open("/Users/jason/Downloads/mov/imgs/0066.png")
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

	// touchArea := image.Rect(0, 0, 100, 35)
	touchArea := image.Rect(0, 0, 50, 50)
	cropImg, _ := CropImage(img, touchArea)

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

	err = SaveImage("/Users/jason/Downloads/mov/crop/0002.png", cropImg)
	if err != nil {
		log.Fatal(err)
	}
}

func TestLoadImage(t *testing.T) {
	ListImageFile("/Users/jason/Developer/epc/op-latency-mobile/out/image/167-png/")

}

func TestAvageHash(t *testing.T) {
	fd, err := os.Open("/Users/jason/Downloads/mov/imgs/0064.png")
	if err != nil {
		log.Fatal(err)
		// return nil, err
	}
	defer fd.Close()
	img, _, _ := image.Decode(fd)
	touchArea := image.Rect(0, 0, 1920, 35)
	cropImg, _ := CropImage(img, touchArea)
	extImgHashT1, _ := goimagehash.ExtDifferenceHash(cropImg, 16, 16)

	fd2, err := os.Open("/Users/jason/Downloads/mov/imgs/0065.png")
	if err != nil {
		log.Fatal(err)
		// return nil, err
	}
	defer fd2.Close()
	img2, _, _ := image.Decode(fd2)
	cropImg2, _ := CropImage(img2, touchArea)
	extImgHashT2, _ := goimagehash.ExtDifferenceHash(cropImg2, 16, 16)
	score, _ := extImgHashT1.Distance(extImgHashT2)
	fmt.Printf("current score: %d", score)
	log.Printf("current score: %d", score)

}

func TestCalcTime(t *testing.T) {
	imagePath := "/Users/jason/Downloads/mov/imgs"
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
