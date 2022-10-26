// Copyright (C) 2015 The Android Open Source Project
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package adb

import (
	"errors"
	"fmt"
	"strings"
)

// DeviceState represents the last queried state of an Android device.
type DeviceState int

// binary: DeviceState#Offline = offline
// binary: DeviceState#Online = device
// binary: DeviceState#Unauthorized = unauthorized
const (
	Offline DeviceState = iota
	Online
	Unauthorized
)

// ErrDeviceNotRooted is returned by Device.Root when the device is running a
// production build as is not 'rooted'.
var ErrDeviceNotRooted = errors.New("Device is not rooted")

// Device represents an attached Android device.
type Device struct {
	Serial      string      `json:"serial"`
	State       DeviceState `json:"state"`
	abi         string      `json:"abi,omitempty"`
	usb         string      `json:"usb,omitempty"`
	Product     string      `json:"product"`
	Model       string      `json:"model"`
	Device      string      `json:"device"`
	transportId int         `json:"transport_id"`
}

// Command returns a new Cmd that will run the command with the specified name
// and arguments on this device.
func (d *Device) Command(path string, args ...string) *Cmd {
	return &Cmd{
		Path:   path,
		Args:   args,
		Device: d,
	}
}

// Root restarts adb as root. If the device is running a production build then
// Root will return ErrDeviceNotRooted.
func (d *Device) Root() error {
	cmd := Cmd{Args: []string{"root"}}
	res, err := cmd.Call()
	switch strings.TrimSpace(res) {
	case "adbd cannot run as root in production builds":
		return ErrDeviceNotRooted
	case "restarting adbd as root", "adbd is already running as root":
		return err
	default:
		if err == nil {
			return errors.New(res)
		} else {
			return err
		}
	}
}

// SELinuxEnforcing returns true if the device is currently in a
// SELinux enforcing mode, or false if the device is currently in a SELinux
// permissive mode.
func (d *Device) SELinuxEnforcing() (bool, error) {
	res, err := d.Command("getenforce").Call()
	return strings.Contains(strings.ToLower(res), "enforcing"), err
}

// SetSELinuxEnforcing changes the SELinux-enforcing mode.
func (d *Device) SetSELinuxEnforcing(enforce bool) error {
	if enforce {
		return d.Command("setenforce", "1").Run()
	} else {
		return d.Command("setenforce", "0").Run()
	}
}

// StartActivity launches the specified action.
func (d *Device) StartActivity(a Action) error {
	return d.Command("am", "start",
		"-S", // Force-stop the target app before starting the activity
		"-W", // Wait for the activity to start
		"-a", a.Name,
		"-n", a.Package.Name+"/"+a.Activity).Run()
}

// String returns a string representing the device.
func (d *Device) Abi() string {
	if d.abi == "" {
		res, err := d.Command("getprop", "ro.product.cpu.abi").Call()
		if err == nil {
			d.abi = strings.TrimSpace(res)
		}
	}
	return d.abi
}

// String returns a string representing the device.
func (d *Device) String() string {
	return fmt.Sprintf("Device<%s>", d.Serial)
}

// Devices returns the list of serial numbers of all the attached Android
// devices.
func Devices() ([]*Device, error) {
	if adb == "" {
		return nil, ErrADBNotFound
	}
	cmd := Cmd{Args: []string{"devices", "-l"}}
	if out, err := cmd.Call(); err == nil {
		return parseDevices(out)
	} else {
		return nil, err
	}
}
func parseDevices(out string) ([]*Device, error) {
	a := strings.SplitAfter(out, "List of devices attached")
	if len(a) != 2 {
		return nil, errors.New("device list not returned")
	}
	lines := strings.Split(a[1], "\n")
	devices := make([]*Device, 0, len(lines))
	for _, line := range lines {
		fields := strings.Fields(line)
		switch len(fields) {
		case 0:
			continue
		case 2:
			state := DeviceState(0)
			if err := state.Parse(fields[1]); err != nil {
				return nil, err
			}
			device := &Device{
				Serial: fields[0],
				State:  state,
			}
			devices = append(devices, device)
		case 7:
			state := DeviceState(0)
			if err := state.Parse(fields[1]); err != nil {
				return nil, err
			}

			productStrs := strings.Split(fields[3], ":")
			product := productStrs[len(productStrs)-1]

			modelStrs := strings.Split(fields[4], ":")
			model := modelStrs[len(modelStrs)-1]

			deviceStrs := strings.Split(fields[5], ":")
			deviceName := deviceStrs[len(deviceStrs)-1]

			device := &Device{
				Serial:  fields[0],
				State:   state,
				Product: product,
				Model:   model,
				Device:  deviceName,
			}
			devices = append(devices, device)
		default:
			return nil, errors.New("could not parse device list")
		}
	}
	return devices, nil
}

// GetDevice get device instance by serial
func GetDevice(serial string) *Device {
	return &Device{Serial: serial}
}
