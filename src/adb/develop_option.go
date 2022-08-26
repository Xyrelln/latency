package adb

import "strings"

func (d *Device) SetPointerLocationOn() error {
	c := strings.Split("settings put system pointer_location 1", " ")
	cmd := Cmd{Args: append([]string{
		"-s", d.Serial}, c...,
	)}
	err := cmd.Run()
	return err
}

func (d *Device) SetPointerLocationOff() error {
	c := strings.Split("settings put system pointer_location 0", " ")
	cmd := Cmd{Args: append([]string{
		"-s", d.Serial}, c...,
	)}
	err := cmd.Run()
	return err
}

func GetDevice(serial string) *Device {
	return &Device{Serial: serial}
}
