package logger

import (
	"io"
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

var productMode bool

// SetLogMode 设置log模式
func SetLogMode(isProduct bool) {
	productMode = isProduct

	logrus.AddHook(&errorHook{})
	logrus.SetReportCaller(true)

	if productMode {
		// 产品模式
		logrus.SetLevel(logrus.ErrorLevel)
		logrus.SetFormatter(&logrus.JSONFormatter{})

		writer, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			log.Fatalf("create file log.txt failed: %v", err)
		}
		logrus.SetOutput(io.MultiWriter(writer))
	} else {
		// 开发模式
		logrus.SetLevel(logrus.InfoLevel)
		logrus.SetFormatter(formatter{})
	}
}
