package adb

func SetPointerLocationOn(d *Device) error {
	cmd := Cmd{Args: []string{
		"-s", d.Serial,
		"settings put system pointer_location 1",
	}}
	err := cmd.Run()
	return err
}

func SetPointerLocationOff(d *Device) error {
	// c := strings.Split("settings put system pointer_location 0", " ")
	cmd := Cmd{Args: []string{
		"-s", d.Serial,
		"settings put system pointer_location 0",
	}}
	err := cmd.Run()
	return err
}
