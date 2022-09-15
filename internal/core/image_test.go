package core

import (
	"image"
	"image/color"
	"image/png"
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

func TestGrayImage(t *testing.T) {
	filename := "/Users/jason/Developer/epc/op-latency-mobile/test_resource/touch_diff/0002.png"
	infile, err := os.Open(filename)

	if err != nil {
		log.Printf("failed opening %s: %s", filename, err)
		log.Fatal(err)
	}
	defer infile.Close()

	imgSrc, _, err := image.Decode(infile)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new grayscale image
	const whiteThreshold = uint8(256 / 2) // 128
	bounds := imgSrc.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	grayScale := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			//imageColor := imgSrc.At(x, y)
			//rr, gg, bb, _ := imageColor.RGBA()
			//r := math.Pow(float64(rr), 2.2)
			//g := math.Pow(float64(gg), 2.2)
			//b := math.Pow(float64(bb), 2.2)
			//m := math.Pow(0.2125*r+0.7154*g+0.0721*b, 1/2.2)
			//Y := uint16(m + 0.5)
			//grayColor := color.Gray{uint8(Y >> 8)}

			pixel := imgSrc.At(x, y)
			originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)
			modifiedColorValue := originalColor.R
			if modifiedColorValue >= whiteThreshold {
				modifiedColorValue = 255
			} else {
				modifiedColorValue = 0
			}

			modifiedColor := color.RGBA{
				R: modifiedColorValue,
				G: modifiedColorValue,
				B: modifiedColorValue,
				A: originalColor.A,
			}

			grayScale.Set(x, y, modifiedColor)

			//oldPixel := imgSrc.At(x, y)
			//r, g, b, _ := oldPixel.RGBA()
			//lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			//grayColor := color.Gray{uint8(lum / 256)}
			//
			//grayScale.Set(x, y, grayColor)
		}
	}

	// Encode the grayscale image to the new file
	newFileName := "/Users/jason/Developer/epc/op-latency-mobile/test_resource/touch_diff/0001-gray4.png"
	newfile, err := os.Create(newFileName)
	if err != nil {
		log.Printf("failed creating %s: %s", newfile.Name(), err)
		log.Fatal(err)


	}
	defer newfile.Close()
	png.Encode(newfile,grayScale)
}

//func TestLoadImage(t *testing.T) {
//	ListImageFile("/Users/jason/Developer/epc/op-latency-mobile/out/image/167-png/")
//
//}
func TestWhiteAndBlackDiff(t *testing.T){
	f1 := "/Users/jason/Developer/epc/op-latency-mobile/test_resource/touch_diff/0001-gray3.png"
	f2 := "/Users/jason/Developer/epc/op-latency-mobile/test_resource/touch_diff/0001-gray4.png"
	img, percept, _ := CompareFiles(f1, f2)

	newFileName := "/Users/jason/Developer/epc/op-latency-mobile/test_resource/touch_diff/diff-gray3-4.png"
	newfile, err := os.Create(newFileName)
	if err != nil {
		log.Printf("failed creating %s: %s", newfile.Name(), err)
		log.Fatal(err)
	}
	defer newfile.Close()
	png.Encode(newfile,img)
	log.Printf("diff percept %f", percept)
}

func TestAvageHash(t *testing.T) {
	fd, err := os.Open("/Users/jason/Developer/epc/op-latency-mobile/test_resource/touch_diff/0001-gray3.png")
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

	fd2, err := os.Open("/Users/jason/Developer/epc/op-latency-mobile/test_resource/touch_diff/0001-gray4.png")
	if err != nil {
		log.Fatal(err)
		// return nil, err
	}
	defer fd2.Close()
	img2, _, _ := image.Decode(fd2)
	cropImg2, _ := CropImage(img2, touchArea)
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

func TestCalcTime(t *testing.T) {
	imagePath := "/Users/jason/Downloads/mov/5"
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
