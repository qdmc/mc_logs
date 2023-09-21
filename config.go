package mc_logs

import (
	"github.com/sirupsen/logrus"
	"sync"
)

var logConf *Config
var confOnce sync.Once

func init() {
	logConf = getConfig()
}

func getConfig() *Config {
	confOnce.Do(func() {
		logConf = defaultLogConfig()
	})
	return logConf
}

type Config struct {
	Level          logrus.Level
	IsStdout       bool
	FileStore      bool
	FileHookConfig *FileHookConf
}

func (c *Config) SetLevel(l uint32) *Config {
	if l <= 6 {
		c.Level = logrus.Level(l)
	}
	return c
}
func (c *Config) SetLevelString(level string) *Config {
	switch level {
	case "panic":
		c.Level = logrus.PanicLevel
	case "fatal":
		c.Level = logrus.FatalLevel
	case "error":
		c.Level = logrus.ErrorLevel
	case "warn":
		c.Level = logrus.WarnLevel
	case "info":
		c.Level = logrus.InfoLevel
	case "debug":
		c.Level = logrus.DebugLevel
	case "trace":
		c.Level = logrus.TraceLevel
	}
	return c
}

type FileHookConf struct {
	FilePath   string
	FileName   string
	MaxSaveDay uint
}

func defaultFileHookConf() *FileHookConf {
	return &FileHookConf{
		FilePath:   "./logs",
		FileName:   "log",
		MaxSaveDay: 7,
	}
}

func defaultLogConfig() *Config {
	return &Config{
		Level:          logrus.InfoLevel,
		IsStdout:       true,
		FileStore:      true,
		FileHookConfig: defaultFileHookConf(),
	}
}

func GetConf() *Config {
	return getConfig()
}

func SetConfig(c *Config) {
	if c != nil {
		logConf = c
	}
}
