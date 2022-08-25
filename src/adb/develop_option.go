package adb

func (d *Device) SetPointerLocationOn() error {
	cmd := Cmd{Args: []string{
		"-s", d.Serial,
		"settings put system pointer_location 1",
	}}
	err := cmd.Run()
	return err
}

func (d *Device) SetPointerLocationOff() error {
	// c := strings.Split("settings put system pointer_location 0", " ")
	cmd := Cmd{Args: []string{
		"-s", d.Serial,
		"settings put system pointer_location 0",
	}}
	err := cmd.Run()
	return err
}

func GetDevice(serial string) *Device {
	return &Device{Serial: serial}
}
