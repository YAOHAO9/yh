package config

// ==========================================
// ConnectorConfig
// ==========================================
var connectorConfig *ConnectorConfig

// ConnectorConfig 服务器配置 配置文件
type ConnectorConfig struct {
	Port uint32 `validate:"gte=1,lte=65535"`
}

// SetConnectorConfig 保存服务器配置
func SetConnectorConfig(cc *ConnectorConfig) {
	connectorConfig = cc
}

// GetConnectorConfig 获取服务器配置
func GetConnectorConfig() *ConnectorConfig {
	return connectorConfig
}
