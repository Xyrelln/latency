package upload

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gitlab.vrviu.com/epc/lighttest-lib/lighttestservice"
	"gitlab.vrviu.com/epc/lighttest-lib/token"
)

func UploadFile(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		log.Errorf("open file failed, err: %v", err)
		return err
	}
	defer f.Close()

	svc := lighttestservice.LightTestService{Endpoint: "https://lighttest.vrviu.com"}
	client := token.ClientInfo{
		Name:     "perftool",
		Version:  "0.0.1",
		Username: "NarakaPlayer",
	}
	log.Infof("prepare upload file: %s, base name: %s", filePath, filepath.Base(filePath))
	uploadPath, err := svc.UploadFile(client, filepath.Base(filePath), f)
	if err != nil {
		log.Errorf("upload file failed, err: %v", err)
		return err
	}
	log.Infof("upload success, path: %s", uploadPath)
	return nil
}

type Result struct {
	CostTime float64 `json:"cost_time"`
}

func UploadResult(costTime float64) error {
	svc := lighttestservice.LightTestService{Endpoint: "http://10.86.3.236:8088"}
	client := token.ClientInfo{
		Name:     "perftool",
		Version:  "0.0.1",
		Username: "NarakaPlayer",
	}
	err := svc.UploadJSONData(client, Result{CostTime: costTime})
	if err != nil {
		log.Errorf("upload result failed, err: %v", err)
	}

	return nil
}
