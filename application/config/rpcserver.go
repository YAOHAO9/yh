package config

// ==========================================
// ServerConfig
// ==========================================
var serverConfig *RPCServerConfig

// RPCServerConfig 服务器配置 配置文件
type RPCServerConfig struct {
	ClusterName  string `validate:"required"`
	ID          string `validate:"required"`
	Kind        string `validate:"required"`
	Host        string `validate:"required"`
	Port        uint32 `validate:"gte=1,lte=65535"`
	IsConnector bool
	Token       string `validate:"required"`
}

// SetSeRPCrverConfig 保存服务器配置
func SetRPCServerConfig(sc *RPCServerConfig) {
	serverConfig = sc
}

// GetServerConfig 获取服务器配置
func GetServerConfig() *RPCServerConfig {
	return serverConfig
}
