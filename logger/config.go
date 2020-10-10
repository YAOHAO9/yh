package logger

import (
	"log"
	"os"
	"time"

	"github.com/YAOHAO9/pine/application/config"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

// LogType 当前日志类型
var LogType string

// SetLogMode 设置log模式
func SetLogMode(logType string) {
	LogType = logType

	logrus.AddHook(&errorHook{})
	logrus.SetReportCaller(true)

	logLevel := logrus.DebugLevel // Default
	switch config.GetServerConfig().LogLevel {
	case LogLevelEnum.Debug: // Debug
		logLevel = logrus.DebugLevel
	case LogLevelEnum.Info: // Info
		logLevel = logrus.InfoLevel
	case LogLevelEnum.Warn: // Warn
		logLevel = logrus.WarnLevel
	case LogLevelEnum.Error: // Error
		logLevel = logrus.ErrorLevel
	}

	logrus.SetLevel(logLevel)
	if LogType == LogTypeEnum.File {
		logrus.SetFormatter(&logrus.JSONFormatter{})

		path, _ := os.Getwd()

		writer, err := rotatelogs.New(
			path+"/log/"+config.GetServerConfig().ID+"_%m_%d.log",
			rotatelogs.WithMaxAge(30*time.Second),
			rotatelogs.WithRotationTime(time.Duration(60)*time.Second),
		)

		if err != nil {
			log.Fatalf("create file log.txt failed: %v", err)
		}

		logrus.SetOutput(writer)
	} else {
		logrus.SetFormatter(ErrorFormatter{})
	}
}
