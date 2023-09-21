package mc_logs

import (
	"errors"
	rotateLogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

func getFileHook(conf *Config) (*lfshook.LfsHook, error) {
	writer, err := getRotateWriter(conf)
	if err != nil {
		return nil, err
	}
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  writer,
		logrus.FatalLevel: writer,
		logrus.DebugLevel: writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.PanicLevel: writer,
		logrus.TraceLevel: writer,
	}
	return lfshook.NewHook(writeMap, new(DefaultJsonFormatter)), nil
}

func getRotateWriter(conf *Config) (*rotateLogs.RotateLogs, error) {
	var hookConf *FileHookConf
	if conf == nil || conf.FileHookConfig == nil {
		hookConf = defaultFileHookConf()
	} else {
		if conf.FileHookConfig.FilePath == "" {
			conf.FileHookConfig.FilePath = "./logs"
		}
		if conf.FileHookConfig.FileName == "" {
			conf.FileHookConfig.FileName = "log"
		}
		if conf.FileHookConfig.MaxSaveDay < 7 || conf.FileHookConfig.MaxSaveDay > 30 {
			conf.FileHookConfig.MaxSaveDay = 7
		}
		hookConf = conf.FileHookConfig
	}
	err := checkAndCreateDir(hookConf.FilePath)
	if err != nil {
		return nil, err
	}
	logFileFullPath := path.Join(hookConf.FilePath, hookConf.FileName)
	logWriter, err := rotateLogs.New(
		logFileFullPath+"_%Y%m%d.log",
		rotateLogs.WithLinkName(logFileFullPath),                               // 生成软链，指向最新日志文件
		rotateLogs.WithMaxAge(time.Duration(hookConf.MaxSaveDay)*24*time.Hour), // 文件最大保存时间
		rotateLogs.WithRotationTime(24*time.Hour),                              // 日志切割时间间隔
	)
	if err != nil {
		return nil, err
	}
	return logWriter, nil
}

func checkAndCreateDir(dir string) error {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return os.MkdirAll(dir, 0777)
	}
	if !info.IsDir() {
		return errors.New("%s file is exist")
	}
	return err
}
