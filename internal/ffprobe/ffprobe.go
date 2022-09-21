package ffprobe

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	// "log"
	// "op-latency-mobile/internal/cmd"

	log "github.com/sirupsen/logrus"
)

var ffprobe string
var ErrFFprobeNotFound = errors.New("ffprobe command not found on PATH")

func init() {
	// Fallback to searching on CurrentDirectory.
	if execPath, err := os.Executable(); err == nil {
		p := filepath.Join(filepath.Dir(execPath), "lib", "ffprobe", ffprobExecFile)
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			ffprobe = p
			return
		} else {
			log.Errorf("ffprobe path check failed: %s, reason: %v ", p, err)
			// log.Errorf("ffmpeg path check failed: %s, reason: v%", p, err)
		}
	}

	// Fallback to searching on PATH.
	if p, err := exec.LookPath(ffprobExecFile); err == nil {
		if p, err = filepath.Abs(p); err == nil {
			ffprobe = p
			return
		}
	}

}

// SetFFProbeBinPath sets the global path to find and execute the ffprobe program
func SetFFProbeBinPath(newBinPath string) {
	ffprobe = newBinPath
}

// ProbeURL is used to probe the given media file using ffprobe. The URL can be a local path, a HTTP URL or any other
// protocol supported by ffprobe, see here for a full list: https://ffmpeg.org/ffmpeg-protocols.html
// This function takes a context to allow killing the ffprobe process if it takes too long or in case of shutdown.
// Any additional ffprobe parameter can be supplied as well using extraFFProbeOptions.
func ProbeURL(ctx context.Context, fileURL string, extraFFProbeOptions ...string) (data *ProbeData, err error) {
	args := append([]string{
		"-loglevel", "fatal",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
	}, extraFFProbeOptions...)

	// Add the file argument
	args = append(args, fileURL)

	cmd := exec.CommandContext(ctx, ffprobe, args...)
	cmd.SysProcAttr = procAttributes()

	return runProbe(cmd)
}

// ProbeReader is used to probe a media file using an io.Reader. The reader is piped to the stdin of the ffprobe command
// and the data is returned.
// This function takes a context to allow killing the ffprobe process if it takes too long or in case of shutdown.
// Any additional ffprobe parameter can be supplied as well using extraFFProbeOptions.
func ProbeReader(ctx context.Context, reader io.Reader, extraFFProbeOptions ...string) (data *ProbeData, err error) {
	args := append([]string{
		"-loglevel", "fatal",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
	}, extraFFProbeOptions...)

	// Add the file from stdin argument
	args = append(args, "-")

	cmd := exec.CommandContext(ctx, ffprobe, args...)
	cmd.Stdin = reader
	cmd.SysProcAttr = procAttributes()

	return runProbe(cmd)
}

// runProbe takes the fully configured ffprobe command and executes it, returning the ffprobe data if everything went fine.
func runProbe(cmd *exec.Cmd) (data *ProbeData, err error) {
	var outputBuf bytes.Buffer
	var stdErr bytes.Buffer

	cmd.Stdout = &outputBuf
	cmd.Stderr = &stdErr

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("error running %s [%s] %w", ffprobe, stdErr.String(), err)
	}

	if stdErr.Len() > 0 {
		return nil, fmt.Errorf("ffprobe error: %s", stdErr.String())
	}

	data = &ProbeData{}
	err = json.Unmarshal(outputBuf.Bytes(), data)
	if err != nil {
		fmt.Print(outputBuf.Bytes())
		return data, fmt.Errorf("error parsing ffprobe output: %w", err)
	}

	if data.Format == nil {
		return data, fmt.Errorf("no format data found in ffprobe output")
	}

	// Populate the old Tags structs for backwards compatibility purposes:
	if len(data.Format.TagList) > 0 {
		data.Format.Tags = &FormatTags{}
		data.Format.Tags.setFrom(data.Format.TagList)
	}
	for _, str := range data.Streams {
		str.Tags.setFrom(str.TagList)
	}

	return data, nil
}

// runProbe takes the fully configured ffprobe command and executes it, returning the ffprobe data if everything went fine.
func runProbe2(cmd *exec.Cmd) (data *PTSPackets, err error) {
	var outputBuf bytes.Buffer
	var stdErr bytes.Buffer

	cmd.Stdout = &outputBuf
	cmd.Stderr = &stdErr

	err = cmd.Run()

	if err != nil {
		return nil, fmt.Errorf("error running %s [%s] %v", ffprobe, stdErr.String(), err)
	}

	if stdErr.Len() > 0 {
		return nil, fmt.Errorf("ffprobe error: %s", stdErr.String())
	}

	// log.Info("pts packets data: %s", outputBuf.String())
	data = &PTSPackets{}

	err = json.Unmarshal(outputBuf.Bytes(), data)
	if err != nil {
		// fmt.Print(outputBuf.Bytes())
		log.Printf("Error getting data stdErr: %v", stdErr.Bytes())
		log.Printf("Error getting data outputBuf: %v", outputBuf.Bytes())
		return data, fmt.Errorf("error parsing ffprobe output: %w", err)
	}

	return data, nil
}

func ProbePTS(ctx context.Context, fileURL string, extraFFProbeOptions ...string) (data *PTSPackets, err error) {
	// ffprobe -v 0 -show_entries packet=pts,duration -of compact=p=0:nk=1 -read_intervals 999999 -select_streams v rec.mp4
	args := append([]string{
		"-loglevel", "fatal",
		// "-print_format", "json",
		"-v", "0",
		"-show_entries", "packet=pts,duration",
		"-of", "compact=p=0:nk=1",
		"-read_intervals", "999999",
		"-select_streams", "v",
		// fileURL,
		"-print_format", "json",
	}, extraFFProbeOptions...)

	// // Add the file argument
	args = append(args, fileURL)

	cmd := exec.CommandContext(ctx, ffprobe, args...)
	cmd.SysProcAttr = procAttributes()

	// cmd := cmd.Cmd{
	// 	Args: []string{
	// 		ffprobe,
	// 		"-loglevel", "fatal",
	// 		// "-print_format", "json",
	// 		"-v", "0",
	// 		"-show_entries", "packet=pts,duration",
	// 		"-of", "compact=p=0:nk=1",
	// 		"-read_intervals", "999999",
	// 		"-select_streams", "v",
	// 		fileURL,
	// 		"-print_format", "json",
	// 	},
	// }
	// out, err := cmd.Call()
	// if err != nil {
	// 	log.Errorf("getting probe data error: %s", out)
	// 	return nil, fmt.Errorf("getting probe data error: %s", err)
	// 	// fmt.Printf("cmd run error:%v", err)
	// }

	// data = &PTSPackets{}
	// log.Info("pts packets data: %s", out)
	// err = json.Unmarshal([]byte(out), data)
	// if err != nil {
	// 	log.Info("getting probe data out: %s", out)
	// 	log.Errorf("getting probe data error: %v", err)
	// 	return nil, fmt.Errorf("error parsing ffprobe output: %w", err)
	// }

	// return data, nil
	return runProbe2(cmd)
}

func IsFFprobeReady() error {
	if ffprobe == "" {
		return ErrFFprobeNotFound
	}
	return nil
}
