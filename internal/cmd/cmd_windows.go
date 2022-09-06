//go:build windows
// +build windows

package cmd

import (
	"fmt"
	"golang.org/x/sys/windows"
	"log"
	// "os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

func init() {
	if p, err := exec.LookPath("scrcpy.exe"); err == nil {
		if p, err = filepath.Abs(p); err == nil {
			scrcpy = p
		}
	}

	if p, err := exec.LookPath("ffmpeg.exe"); err == nil {
		if p, err = filepath.Abs(p); err == nil {
			ffmpeg = p
		}
	}

}

func (c *Cmd) BackendRun(name string) error {
	cmd := exec.Command("cmd", c.Args...)
	// stdout and stderr, and exits with a zero exit status.
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: windows.CREATE_UNICODE_ENVIRONMENT}
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	c.execCmd = cmd
	log.Printf("cmd: %s", name)
	log.Printf("args: %v", c.Args)
	return cmd.Start()
}

// func (c *Cmd) TaskKill(pid string) error {
// 	if taskkill == "" {
// 		return ErrTaskKillNotFound
// 	}
// 	cmd := Cmd{
// 		Args: []string{
// 			"/pid", pid,
// 		},
// 	}
// 	if err := cmd.Run(taskkill); err == nil {
// 		return nil
// 	} else {
// 		return err
// 	}
// }

func (c *Cmd) Kill() error {
	if c.execCmd.Process != nil {
		// signal := os.Interrupt
		// https://github.com/golang/go/issues/46345  windows not implemented signal
		// https://docs.microsoft.com/en-us/windows/win32/procthread/process-creation-flags
		// https://github.com/mattn/goreman/blob/master/proc_windows.go#L16
		// return c.TaskKill(strconv.Itoa(c.execCmd.Process.Pid))
		log.Printf("PID: %d", c.execCmd.Process.Pid)
		for i := 0; i < 5; i++ {
			err := terminateProc(c.execCmd.Process.Pid)
			if err != nil {
				log.Print(err)
			}
			time.Sleep(time.Duration(100) * time.Millisecond)
		}
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

func StartScrcpyRecord(serial, recFile string) (*Cmd, error) {
	if scrcpy == "" {
		return nil, ErrScrcpyNotFound
	}
	cmd := Cmd{
		Args: []string{
			"/c", scrcpy,
			"-s", serial,
			// "-n", // no-control
			// "-w", // stay awake
			"-Nr", recFile,
		},
	}

	if err := cmd.BackendRun(scrcpy); err == nil {
		return &cmd, nil
	} else {
		return nil, err
	}
}

func terminateProc(pid int) error {
	fmt.Println("terminate proc")
	dll, err := windows.LoadDLL("kernel32.dll")
	if err != nil {
		return err
	}
	defer dll.Release()

	f, err := dll.FindProc("AttachConsole")
	if err != nil {
		return err
	}
	r1, _, err := f.Call(uintptr(pid))
	if r1 == 0 && err != syscall.ERROR_ACCESS_DENIED {
		return err
	}

	f, err = dll.FindProc("SetConsoleCtrlHandler")
	if err != nil {
		return err
	}
	r1, _, err = f.Call(0, 1)
	if r1 == 0 {
		return err
	}
	f, err = dll.FindProc("GenerateConsoleCtrlEvent")
	if err != nil {
		return err
	}
	r1, _, err = f.Call(windows.CTRL_C_EVENT, uintptr(pid))
	if r1 == 0 {
		return err
	}
	return nil
}
