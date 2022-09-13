package cmd

import (
    "errors"
    "os"
    "os/exec"
    "path/filepath"
)

var ErrScrcpyNotFound = errors.New("scrcpy command not found on PATH")

var scrcpy string

func init() {
    // Fallback to searching on PATH.
    if p, err := exec.LookPath(scrcpyExecFile); err == nil {
        if p, err = filepath.Abs(p); err == nil {
            scrcpy = p
            return
        }
    }

    // Fallback to searching on CurrentDirectory.
    if execPath, err := os.Executable(); err == nil {
        p := filepath.Join(filepath.Dir(execPath), "lib", "scrcpy", scrcpyExecFile)
        if _, err := os.Stat(p); os.IsNotExist(err) {
            scrcpy = p
            return
        }
    }
}

func StartScrcpy(serial, recFile string) (*Cmd, error) {
    cmd := Cmd{
        Args: []string{
            scrcpy,
            "-s", serial,
            "--max-fps", "60",
            // "-n", // no-control
            // "-w", // stay awake
            "-Nr", recFile,
        },
    }

    if err := cmd.BackendRun(); err == nil {
        return &cmd, nil
    } else {
        return nil, err
    }
}

func IsScrcpyReady() error {
    if scrcpy == "" {
        return ErrScrcpyNotFound
    }
    return nil
}
