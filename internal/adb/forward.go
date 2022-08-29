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
package adb

import "fmt"

// Port is the interface for sockets ports that can be forwarded from an Android
// Device to the local machine.
type Port interface {
	// adbForwardString returns the "port specification" for the adb forward
	// command. See: http://developer.android.com/tools/help/adb.html#commandsummary
	// This method is hidden as its only use is for adb command-line parameters,
	// which this package abstracts.
	adbForwardString() string
}

// TCPPort represents a TCP/IP port on either the local machine or Android
// device. TCPPort implements the Port interface.
type TCPPort int

func (p TCPPort) adbForwardString() string {
	return fmt.Sprintf("tcp:%d", p)
}

// NamedAbstractSocket represents an abstract UNIX domain socket name on either
// the local machine or Android device. NamedAbstractSocket implements the Port
// interface.
type NamedAbstractSocket string

func (p NamedAbstractSocket) adbForwardString() string {
	return fmt.Sprintf("localabstract:%s", p)
}

// Forward will forward the specified device Port to the specified local Port.
func (d *Device) Forward(local, device Port) error {
	cmd := Cmd{Args: []string{
		"-s", d.Serial,
		"forward", local.adbForwardString(), device.adbForwardString(),
	}}
	_, err := cmd.Call()
	return err
}
