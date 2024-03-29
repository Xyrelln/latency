//go:build windows
// +build windows

package cmd

import (
	"golang.org/x/sys/windows"
	"syscall"
)

var cmdStart = []string{"cmd", "/c"}
var procAttrs = &windows.SysProcAttr{
	HideWindow: true,
	// CreationFlags: windows.CREATE_UNICODE_ENVIRONMENT | windows.CREATE_NEW_PROCESS_GROUP,
}

func procAttributes() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		HideWindow: true,
	}
}
