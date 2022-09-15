package core

import (
    "image"
    "image/color"
    "math"
)



// RGBtoBlackAndWhite 黑白处理
func RGBtoBlackAndWhite(originalImage image.Image, whiteThreshold uint8) image.Image {
    originalImage = RGBtoGrayScale(originalImage)

    size := originalImage.Bounds().Size()
    rect := image.Rect(0, 0, size.X, size.Y)

    modifiedImg := image.NewRGBA(rect)

    for x := 0; x < size.X; x++ {
        for y := 0; y < size.Y; y++ {
            pixel := originalImage.At(x, y)
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

            modifiedImg.Set(x, y, modifiedColor)
        }
    }

    return modifiedImg
}

// RGBtoGrayScale 灰阶处理
func RGBtoGrayScale(originalImage image.Image) image.Image {
    size := originalImage.Bounds().Size()
    rect := image.Rect(0,0, size.X, size.Y)

    modifiedImg := image.NewNRGBA(rect)

    for x := 0; x < size.X; x++ {
        for y := 0; y < size.Y; y++ {
            pixel := originalImage.At(x, y)
            originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)

            red := float64(originalColor.R)
            green := float64(originalColor.G)
            blue := float64(originalColor.B)

            grey := uint8(
                math.Round((red+green+blue) / 3),
            )

            modifiedColor := color.RGBA{
                R: grey,
                G: grey,
                B: grey,
                A: originalColor.A,
            }

            modifiedImg.Set(x, y, modifiedColor)
        }
    }

    return modifiedImg
}
