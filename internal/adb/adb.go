// Copyright (C) 2015 The Android Open Source Project
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Package adb provides an interface to the Android Debug Bridge.
package adb

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	// "android.googlesource.com/platform/tools/gpu/maker"
)

// ErrADBNotFound is returned when the ADB executable is not found.
var ErrADBNotFound = errors.New("ADB command not found on PATH")

// ErrDeviceUnauthorized is returned by ADB commands when the device has not
// authorized ADB debugging. Check the confirmation dialog on the device.
var ErrDeviceUnauthorized = errors.New("Device unauthorized")

// The path to the adb executable, or an empty string if the adb executable was
// not found.
var adb string

func init() {
	// Search for ADB using ANDROID_HOME
	if home := os.Getenv("ANDROID_HOME"); home != "" {
		path, err := filepath.Abs(filepath.Join(home, "platform-tools", "adb"))
		if err == nil {
			if _, err := os.Stat(path); err == nil {
				adb = path
				return
			}
		}
	}
	// Fallback to searching on PATH.
	if p, err := exec.LookPath("adb"); err == nil {
		if p, err = filepath.Abs(p); err == nil {
			adb = p
		}
	}
}

// Cmd represents a command that can be run on an Android device.
type Cmd struct {
	// Path is the path of the command to run on the device.
	//
	// If the string is empty, the command is treated as a ADB command for Device.
	Path string
	// Args holds the command line arguments to pass to the command.
	Args []string
	// The device this command should be run on. If nil, then any one of the
	// attached devices will execute the command.
	Device *Device
	// Stdout and Stderr specify the process's standard output and error.
	//
	// If either is nil, Run connects the corresponding file descriptor
	// to the null device (os.DevNull).
	//
	// If Stdout and Stderr are the same writer, at most one
	// goroutine at a time will call Write.
	Stdout io.Writer
	Stderr io.Writer
}

// Call starts the specified command and waits for it to complete, returning the
// all stdout as a string.
// The returned error is nil if the command runs, has no problems copying
// stdout and stderr, and exits with a zero exit status.
func (c *Cmd) Call() (string, error) {
	clone := *c // Don't change c's Stdout
	stdout := &bytes.Buffer{}
	if clone.Stdout != nil {
		clone.Stdout = io.MultiWriter(clone.Stdout, stdout)
	} else {
		clone.Stdout = stdout
	}
	stderr := &bytes.Buffer{}
	if clone.Stdout != nil {
		clone.Stderr = io.MultiWriter(clone.Stdout, stderr)
	} else {
		clone.Stderr = stderr
	}
	err := clone.Run()
	if err != nil && strings.Contains(stderr.String(), "error: device unauthorized.") {
		err = ErrDeviceUnauthorized
	}
	return stdout.String(), err
}
