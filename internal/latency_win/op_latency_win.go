//go:build windows
// +build windows

package latencywin

import (
	"context"
	"fmt"
	"image"
	"op-latency-mobile/internal/core"
	"op-latency-mobile/internal/fs"
	"op-latency-mobile/internal/latency_win/capture"
	"op-latency-mobile/internal/latency_win/input"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

// InputConf ...
type InputConf struct {
	Type   string `json:"type"` // mouse or keyboard
	IsAuto bool   `json:"isAuto"`
	KeyTap string `json:"keyTap,omitempty"`
	// MousePos [2]int `json:"mouse_pos,omitempty"`
}

// Config ...
type Config struct {
	InputConf          InputConf `json:"inputCconf,omitempty"`
	ImageDiffThreshold int       `json:"imageDiff_threshold"`
	Frames             int       `json:"frames,omitempty"`
	StartKey           string    `json:"startKey"`
	// OffsetMs           int       `json:"offset_ms"`
}

// OpLatency ...
func OpLatency(cfg Config, workdir string, printFunc func(string)) (capture.ScreenshotSeq, time.Time, error) {
	keyCode, err := input.KeyToVKCode(cfg.StartKey)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("start_key: %s 配置了无效的按键", cfg.StartKey)
	}

	printFunc(fmt.Sprintf("按 %s 开始测量", cfg.StartKey))
	_, err = input.WindowsInputEv{}.WaitKeyBoardEvent(context.Background(), keyCode)
	// time.Sleep(time.Second)
	if cfg.InputConf.Type == "keyboard" {
		printFunc(fmt.Sprintf("开始测量，请按 %s 键操作游戏", cfg.InputConf.KeyTap))
	} else {
		printFunc(fmt.Sprintf("开始测量，请鼠标点击操作游戏"))
	}

	errg, ctx := errgroup.WithContext(context.Background())
	inputTimeC := make(chan time.Time, 1)

	errg.Go(func() error {
		var (
			t   time.Time
			err error
			wev input.WindowsInputEv
		)
		if cfg.InputConf.Type == "mouse" {
			t, err = wev.WaitMouseEvent(ctx, input.WM_LBUTTONDOWN)
			printFunc(fmt.Sprintf("鼠标按下事件: %d", t.UnixMilli()))
		} else if cfg.InputConf.Type == "keyboard" {
			if cfg.InputConf.KeyTap == "" {
				return fmt.Errorf("配置 input_conf.key_tap 为空")
			}
			keyCode, err := input.KeyToVKCode(cfg.InputConf.KeyTap)
			if err != nil {
				return fmt.Errorf("input_conf.key_tap: %s 配置了无效的按键", cfg.InputConf.KeyTap)
			}
			t, err = wev.WaitKeyBoardEvent(ctx, keyCode)
			printFunc(fmt.Sprintf("键盘事件(%d): %d", keyCode, t.UnixMilli()))
		} else {
			return fmt.Errorf("input_conf.type 必须是 mouse 或 keyboard")
		}

		if err != nil {
			return fmt.Errorf("获取操控输入事件异常\n%w", err)
		}
		inputTimeC <- t
		close(inputTimeC)
		return nil
	})

	rsCap := capture.RustCapture{ExePath: filepath.Join(workdir, "rscapture.exe"), OutputDir: filepath.Join(workdir, "screenshots"), PrintFunc: printFunc}
	screenshotC := make(chan capture.ScreenshotSeq, 1)
	errg.Go(func() error {
		imgs, err := rsCap.CaptureScreenshots(cfg.Frames)
		if err != nil {
			return fmt.Errorf("截图错误\n%w", err)
		}
		screenshotC <- imgs
		close(screenshotC)
		return nil
	})

	if cfg.InputConf.IsAuto {
		time.Sleep(time.Second)
		if cfg.InputConf.Type == "mouse" {
			// input.MouseClick(cfg.InputConf.MousePos[0], cfg.InputConf.MousePos[1]) // 鼠标坐标不准
		} else {
			input.KeyboardPress(cfg.InputConf.KeyTap)
		}
	}

	if err := errg.Wait(); err != nil {
		return nil, time.Time{}, err
	}

	imgs := <-screenshotC
	inputTime := <-inputTimeC
	return imgs, inputTime, nil
}

// OpLatencyWindowsManager ...
type OpLatencyWindowsManager struct {
	Cfg               Config
	InputTime         time.Time
	ScreenshotSeq     capture.ScreenshotSeq
	CurrentImageIndex int
}

// Start ...
func (owm *OpLatencyWindowsManager) Start(cfg Config) error {
	owm.Cfg = cfg
	exeDir, err := fs.GetExecuteRoot()
	if err != nil {
		return err
	}

	screenshots, inputTime, err := OpLatency(cfg, exeDir, func(s string) { fmt.Println(s) })
	if err != nil {
		return err
	}

	owm.InputTime = inputTime
	owm.ScreenshotSeq = screenshots
	return nil
}

func calculateRectangle(guiWidth, guiHeight, x1, y1, x2, y2 int) image.Rectangle {
	rect := image.Rectangle{}
	return rect
}

// CalculateLatencyByImageDiff ...
func (owm *OpLatencyWindowsManager) CalculateLatencyByImageDiff(imageRect core.ImageRectInfo) (respIndex int, responseTime time.Time, latency time.Duration, err error) {
	// selectedRect := rp.imgRegionSelect.SelectedRegion()
	// selectedRect := image.Rectangle{} // TODO
	x0 := imageRect.X * imageRect.SourceWidth / imageRect.PreviewWidth
	y0 := imageRect.Y * imageRect.SourceHeight / imageRect.PreviewHeight
	x1 := imageRect.W*imageRect.SourceWidth/imageRect.PreviewWidth + x0
	y1 := imageRect.H*imageRect.SourceHeight/imageRect.PreviewHeight + y0
	selectedRect := image.Rect(x0, y0, x1, y1)

	respIndex, err = owm.ScreenshotSeq.FindImageHashResponseTime(selectedRect, owm.Cfg.ImageDiffThreshold, owm.InputTime)
	if err != nil {
		log.Errorf("识别画面响应时间失败: %v\n", err)
		return
	}
	responseTime = owm.ScreenshotSeq[respIndex].Time
	latency = responseTime.Sub(owm.InputTime)
	return respIndex, responseTime, latency, nil
}

// CalculateLatencyByIndex ...
func (owm *OpLatencyWindowsManager) CalculateLatencyByIndex(index int) (respIndex int, responseTime time.Time, latency time.Duration, err error) {

	responseTime = owm.ScreenshotSeq[index].Time
	latency = responseTime.Sub(owm.InputTime)
	return index, responseTime, latency, nil
}
