package config

// ==========================================
// LogConfig
// ==========================================
var logConfig *LogConfig

// LogConfig 日志 配置文件
type LogConfig struct {
	Type     string `validate:"oneof=Console File"`
	Level    string `validate:"oneof=Debug Info Warn Error"`
}

// SetLogConfig 保存日志配置
func SetLogConfig(lc *LogConfig) {
	logConfig = lc
}

// GetLogConfig 获取日志配置
func GetLogConfig() *LogConfig {
	return logConfig
}
