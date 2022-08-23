package utils

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func CropImage(img image.Image, cropRect image.Rectangle) (cropImg image.Image, newImg bool) {
	//Interface for asserting whether `img`
	//implements SubImage or not.
	//This can be defined globally.
	type CropableImage interface {
		image.Image
		SubImage(r image.Rectangle) image.Image
	}

	if p, ok := img.(CropableImage); ok {
		// Call SubImage. This should be fast,
		// since SubImage (usually) shares underlying pixel.
		cropImg = p.SubImage(cropRect)
	} else if cropRect = cropRect.Intersect(img.Bounds()); !cropRect.Empty() {
		// If `img` does not implement `SubImage`,
		// copy (and silently convert) the image portion to RGBA image.
		rgbaImg := image.NewRGBA(cropRect)
		for y := cropRect.Min.Y; y < cropRect.Max.Y; y++ {
			for x := cropRect.Min.X; x < cropRect.Max.X; x++ {
				rgbaImg.Set(x, y, img.At(x, y))
			}
		}
		cropImg = rgbaImg
		newImg = true
	} else {
		// Return an empty RGBA image
		cropImg = &image.RGBA{}
		newImg = true
	}

	return cropImg, newImg
}

func GrayImage(img image.Image) (image.Image, error) {
	grayImg := image.NewGray(img.Bounds())
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			grayImg.Set(x, y, img.At(x, y))
		}
	}
	return grayImg, nil
}

func SaveImage(filePath string, src image.Image) error {
	f, err := os.OpenFile(filePath, os.O_SYNC|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	ext := filepath.Ext(filePath)
	if strings.EqualFold(ext, ".jpg") || strings.EqualFold(ext, ".jpeg") {
		err = jpeg.Encode(f, src, &jpeg.Options{Quality: 100})
	} else if strings.EqualFold(ext, ".png") {
		err = png.Encode(f, src)
	}
	return err
}
