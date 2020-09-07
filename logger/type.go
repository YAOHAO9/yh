package logger

// LogTypeEnum 日志类型枚举
var LogTypeEnum = struct {
	Console        int
	File           int
	ConsoleAndFile int
}{
	Console:        1,
	File:           2,
	ConsoleAndFile: 3,
}
