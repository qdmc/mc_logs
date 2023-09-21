package mc_logs

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
)

var logger *logrus.Logger
var once sync.Once

func GetOnce() *logrus.Logger {
	return logger
}

func initLog() *logrus.Logger {
	once.Do(func() {
		conf := getConfig()
		logger = logrus.New()
		logger.SetLevel(conf.Level)
		var logWriter io.Writer
		devNullWriter, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			panic(fmt.Sprintf("os.DevNull_error: %s", err.Error()))
		} else {
			if conf.IsStdout {
				logWriter = io.MultiWriter(devNullWriter, os.Stdout)
			} else {
				logWriter = io.MultiWriter(devNullWriter)
			}
		}
		logger.SetOutput(logWriter)
		logger.SetFormatter(new(DefaultJsonFormatter))
		if conf.FileStore {
			fileHook, err := getFileHook(conf)
			if err != nil {
				panic(fmt.Sprintf("fileHook_error: %s", err.Error()))
			}
			logger.AddHook(fileHook)
		}
	})
	return logger
}
