package application

import (
	"fmt"
	"strings"

	"github.com/YAOHAO9/pine/application/config"
	"github.com/YAOHAO9/pine/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ymlConfig struct {
	Server    *config.ServerConfig
	Zookeeper *config.ZkConfig
}

// 解析命令行参数
func parseFlag() {

	viper.SetDefault("Server.IsConnector", false)
	viper.SetDefault("Zookeeper.Host", "127.0.0.1")
	viper.SetDefault("Zookeeper.Port", "2181")

	// 保存配置
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		logrus.Error("读取配置文件失败: %v", err)
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	configYml := &ymlConfig{}
	viper.Unmarshal(configYml)
	config.SetServerConfig(configYml.Server)
	config.SetZkConfig(configYml.Zookeeper)

	fmt.Printf("%c[%dm%s%c[m\n", 0x1B, logger.Info, fmt.Sprintf("ServerConfig config: %+v", configYml.Server), 0x1B)
	fmt.Printf("%c[%dm%s%c[m\n", 0x1B, logger.Info, fmt.Sprintf("ZooKeeper config: %+v", configYml.Zookeeper), 0x1B)

}
