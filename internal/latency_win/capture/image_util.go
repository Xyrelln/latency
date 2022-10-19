package capture

import (
	// "bytes"
	"image"
	"image/draw"

	"github.com/corona10/goimagehash"
	// "github.com/lucasb-eyer/go-colorful"
)

// CheckAreaRed ...
// func CheckAreaRed(data []byte, area image.Rectangle) (bool, error) {
// 	img, _, err := image.Decode(bytes.NewReader(data))
// 	if err != nil {
// 		return false, err
// 	}

// 	var c colorful.Color
// 	for i := area.Min.X; i < area.Max.X; i++ {
// 		for j := area.Min.Y; j < area.Max.Y; j++ {
// 			c, _ = colorful.MakeColor(img.At(i, j))

// 			if h, s, v := c.Hsv(); (h >= 330 || h <= 30) && s >= 0.3 && v >= 0.3 {
// 				return true, nil
// 			}
// 		}
// 	}
// 	return false, nil
// }

// ImageHashDistance ...
func ImageHashDistance(img1, img2 image.Image, area image.Rectangle) (int, error) {
	img1 = CropImage(img1, area)
	img2 = CropImage(img2, area)

	hash1, err := goimagehash.ExtDifferenceHash(img1, 16, 16)
	hash2, err := goimagehash.ExtDifferenceHash(img2, 16, 16)
	// hash1, err := goimagehash.AverageHash(img1)
	// hash2, err := goimagehash.AverageHash(img2)
	if err != nil {
		return 0, err
	}
	return hash1.Distance(hash2)
}

// CropImage ...
func CropImage(img image.Image, area image.Rectangle) image.Image {
	result := image.NewRGBA(area)
	draw.Draw(result, area, img, area.Min, draw.Src)
	return result
}
