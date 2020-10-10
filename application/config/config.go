package config

// ==========================================
// ServerConfig
// ==========================================
var serverCOnfig *ServerConfig

// ServerConfig 服务器配置 配置文件
type ServerConfig struct {
	SystemName  string `short:"s" long:"sysName" description:"系统名称" required:"true" default:"yh"`
	ID          string `short:"i" long:"serverID" description:"当前服务器ID" required:"true"`
	Kind        string `short:"k" long:"serverKind" description:"服务器类型" required:"true"`
	Host        string `short:"H" long:"host" description:"Host" required:"true" default:"127.0.0.1"`
	Port        string `short:"P" long:"port" description:"Port" required:"true"`
	IsConnector bool   `short:"c" long:"isConnector" description:"是否是Connector服"`
	Token       string `short:"t" long:"token" description:"集群认证Token" required:"true"`
	LogType     string `short:"l" long:"logType" description:"日志类型" default:"Console" choice:"Console" choice:"File"` // 1、控制台输出 2、日志文件
	LogLevel    string `short:"L" long:"logLevel" description:"日志输出等级" default:"Debug" choice:"Debug" choice:"Info" choice:"Warn" choice:"Error"`
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
	Host string `long:"zhost" description:"zookeeper Host" required:"true" default:"127.0.0.1"`
	Port string `long:"zport" description:"Show verbose debug message" required:"true" default:"2181"`
}

// SetZkConfig 配置zookeeper配置
func SetZkConfig(zc *ZkConfig) {
	zkConfig = zc
}

// GetZkConfig 获取zk配置
func GetZkConfig() *ZkConfig {
	return zkConfig
}
