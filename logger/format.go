package logger

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	colorRed    = 31
	colorYellow = 33
	colorBlue   = 36
	colorGray   = 37
)

func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.DebugLevel:
		return colorGray
	case logrus.WarnLevel:
		return colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorBlue
	}
}

// error format
type formatter struct{}

func (f formatter) Format(entry *logrus.Entry) ([]byte, error) {

	// output buffer
	b := &bytes.Buffer{}

	levelColor := getColorByLevel(entry.Level)
	fmt.Fprintf(b, "\x1b[%dm", levelColor)

	timestampFormat := time.RFC3339

	// 时间
	b.WriteString(entry.Time.Format(timestampFormat))

	// 日志等级
	b.WriteString(" [")
	b.WriteString(strings.ToUpper(entry.Level.String()))
	b.WriteString("] ")

	// message
	b.WriteString(strings.TrimSpace(entry.Message))

	stack, hasStack := entry.Data["Stack"]
	if hasStack {
		delete(entry.Data, "Stack")
	}
	// fields
	for key, value := range entry.Data {
		b.WriteString(fmt.Sprint(" [", key, ":", value, "] "))
	}

	if hasStack {
		// stack
		b.WriteString(fmt.Sprint("\nCall stack:\n", stack))
		b.WriteString("\x1b[0m\n")
	} else {
		// caller
		b.WriteString("\x1b[0m")
		b.WriteString(fmt.Sprint(" ", entry.Caller.File, ":", entry.Caller.Line, "\n"))

	}

	return b.Bytes(), nil

}
