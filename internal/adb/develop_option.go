package adb

import (
	log "github.com/sirupsen/logrus"
	"strings"
)

// 开启指针位置显示
func (d *Device) SetPointerLocationOn() error {
	cmd := strings.Split("put system pointer_location 1", " ")
	return d.Command("settings", cmd...).Run()
}

// 关闭指针位置显示
func (d *Device) SetPointerLocationOff() error {
	cmd := strings.Split("put system pointer_location 0", " ")
	return d.Command("settings", cmd...).Run()
}

// 检查指针位置显示是否开启
func (d *Device) IsPointerLocationOn() (bool, error) {
	//cmd := strings.Split("get system pointer_location", " ")
	out, err := d.Command("settings get system pointer_location").Call()
	log.Infof("IsPointerLocationOn out : %s", out)
	log.Infof("IsPointerLocationOn out len : %d", len(out))
	if err != nil {
		return false, err
	}
	
	if strings.TrimSpace(out) == "1" {
		return true, nil
	}

	return false, nil
}
