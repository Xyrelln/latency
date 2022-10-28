package app

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"image"
	"path/filepath"
	"strings"
	"time"

	"op-latency-mobile/internal/adb"
	"op-latency-mobile/internal/cmd"
	"op-latency-mobile/internal/core"
	"op-latency-mobile/internal/ffprobe"
	"op-latency-mobile/internal/fs"
	"op-latency-mobile/internal/upload"

	latencywin "op-latency-mobile/internal/latency_win"

	log "github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	// lighttestServ "gitlab.vrviu.com/epc/lighttest-lib/lighttestservice"
	// lighttestToken "gitlab.vrviu.com/epc/lighttest-lib/token"
	"gitlab.vrviu.com/epc/lighttest-lib/lighttestservice"
	lighttestUpdate "gitlab.vrviu.com/epc/lighttest-lib/update"
	lighttestVer "gitlab.vrviu.com/epc/lighttest-lib/version"
)

const (
	appName                  = "latency"
	localPassFile            = "latency-pass"
	lighttestServiceEndpoint = "https://lighttest.vrviu.com"
	defaultStateKey          = "state_default"
	defaultWorkspaceKey      = "wksp_default"
	recordFile               = "rec.mp4"
	firstImageFile           = "0001.png"
	sceneKeyPrefix           = "scene_"
)

var autorun = true

// Api struct
type Api struct {
	ctx       context.Context `json:"ctx,omitempty"`
	Cmd       *cmd.Cmd        `json:"cmd,omitempty"`
	VideoDir  string          `json:"video_dir,omitempty"`
	ImagesDir string          `json:"images_dir,omitempty"`
	AppData   string          `json:"app_data,omitempty"`
	state     *workspaceState
	logger    *logger.Logger
	store     *store
	cmdRunner *cmd.CmdRunner

	// windows op latency
	latencyWinManager *latencywin.OpLatencyWindowsManager
}

type Record struct {
	VideoDir  string `json:"video_dir,omitempty"`
	ImagesDir string `json:"images_dir,omitempty"`
}

type storeLogger struct {
	logger.Logger
}

func (s storeLogger) Warningf(message string, args ...interface{}) {
	s.Warningf(message, args...)
}

// NewApp creates a new Api application struct
func NewApp() *Api {
	return &Api{
		latencyWinManager: latencywin.NewOpLatencyWindowsManager(),
	}
}

// startup is called when the Api starts. The context is saved
// so we can call the runtime methods
func (a *Api) startup(ctx context.Context) {
	a.ctx = ctx

	store, err := newStore(a.AppData)
	a.store = store
	if err != nil {
		log.Errorf("app: failed to create database: %v", err)
		// fmt.Errorf("app: failed to create database: %v", err)
	} else {
		a.state = a.getCurrentState()
	}
}

func (a *Api) shutdown(ctx context.Context) {
	a.store.close()
}

func (a *Api) domready(ctx context.Context) {
	a.IsAppReady()
}

func (a *Api) getCurrentState() *workspaceState {
	rtn := &workspaceState{
		CurrentID: defaultWorkspaceKey,
	}
	val, err := a.store.get([]byte(defaultStateKey))
	if err != nil && err != errKeyNotFound {
		log.Errorf("failed to get current state from store: %v", err)
	}
	if len(val) == 0 {
		return rtn
	}
	dec := gob.NewDecoder(bytes.NewBuffer(val))
	if err := dec.Decode(rtn); err != nil {
		log.Errorf("failed to decode state: %v", err)
	}
	return rtn
}

// 获取设备列表
func (a *Api) ListDevices() ([]*adb.Device, error) {
	devices, err := adb.Devices()
	if err != nil {
		log.Errorf("ListDevices failed: %v", err)
		return nil, err
	}
	return devices, nil
}

type UserAction struct {
	Auto  bool   `json:"auto"`
	Type  string `json:"type"`
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Tx    int    `json:"tx"`
	Ty    int    `json:"ty"`
	Speed int    `json:"speed"`
}

type CropInfo struct {
	Top    int `json:"top"`
	Left   int `json:"left"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type DeviceInfo struct {
	DeviceName   string `json:"device_name"`
	ScreenWidth  int    `json:"screen_width"`
	ScreenHeight int    `json:"screen_height"`
}

type UserScene struct {
	Name           string     `json:"name"`
	Key            string     `json:"key"`
	Device         DeviceInfo `json:"device"`
	CropCoordinate CropInfo   `json:"crop_coordinate"`
	Action         UserAction `json:"action"`
}

// 获取设备列表
func (a *Api) ListScens() ([]UserScene, error) {
	log.Infof("list scenes")
	var userScenes []UserScene

	items, err := a.store.list([]byte(sceneKeyPrefix))
	if err != nil {
		log.Errorf("failed to get scenes from store: %v", err)
		return nil, err
	}

	for _, val := range items {
		scene := UserScene{}
		dec := gob.NewDecoder(bytes.NewBuffer(val))
		if err = dec.Decode(&scene); err != nil {
			return userScenes, err
		}
		userScenes = append(userScenes, scene)
	}

	log.Infof("scenes content: %v", userScenes)
	return userScenes, nil
}

func (a *Api) SetScene(userScene UserScene) {
	log.Infof("set scene content: %v", userScene)
	var val bytes.Buffer

	name := fs.GetRandString(3)
	userScene.Key = name

	enc := gob.NewEncoder(&val)
	enc.Encode(userScene)

	a.store.set([]byte(sceneKeyPrefix+name), val.Bytes())
}

func (a *Api) DeleteScene(key string) {
	log.Infof("delete scene name: %s", key)
	if key == "" {
		return
		// items, err := a.store.list([]byte(sceneKeyPrefix))
		// if err != nil {
		// 	log.Errorf("failed to get scenes from store: %v", err)
		// 	// return nil, err
	}
	a.store.del([]byte(sceneKeyPrefix + key))
}

// 检查 app 依赖包环境信息
func (a *Api) IsAppReady() error {
	p, err := adb.IsAdbReady()
	if err != nil {
		log.Error("adb path wrong")
		return err
	}
	log.Infof("adb path: %s", p)

	p, err = cmd.IsFFmpegReady()
	if err != nil {
		log.Error("ffmpeg path wrong")
		return err
	}
	log.Infof("ffmpeg path: %s", p)

	p, err = ffprobe.IsFFprobeReady()
	if err != nil {
		log.Error("ffprobe path wrong")
		return err
	}
	log.Infof("ffprobe path: %s", p)

	p, err = cmd.IsScrcpyReady()
	if err != nil {
		log.Error("scrcpy path wrong")
		return err
	}
	log.Infof("scrcpy path: %s", p)
	return nil
}

// 获取屏幕信息
func (a *Api) GetDisplay(serial string) (*adb.Display, error) {
	device := adb.GetDevice(serial)
	display, err := device.DisplaySize()
	if err != nil {
		return nil, err
	}
	return display, nil
}

func (a *Api) GetPhysicalSize(serial string) (*adb.Display, error) {
	device := adb.GetDevice(serial)
	display, err := device.PhysicalSize()
	if err != nil {
		return nil, err
	}
	return display, nil
}

func (a *Api) SetAutoSwipeOn(sw adb.SwipeEvent, interval int) error {
	autorun = true
	devices, err := adb.Devices()
	if err != nil {
		log.Errorf("ListDevices failed: %v", err)
		return err
	}
	for {
		for _, device := range devices {
			go device.AutoSwipe(sw)
		}
		time.Sleep(time.Duration(interval) * time.Millisecond)
		if !autorun {
			break
		}
	}

	return nil
}

// 加载截图
func (a *Api) LoadScreenshot(serial string) (core.ImageInfo, error) {
	var img core.ImageInfo
	device := adb.GetDevice(serial)
	imgPath, err := device.GetScreenshot()
	if err != nil {
		return img, err
	}
	screenshotDir, err := fs.GetScreenshotDir()
	if err != nil {
		log.Errorf("get screenshot dir err %v", err)
		return img, err
	}
	localPath := filepath.Join(screenshotDir, "screenshot.png")
	err = device.Pull(imgPath, localPath)
	if err != nil {
		return img, err
	}

	mInfo, err := core.GetImageInfo(localPath)
	if err != nil {
		return mInfo, err
	}
	// path to uri format for FileLoader on win
	if fs.IsWindowsDrivePath(mInfo.Path) {
		mInfo.Path = "/" + strings.ReplaceAll(mInfo.Path, "\\", "/")
	}
	return mInfo, nil
}

func (a *Api) SetAutoSwipeOff() error {
	autorun = false
	return nil
}

// IsPointerLocationOn 查询指针开启状态
func (a *Api) IsPointerLocationOn(serial string) (bool, error) {
	device := adb.GetDevice(serial)
	on, err := device.IsPointerLocationOn()
	if err != nil {
		return false, nil
	}
	return on, nil
}

// 开启指针位置显示
func (a *Api) SetPointerLocationOn(serial string) error {
	device := adb.GetDevice(serial)
	if err := device.SetPointerLocationOn(); err == nil {
		on, err := device.IsPointerLocationOn()
		if err != nil {
			return err
		}
		if !on {
			return errors.New("set pointer location failed")
		}
	} else {
		return err
	}

	return nil
}

// 关闭指针位置显示
func (a *Api) SetPointerLocationOff(serial string) error {
	log.Info("set pointer location off")
	device := adb.GetDevice(serial)
	return device.SetPointerLocationOff()
}

func (a *Api) StartRecord(serial string, userAction UserAction) (rerr error) {
	defer func() {
		if rerr != nil {
			log.Error(rerr)
			a.emitInfo(eventRecordStartError)
		}
	}()

	a.VideoDir, a.ImagesDir = fs.CreateWorkDir()
	recFile := filepath.Join(a.VideoDir, recordFile)
	log.Infof("starting record, store path: %s", recFile)

	runner, rerr := cmd.ScrcpyStart(serial, recFile)
	a.cmdRunner = runner

	// record file exists check
	go func() {
		var nums [15][0]int
		interval := 200

		for range nums {
			isExists := fs.FileSizeGreaterThan(recFile, 1024*1024)
			if isExists {
				log.Info("record file exists: %s", recFile)
				a.emitInfo(eventRecordFileExists)

				if userAction.Auto {
					a.AutoInput(serial, userAction)
				}

				break
			}
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}()

	return
}

func (a *Api) AutoInput(serial string, userAction UserAction) error {
	if userAction.Type == "swipe" {

		a.InputSwipe(serial, adb.SwipeEvent{
			Sx:    userAction.X,
			Sy:    userAction.Y,
			Dx:    userAction.Tx,
			Dy:    userAction.Ty,
			Speed: userAction.Speed,
		})
	} else if userAction.Type == "click" {
		a.InputTap(serial, adb.TapEvent{
			X: userAction.X,
			Y: userAction.Y,
		})
	}
	return nil
}

// 发送拖动事件
func (a *Api) InputSwipe(serial string, sw adb.SwipeEvent) error {
	log.Infof("get input swipe event： %v on %s", sw, serial)
	device := adb.GetDevice(serial)
	return device.InputSwipe(sw)
}

// 发送点击事件
func (a *Api) InputTap(serial string, tap adb.TapEvent) error {
	log.Infof("get input tap event： %v on %s", tap, serial)
	device := adb.GetDevice(serial)
	return device.InputTap(tap)
}

func (a *Api) StopRunner() error {
	a.cmdRunner.CancelFunc()
	return nil
}

// Start 启动延迟测试
func (a *Api) Start(serial string, recordSecond int64, userAction UserAction) error {
	err := a.StartRecord(serial, userAction)
	if err != nil {
		a.emitInfo(eventRecordStartError)
		return err
	}
	// wait for filish
	time.Sleep(time.Duration(recordSecond) * time.Second)
	err = a.StopScrcpyServer(serial)
	if err != nil {
		a.emitInfo(eventRecordStopError)
		return err
	}
	err = a.StartTransform()
	if err != nil {
		a.emitInfo(eventTransformStartError)
		return err
	}
	return nil
}

func (a *Api) StopTransform() error {
	return nil
}

func (a *Api) ListRecords() ([]fs.RecordFile, error) {
	root, _ := fs.GetExecuteRoot()
	workDir := filepath.Join(root, "cache")
	files, err := fs.GetRecordFiles(workDir)
	if err != nil {
		log.Errorf("ListRecords error: %v", err)
		return []fs.RecordFile{}, nil
	}
	return files, nil
}

// StartWithVideo 启动延迟测试
func (a *Api) StartWithVideo(videoPath string) error {
	err := a.Transform(videoPath)
	if err != nil {
		a.emitInfo(eventTransformStartError)
		return err
	}
	return nil
}

// 停止 scrcpy server
func (a *Api) StopScrcpyServer(serial string) error {
	device := adb.GetDevice(serial)
	err := device.KillScrcyServer()
	if err != nil {
		a.emitInfo(eventRecordStopError)
		return err
	}
	a.emitInfo(eventRecordFilish)
	return nil
}

//func (a *Api) StopRecord(serial string) error {
//	log.Infof("stop monitor")
//	a.emitInfo(eventRecordFilish)
//	return a.Cmd.Kill()
//
//}

//	func (a *Api) StopProcessing() error {
//		log.Printf("stop monitor")
//		a.emitInfo(eventRecordFilish)
//
//		return a.Cmd.Kill()
//	}
func (a *Api) Transform(videoPath string) error {
	//srcVideoPath := filepath.Join(a.VideoDir, recordFile)
	a.VideoDir, a.ImagesDir = fs.CreateWorkDir()
	_, err := fs.Copy(videoPath, filepath.Join(a.VideoDir, "rec.mp4"))
	if err != nil {
		log.Error(err)
	}
	destImagePath := filepath.Join(a.ImagesDir, "%4d.png")
	a.emitInfo(eventTransformStart)
	cmd, err := cmd.StartFFmpeg(videoPath, destImagePath)
	if err != nil {
		log.Error(err)
		a.emitInfo(eventTransformStartError)
		return err
	}

	a.emitInfo(eventTransformFilish)
	a.Cmd = cmd
	return nil
}

func (a *Api) StartTransform() (rerr error) {
	defer func() {
		if rerr != nil {
			log.Error(rerr)
			a.emitInfo(eventTransformStartError)
		}
	}()

	srcVideoPath := filepath.Join(a.VideoDir, recordFile)
	destImagePath := filepath.Join(a.ImagesDir, "%4d.png")
	a.emitInfo(eventTransformStart)
	log.Infof("prepare data")

	runner, rerr := cmd.FFmpegStart(srcVideoPath, destImagePath, func() error {
		a.emitInfo(eventTransformFilish)
		return nil
	})
	a.cmdRunner = runner
	return
}

func (a *Api) GetFirstImageInfo() (core.ImageInfo, error) {
	firstImage := filepath.Join(a.ImagesDir, firstImageFile)
	mInfo, err := core.GetImageInfo(firstImage)
	log.Infof("Get first image: %v", mInfo)
	if err != nil {
		return mInfo, err
	}
	// path to uri format for FileLoader on win
	if fs.IsWindowsDrivePath(mInfo.Path) {
		mInfo.Path = "/" + strings.ReplaceAll(mInfo.Path, "\\", "/")
	}
	return mInfo, nil
}

func (a *Api) StartAnalyse(imageRect core.ImageRectInfo, diffScore int) error {
	// log.Infof("current rect: %v", imageRect)
	// log.Infof("workdir: %s", a.ImagesDir)
	// a.emitInfo(eventAnalyseStart)
	// responseTimes, _ := core.CalcTime(a.ImagesDir, imageRect, diffScore)

	dm := core.NewDelayMonitor()
	dm.VideoPath = filepath.Join(a.VideoDir, recordFile)
	dm.ImagesFolder = a.ImagesDir
	dm.PointerRect = image.Rect(0, 0, 100, 35) // 指针位置观察点
	dm.SceneRect = imageRect
	costTime, err := dm.Run()
	if err != nil {
		fmt.Printf("delay monitor run failed:%v", err)
	}
	a.emitData(eventAnalyseFilish, costTime)
	return nil
}

func (a *Api) emitInfo(eventName string) {
	runtime.EventsEmit(a.ctx, eventName)
}

func (a *Api) emitData(eventName string, data interface{}) {
	runtime.EventsEmit(a.ctx, eventName, data)
}

//func (a *Api) StopTransform() {
//	a.Cmd.Kill()
//}

func (a *Api) GetImageFiles() ([]string, error) {
	var imgs []string
	imgs, err := fs.GetImageFiles(a.ImagesDir, imgs)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return imgs, nil
}

// UploadFile 文件上传
func (a *Api) UploadFile(filePath string) error {
	err := upload.UploadFile(filePath)
	if err != nil {
		log.Errorf("upload file failed: %s", err)
		return err
	}
	return nil
}

// VersionInfo ...
type VersionInfo struct {
	Version        string `json:"version"`
	CommitShortSHA string `json:"commitShortSHA"`
	BuildTimestamp string `json:"buildTimestamp"`
}

// GetVersionInfo ...
func (a *Api) GetVersionInfo() VersionInfo {
	return VersionInfo{
		Version:        lighttestVer.Version,
		CommitShortSHA: lighttestVer.CommitShortSHA,
		BuildTimestamp: lighttestVer.BuildTimestamp,
	}
}

// UpdateInfo ...
type UpdateInfo struct {
	LatestVersion string `json:"latestVersion"`
	NeedUpdate    bool   `json:"needUpdate"`
	Err           string `json:"err,omitempty"`
}

// 获取当前版本号
func (a *Api) GetCurrentVersion() string {
	return lighttestVer.Version
}

// CheckUpdate ...
func (a *Api) CheckUser() (lighttestservice.UserInfo, error) {
	// userInfo, err := upload.GetUserInfo()
	return upload.GetUserInfo()
}

// CheckUpdate ...
func (a *Api) CheckUpdate() UpdateInfo {
	mgr := lighttestUpdate.LighttestServiceUpdateManager{Endpoint: lighttestServiceEndpoint}
	latestVersion, needUpdate, err := mgr.NeedUpdate(appName, lighttestVer.Version)
	if err != nil {
		return UpdateInfo{Err: err.Error()}
	}
	log.Infof("check last version: %s need update: %t", latestVersion, needUpdate)
	return UpdateInfo{
		LatestVersion: latestVersion,
		NeedUpdate:    needUpdate,
	}
}

// DoUpdate ...
func (a *Api) DoUpdate(version string) {
	mgr := lighttestUpdate.LighttestServiceUpdateManager{Endpoint: lighttestServiceEndpoint}
	go func() {
		log.Infof("upgrade to last version: %s", version)
		err := mgr.DoUpdate(appName, version)
		if err != nil {
			log.Errorf("update error: %v", err)
			a.emitData(eventUpdateError, err.Error())
		} else {
			a.emitInfo(eventUpdateSuccess)
		}
	}()
}

func (a *Api) ClearMobleCache() {
	fs.ClearCacheDir("mobile")
}

func (a *Api) ClearPCCache() {
	fs.ClearCacheDir("pc")
}
