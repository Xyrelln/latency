package app

import (
	"bytes"
	"embed"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"op-latency-mobile/internal/fs"
	"op-latency-mobile/internal/logger"

	log "github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//var assets embed.FS

type FileLoader struct {
	http.Handler
}

func NewFileLoader() *FileLoader {
	return &FileLoader{}
}

func (h *FileLoader) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var err error
	requestedFilename := req.URL.Path
	// requestedFilename := strings.TrimPrefix(req.URL.Path, "/")
	println("Requesting file:", requestedFilename)
	if fs.IsWindowsDrivePathURI(requestedFilename) {
		requestedFilename = strings.Replace(requestedFilename, "/", "", 1)
		requestedFilename = strings.ReplaceAll(requestedFilename, "/", "\\")
	}

	fileData, err := os.ReadFile(requestedFilename)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf("Could not load file %s", requestedFilename)))
	}

	res.Write(fileData)
}

func Run(assets embed.FS) int {
	// Create an instance of the app structure
	appData, err := appDataLocation(appName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open add data directory: %v\n", err)
		return 1
	}
	defer crashlog(appData)

	app := NewApp()
	app.AppData = appData

	ops := &options.App{
		Title:     "latency-mobile",
		Width:     1080,
		Height:    720,
		MinWidth:  820,
		MinHeight: 620,
		MaxWidth:  1920,
		MaxHeight: 1080,
		Assets:    assets,
		Logger:    logger.WailsLogger{},
		// BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:     app.startup,
		OnShutdown:    app.shutdown,
		OnDomReady:    app.domready,
		AssetsHandler: NewFileLoader(),
		Bind: []interface{}{
			app,
		},
	}

	// Create application with options
	// err = wails.Run(ops)
	log.Error(wails.Run(ops))

	// if err != nil {
	// 	println("Error:", err.Error())
	// }

	return 0
}

func crashlog(appData string) {
	// if wails.BuildMode != cmd.BuildModeProd {
	// 	return
	// }
	if r := recover(); r != nil {
		if _, err := os.Stat(appData); os.IsNotExist(err) {
			os.MkdirAll(appData, 0700)
		}
		var b bytes.Buffer
		b.WriteString(fmt.Sprintf("%+v\n\n", r))
		buf := make([]byte, 1<<20)
		s := runtime.Stack(buf, true)
		b.Write(buf[0:s])
		ioutil.WriteFile(filepath.Join(appData, "crash.log"), b.Bytes(), 0644)
	}
}
