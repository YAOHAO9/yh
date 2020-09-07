package application

import (
	"os"

	"github.com/YAOHAO9/yh/application/config"
	"github.com/jessevdk/go-flags"
)

type options struct {
	*config.ServerConfig
	*config.ZkConfig
}

// 解析命令行参数
func parseFlag() {

	opts := &options{}
	_, err := flags.Parse(opts)

	if err != nil {
		os.Exit(0)
	}

	// 保存配置
	config.SetServerConfig(opts.ServerConfig)
	config.SetZkConfig(opts.ZkConfig)

}
