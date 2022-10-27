package adb

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type SwipeEvent struct {
	Sx    int `json:"sx"`
	Sy    int `json:"sy"`
	Dx    int `json:"dx"`
	Dy    int `json:"dy"`
	Speed int `json:"speed"`
}

type TapEvent struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Display struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func (d *Device) AutoSwipe(se SwipeEvent) error {
	cmd := d.Command(fmt.Sprintf("input swipe %d %d %d %d %d", se.Sx, se.Sy, se.Dx, se.Dy, se.Speed))
	err := cmd.BackendRun()
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (d *Device) InputSwipe(se SwipeEvent) error {
	cmd := d.Command(fmt.Sprintf("input swipe %d %d %d %d %d", se.Sx, se.Sy, se.Dx, se.Dy, se.Speed))
	err := cmd.BackendRun()
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (d *Device) InputTap(tap TapEvent) error {
	cmd := d.Command(fmt.Sprintf("input tap %d %d", tap.X, tap.Y))
	err := cmd.BackendRun()
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (d *Device) AutoTouch() error {
	cmd := d.Command("input swipe 1700 500 2100 500 500")

	err := cmd.BackendRun()
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (d *Device) PhysicalSize() (*Display, error) {
	cmd := d.Command("wm size")
	out, err := cmd.Call()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	display, err := parsePhysicalSize(out)
	if err != nil {
		log.Info(out) // check out when error
		log.Error(err)
		return nil, err
	}
	log.Infof("display: %v", display)
	return display, nil
}

func parsePhysicalSize(out string) (*Display, error) {
	var dispaly Display
	values := strings.Split(out, ":")
	if len(values) != 2 {
		return nil, fmt.Errorf("get physical size error")
	}

	sizes := strings.Split(values[len(values)-1], "x")
	dispaly.Width, _ = strconv.Atoi(strings.TrimSpace(sizes[0]))
	dispaly.Height, _ = strconv.Atoi(strings.TrimSpace(sizes[1]))
	return &dispaly, nil

}

// 获取屏幕分辨率
func (d *Device) DisplaySize() (*Display, error) {
	cmd := d.Command("dumpsys window displays | head -n 3")

	out, err := cmd.Call()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	display, err := parseSize(out)
	if err != nil {
		log.Info(out) // check out when error
		log.Error(err)
		return nil, err
	}
	log.Infof("display: %v", display)
	return display, nil
}

func parseSize(out string) (*Display, error) {
	var dispaly Display
	reAppSize := regexp.MustCompile(`app=[0-9]+x[0-9]+`)
	appSizes := reAppSize.Find([]byte(out))
	if len(appSizes) == 0 {
		return nil, errors.New("display info read error")
	}

	values := strings.Split(string(appSizes), "=")
	if len(values) != 2 {
		return nil, fmt.Errorf("display size read error")
	}

	sizes := strings.Split(values[len(values)-1], "x")
	dispaly.Width, _ = strconv.Atoi(strings.TrimSpace(sizes[0]))
	dispaly.Height, _ = strconv.Atoi(strings.TrimSpace(sizes[1]))
	return &dispaly, nil
}

// func parseSize(out string) (*Display, error) {
// 	var dispaly Display
// 	reW := regexp.MustCompile(`displayWidth=[0-9]+`)
// 	reH := regexp.MustCompile(`displayHeight=[0-9]+`)

// 	reW1 := regexp.MustCompile(`width=[0-9]+`)
// 	reH1 := regexp.MustCompile(`height=[0-9]+`)

// 	resultW := reW.Find([]byte(out))
// 	resultH := reH.Find([]byte(out))

// 	resultW1 := reW1.Find([]byte(out))
// 	resultH1 := reH1.Find([]byte(out))
// 	if len(resultW) > 0 && len(resultH) > 0 {
// 		valueW := strings.Split(string(resultW), "=")
// 		w1 := valueW[len(valueW)-1]
// 		width, _ := strconv.Atoi(w1)
// 		dispaly.Width = width

// 		valueH := strings.Split(string(resultH), "=")
// 		h1 := valueH[len(valueH)-1]
// 		height, _ := strconv.Atoi(h1)
// 		dispaly.Height = height
// 	} else if len(resultW1) > 0 && len(resultH1) > 0 {
// 		valueW := strings.Split(string(resultW1), "=")
// 		w1 := valueW[len(valueW)-1]
// 		width, _ := strconv.Atoi(w1)
// 		dispaly.Width = width

// 		valueH := strings.Split(string(resultH1), "=")
// 		h1 := valueH[len(valueH)-1]
// 		height, _ := strconv.Atoi(h1)
// 		dispaly.Height = height
// 	} else {
// 		return nil, errors.New("display info read failed")
// 	}
// 	return &dispaly, nil

// }

func (d *Device) IsHorizontal() error {
	cmd := d.Command("dumpsys input | grep SurfaceOrientation")

	err := cmd.Run()
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
