//go:build !windows
// +build !windows

package cmd

import (
	"syscall"

	"golang.org/x/sys/unix"
)

var cmdStart = []string{"/bin/sh", "-c"}
var procAttrs = &unix.SysProcAttr{Setpgid: true}

func procAttributes() *syscall.SysProcAttr {
	return nil
}
