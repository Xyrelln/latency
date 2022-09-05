//go:build windows
// +build windows

package cmd

import (
	"log"
	"os"
	"os/exec"
	"syscall"

	"golang.org/x/sys/windows"
)

func (c *Cmd) BackendRun(name string) error {
	cmd := exec.Command(name, c.Args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: windows.CREATE_NEW_PROCESS_GROUP}
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	c.execCmd = cmd
	log.Printf("cmd: %s", name)
	log.Printf("args: %v", c.Args)
	return cmd.Start()
}

func (c *Cmd) TaskKill(pid string) error {
	if taskkill == "" {
		return ErrTaskKillNotFound
	}
	cmd := Cmd{
		Args: []string{
			"/pid", pid,
		},
	}
	if err := cmd.Run(taskkill); err == nil {
		return nil
	} else {
		return err
	}
}

func (c *Cmd) Kill() error {
	if c.execCmd.Process != nil {
		signal := os.Interrupt
		// https://github.com/golang/go/issues/46345  windows not implemented signal
		// return c.TaskKill(strconv.Itoa(c.execCmd.Process.Pid))
		log.Printf("PID: %d", c.execCmd.Process.Pid)
		return terminateProc(c.execCmd.Process.Pid, signal)
		// return SendCtrlBreak(c.execCmd.Process.Pid)

	}
	return nil
}

// SendCtrlBreak sends a Ctrl-Break event to the process with id pid
func SendCtrlBreak(pid int) error {
	// d, e := syscall.LoadDLL("kernel32.dll")
	// if e != nil {
	// 	return e
	// }
	// p, e := d.FindProc("GenerateConsoleCtrlEvent")
	// if e != nil {
	// 	return e
	// }
	// r, _, e := p.Call(uintptr(syscall.CTRL_BREAK_EVENT), uintptr(pid))
	// if r == 0 {
	// 	return e // syscall.GetLastError()
	// }

	d, e := syscall.LoadDLL("kernel32.dll")
	if e != nil {
		log.Printf("LoadDLL: %v\n", e)
	}
	p, e := d.FindProc("GenerateConsoleCtrlEvent")
	if e != nil {
		log.Printf("FindProc: %v\n", e)
	}
	r, _, e := p.Call(syscall.CTRL_BREAK_EVENT, uintptr(pid))
	if r == 0 {
		log.Printf("GenerateConsoleCtrlEvent: %v\n", e)
	}

	return nil
}
