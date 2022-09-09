//go:build windows
// +build windows

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
	"log"
	"os/exec"
	"syscall"
)

func (c *Cmd) Run() error {
	args := []string{}
	if c.Device != nil {
		args = append(args, "-s", c.Device.Serial)
	}
	if c.Path != "" {
		args = append(args, "shell", c.Path)
	}
	args = append(args, c.Args...)
	log.Printf("adb: %s", adb)
	log.Printf("args: %v", args)
	cmd := exec.Command(adb, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	return cmd.Run()
}

func (c *Cmd) BackendRun() error {
	args := []string{}
	if c.Device != nil {
		args = append(args, "-s", c.Device.Serial)
	}
	if c.Path != "" {
		args = append(args, "shell", c.Path)
	}
	args = append(args, c.Args...)
	log.Printf("adb: %s", adb)
	log.Printf("args: %v", args)
	cmd := exec.Command(adb, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	return cmd.Start()
}
