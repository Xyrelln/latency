package adb

import "strings"

func (d *Device) SetPointerLocationOn() error {
	cmd := strings.Split("put system pointer_location 1", " ")
	return d.Command("settings", cmd...).BackendRun()
}

func (d *Device) SetPointerLocationOff() error {
	cmd := strings.Split("put system pointer_location 0", " ")
	return d.Command("settings", cmd...).BackendRun()
}
