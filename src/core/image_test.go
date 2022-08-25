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

func TestLoadImage2(t *testing.T) {
	CurseRead()

}

// func TestCurseRead(t *testing.T) {
// 	files, _ := GetFilesByPath("/Users/jason/Developer/epc/op-latency-mobile/out/image/167-png", ".png")
// 	sort.Slice(files, func(i, j int) bool {
// 		return files[i] < files[j]
// 	})

// 	var lastImg image.Image
// 	var startCounter = false
// 	var touchIndex = 0
// 	for index, fileName := range files {
// 		f, err := os.Open(path.Join("../../out/image/167-png", fileName))
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer f.Close()
// 		img, _ := png.Decode(f)

// 		var cropImg image.Image
// 		if !startCounter {
// 			touchArea := image.Rect(0, 0, 100, 35)
// 			cropImg, _ = CropImage(img, touchArea)
// 		} else {
// 			width := img.Bounds().Dx()
// 			height := img.Bounds().Dy()
// 			centerRect := image.Rect(width/4, height/4, width/4*3, height/4*3)
// 			cropImg, _ = CropImage(img, centerRect)
// 		}

// 		if lastImg != nil {
// 			score, _ := ImageDiff(lastImg, cropImg)
// 			if score > 20 {
// 				// 检测到点击, 开始截取中心位置区域，同时开始统计
// 				if startCounter == true {
// 					diffIndex := index - touchIndex
// 					costTime := diffIndex * (1000 / 60)
// 					log.Printf("costTime: %d", costTime)
// 					startCounter = false
// 				} else {
// 					touchIndex = index
// 					startCounter = true
// 					width := img.Bounds().Dx()
// 					height := img.Bounds().Dy()
// 					centerRect := image.Rect(width/4, height/4, width/4*3, height/4*3)
// 					cropImg, _ = CropImage(img, centerRect)
// 					log.Printf("file: %s", fileName)
// 					log.Printf("score: %d", score)
// 				}

// 			}
// 		}
// 		lastImg = cropImg
// 	}

// }
