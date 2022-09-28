package utils

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode"

	log "github.com/sirupsen/logrus"
)

func GetTimeStamp() string {
	return time.Now().Format("20060102150405.000")
}

func GetExecuteRoot() string {
	p, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Dir(p)
}

func CreateWorkDir() (string, string) {
	root := GetExecuteRoot()
	timestamp := GetTimeStamp()
	workDir := filepath.Join(root, "cache", timestamp)
	videoDir := filepath.Join(workDir, "video")
	imagesDir := filepath.Join(workDir, "images")

	if _, err := os.Stat(videoDir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(videoDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat(imagesDir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(imagesDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	return videoDir, imagesDir
}

func ClearCacheDir() {
	root := GetExecuteRoot()
	workDir := filepath.Join(root, "cache")
	go os.RemoveAll(workDir)
}

// isWindowsDrivePath returns true if the file path is of the form used by
// Windows. We check if the path begins with a drive letter, followed by a ":".
// For example: C:/x/y/z.
func IsWindowsDrivePath(path string) bool {
	if len(path) < 3 {
		return false
	}
	return unicode.IsLetter(rune(path[0])) && path[1] == ':'
}

func IsWindowsDrivePathURI(path string) bool {
	if len(path) < 3 {
		return false
	}
	return unicode.IsLetter(rune(path[1])) && path[2] == ':'
}

func GetImageFiles(pathname string, s []string) ([]string, error) {
	// var imgs []string
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return s, err
	}

	for _, fi := range rd {
		if !fi.IsDir() && strings.HasSuffix(fi.Name(), "png") {
			fullName := pathname + "/" + fi.Name()
			s = append(s, fullName)
		}
	}
	// sorted
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j] // filename as 0001   0002
	})

	return s, nil
}

func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
