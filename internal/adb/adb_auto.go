package adb

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type SwipeEvent struct {
	Sx    int `json:"sx"`
	Sy    int `json:"sy"`
	Dx    int `json:"dx"`
	Dy    int `json:"dy"`
	Speed int `json:"speed"`
}

type Display struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func (d *Device) AutoSwipe(se SwipeEvent) error {
	cmd := d.Command(fmt.Sprintf("input swipe %d %d %d %d %d", se.Sx, se.Sy, se.Dx, se.Dy, se.Speed))
	err := cmd.BackendRun()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (d *Device) AutoTouch() error {
	cmd := d.Command("input swipe 1700 500 2100 500 500")

	err := cmd.BackendRun()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (d *Device) DisplaySize() (*Display, error) {
	cmd := d.Command("dumpsys window displays | grep display")

	out, err := cmd.Call()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Print(out)
	display, err := parseSize(out)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Printf("display: %v", display)
	return display, nil
}

func parseSize(out string) (*Display, error) {
	var dispaly Display
	reW := regexp.MustCompile(`displayWidth=[0-9]+`)
	reH := regexp.MustCompile(`displayHeight=[0-9]+`)

	resultW := reW.Find([]byte(out))
	resultH := reH.Find([]byte(out))
	if len(resultW) > 0 && len(resultH) > 0 {
		valueW := strings.Split(string(resultW), "=")
		w1 := valueW[len(valueW)-1]
		width, _ := strconv.Atoi(w1)
		dispaly.Width = width

		valueH := strings.Split(string(resultH), "=")
		h1 := valueH[len(valueH)-1]
		height, _ := strconv.Atoi(h1)
		dispaly.Height = height
	} else {
		return nil, errors.New("dispaly info read fialed")
	}
	return &dispaly, nil

}

func (d *Device) IsHorizontal() error {
	cmd := d.Command("dumpsys input | grep SurfaceOrientation")

	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
