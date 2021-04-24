package config

// ==========================================
// ZkConfig
// ==========================================
var zkConfig *ZooKeeperConfig

// ZooKeeperConfig zk 配置文件
type ZooKeeperConfig struct {
	Host string `validate:"required"`
	Port uint32 `validate:"gte=1,lte=65535"`
}

// SetZkConfig 配置zookeeper配置
func SetZkConfig(zc *ZooKeeperConfig) {
	zkConfig = zc
}

// GetZkConfig 获取zk配置
func GetZkConfig() *ZooKeeperConfig {
	return zkConfig
}
