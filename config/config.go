package config

// ==========================================
// ServerConfig
// ==========================================
var serverCOnfig *ServerConfig

// ServerConfig 服务器配置 配置文件
type ServerConfig struct {
	SystemName string
	ID         string
	Kind       string
	Host       string
	Port       string
	ClientPort string
	Token      string
}

// SetServerConfig 保存服务器配置
func SetServerConfig(sc *ServerConfig) {
	serverCOnfig = sc
}

// GetServerConfig 获取服务器配置
func GetServerConfig() *ServerConfig {
	return serverCOnfig
}

// ==========================================
// ZkConfig
// ==========================================
var zkConfig *ZkConfig

// ZkConfig zk 配置文件
type ZkConfig struct {
	Host string
	Port string
}

// SetZkConfig 配置zookeeper配置
func SetZkConfig(zc *ZkConfig) {
	zkConfig = zc
}

// GetZkConfig 获取zk配置
func GetZkConfig() *ZkConfig {
	return zkConfig
}
