package core

import (
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"testing"

	"github.com/corona10/goimagehash"
	// "github.com/otiai10/gosseract/v2"
)

func TestCropCurserArea(t *testing.T) {
	var imageRect ImageRectInfo
	fd, err := os.Open("/Users/jason/Developer/epc/op-latency-mobile/test_resource/touch_diff/0008.png")
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
	touchArea := image.Rect(0, 0, 100, 35)
	cropImg := CropImage(img, touchArea)

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

	err = SaveImage("/Users/jason/Downloads/mov/crop/0008.png", cropImg)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGrayImage(t *testing.T) {
	filename := "/Users/jason/Developer/epc/op-latency-mobile/build/bin/op-latency-mobile.app/Contents/MacOS/cache/mobile/20221101215156/images/0002.png"
	infile, err := os.Open(filename)
	defer infile.Close()
	if err != nil {
		log.Printf("failed opening %s: %s", filename, err)
		log.Fatal(err)
	}
	img, _, _ := image.Decode(infile)
	if err != nil {
		log.Fatal(err)
	}

	// touchArea := image.Rect(0, 0, 300, 300)
	// touchArea := image.Rect(0, 103, 100, 138)
	touchArea := image.Rect(0, 102, 137, 142)
	// cropImg, _ := CropImage(img, touchArea)

	result := image.NewRGBA(touchArea)
	draw.Draw(result, touchArea, img, touchArea.Min, draw.Src)

	// bwImg := RGBtoBlackAndWhite(result, 10)
	// bwImg := RGBtoGrayScale(result)

	bwImg := RGBtoGray(result)

	// func CropImage(img image.Image, area image.Rectangle) image.Image {
	// result := image.NewRGBA(touchArea)
	// draw.Draw(result, touchArea, img, touchArea.Min, draw.Src)
	// return result
	// }

	// Encode the grayscale image to the new file
	newFileName := "/Users/jason/Downloads/102.png"
	newfile, err := os.Create(newFileName)
	if err != nil {
		log.Printf("failed creating %s: %s", newfile.Name(), err)
		log.Fatal(err)

	}
	defer newfile.Close()
	png.Encode(newfile, result)

	newFileName2 := "/Users/jason/Downloads/102-bw.png"
	newfile2, err := os.Create(newFileName2)
	if err != nil {
		log.Printf("failed creating %s: %s", newfile2.Name(), err)
		log.Fatal(err)

	}
	defer newfile2.Close()
	png.Encode(newfile2, bwImg)
	// newfile.Close()
}

//	func TestLoadImage(t *testing.T) {
//		ListImageFile("/Users/jason/Developer/epc/op-latency-mobile/out/image/167-png/")
//
// }
func TestWhiteAndBlackDiff(t *testing.T) {
	f1 := "/Users/jason/Downloads/mov/0033-bw.png"
	f2 := "/Users/jason/Downloads/mov/0034-bw.png"
	img, percept, _ := CompareFiles(f1, f2)

	newFileName := "/Users/jason/Downloads/mov/diff-gray33-34.png"
	newfile, err := os.Create(newFileName)
	if err != nil {
		log.Printf("failed creating %s: %s", newfile.Name(), err)
		log.Fatal(err)
	}
	defer newfile.Close()
	png.Encode(newfile, img)
	log.Printf("diff percept %f", percept)
}

func TestAvageHash(t *testing.T) {
	fd, err := os.Open("/Users/jason/Downloads/mov/0033.png")
	if err != nil {
		log.Fatal(err)
		// return nil, err
	}
	defer fd.Close()
	img, _, _ := image.Decode(fd)
	touchArea := image.Rect(0, 0, 100, 35)
	//cropImg, _ := CropImage(img, touchArea)
	cropImg := img
	extImgHashT1, _ := goimagehash.ExtDifferenceHash(cropImg, 16, 16)

	fd2, err := os.Open("/Users/jason/Downloads/mov/0034.png")
	if err != nil {
		log.Fatal(err)
		// return nil, err
	}
	defer fd2.Close()
	img2, _, _ := image.Decode(fd2)
	cropImg2 := CropImage(img2, touchArea)
	extImgHashT2, _ := goimagehash.ExtDifferenceHash(cropImg2, 16, 16)
	score, _ := extImgHashT1.Distance(extImgHashT2)
	log.Printf("dHash 16* 16 score: %d \n", score)

	avgHash, _ := goimagehash.AverageHash(cropImg)
	avgHash2, _ := goimagehash.AverageHash(cropImg2)
	avgScore, _ := avgHash.Distance(avgHash2)
	log.Printf("avgScore: %d \n", avgScore)

	pHash, _ := goimagehash.PerceptionHash(cropImg)
	pHash2, _ := goimagehash.PerceptionHash(cropImg2)
	pScore, _ := pHash.Distance(pHash2)
	log.Printf("pScore: %d \n", pScore)

}

// func TestOcr(t *testing.T)  {
// 	client := gosseract.NewClient()
// 	defer client.Close()
// 	client.SetImage("/Users/jason/Downloads/mov/0034-bw.png")
// 	text, _ := client.Text()
// 	fmt.Println(text)
// }

//func TestCalcTime(t *testing.T) {
//	imagePath := "/Users/jason/Downloads/mov/6/images"
//	imgRect := ImageRectInfo{
//		X:             20,
//		Y:             26,
//		W:             446,
//		H:             70,
//		PreviewWidth:  500,
//		PreviewHeight: 281,
//		SourceWidth:   1920,
//		SourceHeight:  1080,
//	}
//
//	CalcTime(imagePath, imgRect, 20)
//}
