package cmd

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

// var scrcpyCmd = "scrcpy"
// var ffmpegCmd = "ffmpeg"
// var scrcpyCmdOnce sync.Once

// func getScrcpyCommand() {
// 	scrcpyEnv := os.Getenv("scrcpy")
// 	if len(scrcpyEnv) > 0 {
// 		scrcpyCmd = scrcpyEnv
// 	}
// }

// func getFfmpegCommand() {
// 	ffmpegEnv := os.Getenv("ffmpeg")
// 	if len(ffmpegEnv) > 0 {
// 		ffmpegCmd = ffmpegEnv
// 	}
// }

// func cmdExec(serial string, params ...string) error {
// 	if cmd, err := execAsync(serial, params...); err != nil {
// 		return err
// 	} else {
// 		return cmd.Wait()
// 	}
// }

var ErrScrcpyNotFound = errors.New("scrcpy command not found on PATH")
var ErrFfmpegNotFound = errors.New("ffmpeg command not found on PATH")

var scrcpy string
var ffmpeg string

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
	Path   string
	Args   []string
	Cmder  *exec.Cmd
	Stdout io.Writer
	Stderr io.Writer
}

func (c *Cmd) Run() error {
	cmd := exec.Command(scrcpy, c.Args...)
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	c.Cmder = cmd
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

func (c *Cmd) Kill() error {
	err := c.Cmder.Process.Kill()
	return err
}

func RunCommand(timeout int, command string, params ...string) (stdout, stderr string, isKilled bool) {
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := exec.Command(command, params...)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	cmd.Start()

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()
	after := time.After(time.Duration(timeout) * time.Second)
	select {
	case <-after:
		cmd.Process.Signal(syscall.SIGINT)
		time.Sleep(time.Duration(timeout) * time.Second)
		cmd.Process.Kill()
		isKilled = true
	case <-done:
		isKilled = false

	}
	stdout = string(bytes.TrimSpace(stdoutBuf.Bytes()))
	stderr = string(bytes.TrimSpace(stderrBuf.Bytes()))
	return
}

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

func StartScrcpyRecord(serial string, path string) error {
	cmd := Cmd{Args: []string{
		"-s", serial,
		"-r", path,
	}}

	msg, err := cmd.Call()
	log.Printf("scrcpy: %s", msg)
	return err
}

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
