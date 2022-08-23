package utils

import (
	"image"

	"github.com/corona10/goimagehash"
)

func ImageDiff(image1, image2 image.Image) (int, error) {
	// firstImageFile, err := os.Open(firstFile)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// secondImageFile, err := os.Open(secondiFile)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer firstImageFile.Close()
	// defer secondImageFile.Close()

	// firstImage, _ := jpeg.Decode(firstImageFile)
	firstHash, _ := goimagehash.ExtDifferenceHash(image1, 16, 16)

	// secondImage, _ := jpeg.Decode(secondImageFile)
	secondHash, _ := goimagehash.ExtDifferenceHash(image2, 16, 16)

	score, _ := secondHash.Distance(firstHash)
	// log.Printf("score: %d", score)
	return score, nil
}
