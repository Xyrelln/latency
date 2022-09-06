package utils

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"
	"unicode"
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
