package core

import (
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path"
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

func GetFileNamesByPath(dirPath, extension string) ([]string, error) {
	var files []string
	fs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return files, err
	}

	for _, file := range fs {
		ext := filepath.Ext(file.Name())
		if strings.EqualFold(ext, extension) {
			files = append(files, file.Name())
		}
	}
	// sort it
	sort.Slice(files, func(i, j int) bool {
		return files[i] < files[j]
	})
	return files, nil
}

func LoadImage(path string, info os.FileInfo, e error) (image.Image, error) {
	if e != nil {
		return nil, e
	}
	if !info.IsDir() {
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

func ListImageFile(dirName string) error {
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
					extImgHashC, _ := goimagehash.ExtDifferenceHash(img, 16, 16)
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
	} else {
		// write duration of reading images and show images in a collage
	}
	log.Printf("image count: %d", len(imgs))
	return nil
}

// func ImageReader(imgPath string) (image.Image, error) {
// 	f, err := os.Open(path.Join("/Users/jason/Developer/epc/op-latency-mobile/out/image/167-png", imgPath))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer f.Close()

// 	ext := filepath.Ext(imgPath)
// 	var img image.Image
// 	if strings.EqualFold(ext, ".jpg") || strings.EqualFold(ext, ".jpeg") {
// 		img, err = jpeg.Decode(f)
// 	} else if strings.EqualFold(ext, ".png") {
// 		img, err = png.Decode(f)
// 	}
// 	return img, nil

// }
func CurseRead() {
	files, _ := GetFileNamesByPath("/Users/jason/Developer/epc/op-latency-mobile/out/image/167-png", ".png")
	// sort.Slice(files, func(i, j int) bool {
	// 	return files[i] < files[j]
	// })

	// var lastImg image.Image
	// var startCounter = false
	// var touchIndex = 0
	for _, fileName := range files {
		f, err := os.Open(path.Join("/Users/jason/Developer/epc/op-latency-mobile/out/image/167-png/", fileName))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		png.Decode(f)

		// var cropImg image.Image
		// if !startCounter {
		// 	touchArea := image.Rect(0, 0, 100, 35)
		// 	cropImg, _ = CropImage(img, touchArea)
		// } else {
		// 	width := img.Bounds().Dx()
		// 	height := img.Bounds().Dy()
		// 	centerRect := image.Rect(width/4, height/4, width/4*3, height/4*3)
		// 	cropImg, _ = CropImage(img, centerRect)
		// }

		// if lastImg != nil {
		// 	score, _ := ImageDiff(lastImg, cropImg)
		// 	if score > 20 {
		// 		// 检测到点击, 开始截取中心位置区域，同时开始统计
		// 		if startCounter == true {
		// 			diffIndex := index - touchIndex
		// 			costTime := diffIndex * (1000 / 60)
		// 			log.Printf("costTime: %d", costTime)
		// 			startCounter = false
		// 		} else {
		// 			touchIndex = index
		// 			startCounter = true
		// 			width := img.Bounds().Dx()
		// 			height := img.Bounds().Dy()
		// 			centerRect := image.Rect(width/4, height/4, width/4*3, height/4*3)
		// 			cropImg, _ = CropImage(img, centerRect)
		// 			log.Printf("file: %s", fileName)
		// 			log.Printf("score: %d", score)
		// 		}

		// 	}
		// }
		// lastImg = cropImg
	}
}
