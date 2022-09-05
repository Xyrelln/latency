package core

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/corona10/goimagehash"
	"golang.org/x/sync/errgroup"
)

var (
	threshold_t = 20
	threshold_c = 20
)

type ImageFile struct {
	Path string
	Img  image.Image
	// ExtImgHash  *goimagehash.ExtImageHash
	ExtImgHashT *goimagehash.ExtImageHash // touch area
	ExtImgHashC *goimagehash.ExtImageHash // center area
}

type ImageRectInfo struct {
	X             int `json:"x"`
	Y             int `json:"y"`
	W             int `json:"w"`
	H             int `json:"h"`
	PreviewWidth  int `json:"preview_width"`
	PreviewHeight int `json:"preview_height"`
	SourceWidth   int `json:"source_width"`
	SourceHeight  int `json:"source_height"`
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

func GetCropRect(imageRect ImageRectInfo) (image.Rectangle, error) {
	// 0 0 200 200 600 338
	// 111 111 200 200 600 338
	// proportion := imageRect.SourceWidth / imageRect.PreviewWidth
	// yProportion := imageRect.SourceHeight / imageRect.PreviewHeight

	// if (proportion - yProportion) > 0 {
	// 	return image.Rect(0, 0, 0, 0), errors.New("image with wrong scaling")
	// }

	x0 := imageRect.X * imageRect.SourceWidth / imageRect.PreviewWidth
	y0 := imageRect.Y * imageRect.SourceHeight / imageRect.PreviewHeight
	x1 := imageRect.W*imageRect.SourceWidth/imageRect.PreviewWidth + x0
	y1 := imageRect.H*imageRect.SourceHeight/imageRect.PreviewHeight + y0
	cropRect := image.Rect(x0, y0, x1, y1)

	return cropRect, nil
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
					centerArea := image.Rect(width/4, height/4, width/4*3, height/4*3)
					cropImgC, _ := CropImage(img, centerArea)
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

func ListImageFileWithCrop(dirName string, rect image.Rectangle) ([]ImageFile, error) {
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

					// width := img.Bounds().Dx() // @todo get x,y by phone
					// height := img.Bounds().Dy()
					// centerArea := image.Rect(width/4, height/4, width/4*3, height/4*3)
					cropImgC, _ := CropImage(img, rect)
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

func GetImageInfo(imagePath string) (ImageInfo, error) {
	log.Printf("image path: %s", imagePath)
	var imgInfo ImageInfo
	fimg, err := os.Open(imagePath)
	if err != nil {
		log.Fatal(err)
	}
	defer fimg.Close()
	img, _, err := image.Decode(fimg)
	if err != nil {
		log.Fatal(err)
		return imgInfo, err
	}
	imgInfo.Path = imagePath
	imgInfo.Width = img.Bounds().Dx()
	imgInfo.Height = img.Bounds().Dy()
	return imgInfo, nil

}

func CalcTime(imgPath string, imageRect ImageRectInfo) ([]int, error) {
	// dir := "/Users/jason/Developer/epc/op-latency-mobile/out/image/167-png/"

	rect, err := GetCropRect(imageRect)
	if err != nil {
		log.Fatal("image with wrong scaling")
	}
	log.Printf("rect: %v", rect)
	imgs, err := ListImageFileWithCrop(imgPath, rect)
	if err != nil {
		log.Fatal("Specified directory with images inside does not exists or is corrupted")
	}

	var previousImg ImageFile
	var touched = false
	var touchedIndex = 0
	// var touchIndex = 0
	// 1. 从 touch 展示区域评分判断触控操作
	/**
	1. 从 scoreT 评分判断触控操作
	2. 触控开始后检测中心内容区域评分
	3. 计算中间时间差
	**/
	var responseTime []int
	for index, imageFile := range imgs {
		if index > 0 {
			// start diff in second item
			scoreT, _ := imageFile.ExtImgHashT.Distance(previousImg.ExtImgHashT)
			scoreC, _ := imageFile.ExtImgHashC.Distance(previousImg.ExtImgHashC)
			log.Printf("scoreT: %d", scoreT)
			log.Printf("scoreC: %d", scoreC)

			if touched {
				if scoreC >= threshold_c {
					costTime := (index - touchedIndex) * (1000 / 60)
					responseTime = append(responseTime, costTime)

					// reset touched
					log.Printf("current file: %s", imageFile.Path)
					touched = false
				}

			} else {
				if scoreT >= threshold_t {
					touched = true
					touchedIndex = index
				}
			}

		}
		previousImg = imageFile
	}
	log.Printf("responseTime: %v", responseTime)
	return responseTime, nil
}