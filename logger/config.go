package logger

import (
	"log"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

// LogType 当前日志类型
var LogType int

// SetLogMode 设置log模式
func SetLogMode(logType int) {
	LogType = logType

	logrus.AddHook(&errorHook{})
	logrus.SetReportCaller(true)

	if LogType == LogTypeEnum.File {
		// 产品模式
		logrus.SetLevel(logrus.ErrorLevel)
		logrus.SetFormatter(&logrus.JSONFormatter{})

		path := "D:\\Projects\\go-trial"
		writer, err := rotatelogs.New(
			path+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(path),
			rotatelogs.WithMaxAge(30*time.Second),
			rotatelogs.WithRotationTime(time.Duration(60)*time.Second),
		)

		if err != nil {
			log.Fatalf("create file log.txt failed: %v", err)
		}

		logrus.SetOutput(writer)
	}

	if LogType == LogTypeEnum.Console {
		// 开发模式
		logrus.SetLevel(logrus.TraceLevel)
		logrus.SetFormatter(ErrorFormatter{})
	}
}
