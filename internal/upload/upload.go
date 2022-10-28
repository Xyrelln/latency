package upload

import (
	"os"
	"path/filepath"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.vrviu.com/epc/lighttest-lib/lighttestservice"
	"gitlab.vrviu.com/epc/lighttest-lib/token"
)

type lightService struct {
	Service lighttestservice.LightTestService
	Client  token.ClientInfo
}

var lightServiceIns *lightService
var once sync.Once

const secretKey = "EBIpWLd1cC29gQl6"

// GetLightService 返回 lightService 实例
func GetLightService() *lightService {
	once.Do(func() {
		lightServiceIns = &lightService{
			Service: lighttestservice.LightTestService{Endpoint: "https://lighttest.vrviu.com"},
			Client: token.ClientInfo{
				Name:     "op-latency",
				Version:  "0.2.1",
				Username: "vrviu",
			},
		}
	})
	return lightServiceIns
}
func GetUserInfo() (lighttestservice.UserInfo, error) {
	ls := GetLightService()
	userInfo, err := ls.Service.GetUserInfo(token.ClientInfo{
		Name:     "op-latency",
		Version:  "0.2.1",
		Username: "vrviu",
	}, secretKey)

	if err != nil {
		return lighttestservice.UserInfo{}, err
	}
	return userInfo, nil
}

// UploadFile 上传文件
func UploadFile(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		log.Errorf("open file failed, err: %v", err)
		return err
	}
	defer f.Close()

	ls := GetLightService()

	// svc := lighttestservice.LightTestService{Endpoint: "https://lighttest.vrviu.com"}
	// client := token.ClientInfo{
	// 	Name:     "perftool",
	// 	Version:  "0.0.1",
	// 	Username: "NarakaPlayer",
	// }
	log.Infof("prepare upload file: %s, base name: %s", filePath, filepath.Base(filePath))
	uploadPath, err := ls.Service.UploadFile(ls.Client, secretKey, filepath.Base(filePath), f)
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

// UploadResult 上传结果
func UploadResult(costTime float64) error {
	ls := GetLightService()
	// svc := lighttestservice.LightTestService{Endpoint: "http://10.86.3.236:8088"}
	// client := token.ClientInfo{
	// 	Name:     "perftool",
	// 	Version:  "0.0.1",
	// 	Username: "NarakaPlayer",
	// }
	err := ls.Service.UploadJSONData(ls.Client, secretKey, Result{CostTime: costTime})
	if err != nil {
		log.Errorf("upload result failed, err: %v", err)
	}

	return nil
}
