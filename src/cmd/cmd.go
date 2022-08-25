package cmd

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"os/exec"
	"path/filepath"
	"syscall"
)

var ErrScrcpyNotFound = errors.New("scrcpy command not found on PATH")
var ErrFfmpegNotFound = errors.New("ffmpeg command not found on PATH")

var scrcpy string
var ffmpeg string
var jobs map[string]interface{}

var ctx context.Context
var cancel context.CancelFunc

func init() {
	if p, err := exec.LookPath("scrcpy"); err == nil {
		if p, err = filepath.Abs(p); err == nil {
			scrcpy = p
		}
	}

	if p, err := exec.LookPath("ffmpeg"); err == nil {
		if p, err = filepath.Abs(p); err == nil {
			ffmpeg = p
		}
	}
}

type Cmd struct {
	Path    string
	Args    []string
	ExecCmd *exec.Cmd
	Timeout int
	Ctx     context.Context
	Cancel  context.CancelFunc
	Pid     int
	Stdout  io.Writer
	Stderr  io.Writer
}

type Job struct {
	Serial      string
	CmdInstance Cmd
}

func (c *Cmd) Run() error {
	cmd := exec.Command(scrcpy, c.Args...)
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	c.ExecCmd = cmd
	return cmd.Run()
}

func (c *Cmd) Call() (string, error) {
	clone := *c
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
	// if err != nil && strings.Contains(stderr.String(), "error: device unauthorized.") {
	// 	err = ErrDeviceUnauthorized
	// }
	return stdout.String(), err
}

func (c *Cmd) ContextRun() {
	// timeout := 10
	// ctx := context.Background()
	// var cancel context.CancelFunc
	// ctx, cancel = context.WithCancel(context.Background())
	// defer cancel()

	cmd := exec.CommandContext(c.Ctx, scrcpy, c.Args...)
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	// Request the OS to assign process group to the new process, to which all its children will belong
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	c.ExecCmd = cmd

	cmd.Start()
	log.Printf("current pid: %d", cmd.Process.Pid)
	c.Pid = cmd.Process.Pid
	cmd.Wait()

	// log.Printf("current pid: %d", cmd.Process.Pid)
	// go func() {
	// 	cmd.Run()
	// }()
	// select {
	// case <-ctx.Done():
	// 	log.Println("context run down")
	// }
}

func (c *Cmd) Start() error {
	cmd := exec.Command(scrcpy, c.Args...)
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}
	c.ExecCmd = cmd
	return nil
}

func CancelProcess() {
	log.Fatalln("cancel processing")
	cancel()
}

func StartScrcpyRecord(serial string, path string) (Cmd, error) {
	cmd := Cmd{
		Args: []string{
			"-s", serial,
			"-r", path,
		},
	}

	err := cmd.Start()
	if err != nil {
		// log.Fatalf("cmd exec failed: %v", err)
		return cmd, err
	}
	return cmd, nil
}

func StopScrcpyRecord(serial string) error {
	ins := jobs[serial]
	if ins != nil {
		err := ins.(Cmd).ExecCmd.Process.Kill()
		return err
	}
	return nil
}

func StartTransform() {

}

// func RunCommand(timeout int, command string, params ...string) (stdout, stderr string, isKilled bool) {
// 	var stdoutBuf, stderrBuf bytes.Buffer
// 	cmd := exec.Command(command, params...)
// 	cmd.Stdout = &stdoutBuf
// 	cmd.Stderr = &stderrBuf
// 	cmd.Start()

// 	done := make(chan error)
// 	go func() {
// 		done <- cmd.Wait()
// 	}()
// 	after := time.After(time.Duration(timeout) * time.Second)
// 	select {
// 	case <-after:
// 		cmd.Process.Signal(syscall.SIGINT)
// 		time.Sleep(time.Duration(timeout) * time.Second)
// 		cmd.Process.Kill()
// 		isKilled = true
// 	case <-done:
// 		isKilled = false

// 	}
// 	stdout = string(bytes.TrimSpace(stdoutBuf.Bytes()))
// 	stderr = string(bytes.TrimSpace(stderrBuf.Bytes()))
// 	return
// }

// func execAsync(serial string, params ...string) (*exec.Cmd, error) {
// 	args := make([]string, 0, 8)
// 	if len(serial) > 0 {
// 		args = append(args, "-s", serial)
// 	}
// 	scrcpyCmdOnce.Do(getScrcpyCommand)

// 	cmd := exec.Command(scrcpyCmd, args...)
// 	cmd.Stderr = os.Stderr
// 	cmd.Stdout = os.Stdout

// 	if err := cmd.Start(); err != nil {
// 		return nil, err
// 	}
// 	return cmd, nil
// }

// func videoToImage(path string) error {
// 	return nil
// }

// func GetFileListByPath(dirPath, extension string) ([]string, error) {
// 	var files []string
// 	fs, err := ioutil.ReadDir(dirPath)
// 	if err != nil {
// 		return files, err
// 	}

// 	for _, file := range fs {
// 		ext := filepath.Ext(file.Name())
// 		if strings.EqualFold(ext, extension) {
// 			files = append(files, file.Name())
// 		}
// 	}
// 	return files, nil
// }

// // func HasSuffixIgnoreCapitalization(s, suffix string) bool {
// // 	return len(s) >= len(suffix) && strings.ToLower(s[len(s)-len(suffix):]) == suffix
// // }

// func ImageDiff(fileList []string) (int, error) {

// 	for index, fileName := range fileList {
// 		if index == 0 {
// 			continue
// 		}
// 		previousFileName := fileList[index-1]
// 		previousFile, _ := os.Open(path.Join("./out/image/", previousFileName))
// 		currentFile, _ := os.Open(path.Join("./out/image/", fileName))
// 		defer previousFile.Close()
// 		defer currentFile.Close()

// 		previousImage, _ := png.Decode(previousFile)
// 		previousHash, _ := goimagehash.ExtDifferenceHash(previousImage, 16, 16)
// 		// previousHash, _ := goimagehash.PerceptionHash(previousImage)
// 		currentImage, _ := png.Decode(currentFile)
// 		currentHash, _ := goimagehash.ExtDifferenceHash(currentImage, 16, 16)
// 		// currentHash, _ := goimagehash.PerceptionHash(currentImage)
// 		score, _ := currentHash.Distance(previousHash)
// 		log.Printf("file: %s", fileName)
// 		log.Printf("score: %d", score)
// 		// return score, err
// 	}
// 	return 100, nil
// }

// func main() {
// out, err, stat := RunCommand(5, scrcpyCmd, "-s", "b9f8ef93", "-Nr", "./out/video/file.mp4")
// log.Print("is killed: ", stat)
// log.Print("res: \n", out, "\n")
// log.Print("err: \n", err, "\n")

// out, err, stat = RunCommand(2*60, ffmpegCmd, "-i", "./out/video/file.mp4", `./out/image/%04d.png`)
// log.Print("is killed: ", stat)
// log.Print("res: \n", out, "\n")
// log.Print("err: \n", err, "\n")

// fileList, _ := GetFileListByPath("./out/image")
// if err != nil {
// 	log.Print("err: %v", err)
// }
// log.Print("files: %v", fileList)
// sort.Slice(fileList, func(i, j int) bool {
// 	return fileList[i] < fileList[j]
// })

// log.Printf("files: %v", fileList)
// ImageDiff(fileList)

// }
