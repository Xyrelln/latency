package core

import (
	"image"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/corona10/goimagehash"
	"golang.org/x/sync/errgroup"
)

type ImageFile struct {
	Path string
	Img  image.Image
	// ExtImgHash  *goimagehash.ExtImageHash
	ExtImgHashT *goimagehash.ExtImageHash // touch area
	ExtImgHashC *goimagehash.ExtImageHash // center area
}

// func GetFileNamesByPath(dirPath, extension string) ([]string, error) {
// 	var files []string
// 	fs, err := ioutil.ReadDir(dirPath)
// 	if err != nil {
// 		return files, err
// 	}

// 	for _, file := range fs {
// 		ext := filepath.Ext(file.Name())
// 		if strings.EqualFold(ext, extension) {
// 			files = append(files, file.Name())
// 		}
// 	}
// 	// sort it
// 	sort.Slice(files, func(i, j int) bool {
// 		return files[i] < files[j]
// 	})
// 	return files, nil
// }

func IsImage(info os.FileInfo) bool {
	var pngExt = ".png"
	var jpegExt = ".jpeg"
	ext := filepath.Ext(info.Name())
	return strings.EqualFold(ext, pngExt) || strings.EqualFold(ext, jpegExt)
}

func LoadImage(path string, info os.FileInfo, e error) (image.Image, error) {
	if e != nil {
		return nil, e
	}
	if !info.IsDir() && IsImage(info) {
		fimg, _ := os.Open(path)
		defer fimg.Close()
		img, _, imageError := image.Decode(fimg)
		if imageError == nil {
			return img, nil
		} else {
			return nil, imageError
		}
	}
	return nil, nil
}

func ListImageFile(dirName string) ([]ImageFile, error) {
	// var images []image.Image
	var imgs []ImageFile
	var eg errgroup.Group
	touchArea := image.Rect(0, 0, 100, 35)

	_ = filepath.Walk(
		dirName,
		func(path string, info os.FileInfo, e error) error {
			eg.Go(func() error {
				img, err := LoadImage(path, info, e)
				if img != nil {
					// log.Println(path)
					// images = append(images, img)
					cropImgT, _ := CropImage(img, touchArea)
					extImgHashT, _ := goimagehash.ExtDifferenceHash(cropImgT, 16, 16)

					width := img.Bounds().Dx() // @todo get x,y by phone
					height := img.Bounds().Dy()
					centerRect := image.Rect(width/4, height/4, width/4*3, height/4*3)
					cropImgC, _ := CropImage(centerRect, touchArea)
					extImgHashC, _ := goimagehash.ExtDifferenceHash(cropImgC, 16, 16)
					imgs = append(imgs, ImageFile{
						Path:        path,
						Img:         img,
						ExtImgHashC: extImgHashC,
						ExtImgHashT: extImgHashT,
					})
				}
				return err
			})
			return nil
		},
	)
	err := eg.Wait()
	if err != nil {
		log.Fatal("Specified directory with images inside does not exists or is corrupted")
	}
	// sorted
	sort.Slice(imgs, func(i, j int) bool {
		return imgs[i].Path < imgs[j].Path // filename as 0001   0002
	})

	log.Printf("image count: %d", len(imgs))
	return imgs, nil
}

func CalcTime() {
	dir := "/Users/jason/Developer/epc/op-latency-mobile/out/image/167-png/"
	imgs, err := ListImageFile(dir)
	if err != nil {
		log.Fatal("Specified directory with images inside does not exists or is corrupted")
	}

	var previousImg ImageFile
	// var touchIndex = 0
	for index, imageFile := range imgs {
		if index > 0 {
			scoreT, _ := imageFile.ExtImgHashT.Distance(previousImg.ExtImgHashT)
			scoreC, _ := imageFile.ExtImgHashC.Distance(previousImg.ExtImgHashC)
			log.Printf("scoreT: %d", scoreT)
			log.Printf("scoreC: %d", scoreC)
		}
		// if scoreT > 20 {
		// 	// 检测到点击, 开始截取中心位置区域，同时开始统计
		// 	if startCounter == true {
		// 		diffIndex := index - touchIndex
		// 		costTime := diffIndex * (1000 / 60)
		// 		log.Printf("costTime: %d", costTime)
		// 		startCounter = false
		// 	} else {
		// 		touchIndex = index
		// 		startCounter = true
		// 		log.Printf("file: %s", fileName)
		// 		log.Printf("score: %d", score)
		// 	}
		// }
		previousImg = imageFile
	}
}
