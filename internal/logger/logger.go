package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type myFormatter struct {
	log.TextFormatter
}

func (f *myFormatter) Format(entry *log.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("[%s][%s][%s:%d] %s\n", entry.Time.Format(time.RFC3339), strings.ToUpper(entry.Level.String()), entry.Caller.File, entry.Caller.Line, entry.Message)), nil
}

// SetupLog ...
func init() {
	p, err := os.Executable()
	if err != nil {
		panic(err)
	}
	logdir := filepath.Dir(p)
	log.SetOutput(&lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/latency-mobile.log", logdir),
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
		LocalTime:  true,
	})

	log.SetFormatter(&myFormatter{})
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(true)
}

// WailsLogger ...
type WailsLogger struct{}

// Print ...
func (l WailsLogger) Print(message string) {
	log.Print(message)
}

// Trace ...
func (l WailsLogger) Trace(message string) {
	log.Trace(message)
}

// Debug ...
func (l WailsLogger) Debug(message string) {
	log.Debug(message)
}

// Info ...
func (l WailsLogger) Info(message string) {
	log.Info(message)
}

// Warning ...
func (l WailsLogger) Warning(message string) {
	log.Warning(message)
}

// Error ...
func (l WailsLogger) Error(message string) {
	log.Error(message)
}

// Fatal ...
func (l WailsLogger) Fatal(message string) {
	log.Fatal(message)
}
