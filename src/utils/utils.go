package utils

import (
	"errors"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

func GetTimeStamp() string {
	return time.Now().Format(time.RFC3339)
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
	workDir := path.Join(root, timestamp)
	videoDir := path.Join(workDir, "video")
	imagesDir := path.Join(workDir, "images")

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
