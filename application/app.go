package application

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/YAOHAO9/yh/application/config"
	RpcServer "github.com/YAOHAO9/yh/rpc/server"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Application app
type Application struct {
}

// Start start application
func (app Application) Start() {
	RpcServer.Start()
}

var app *Application

// CreateApp 创建app
func CreateApp() *Application {
	if !parseFlag() {
		return nil
	}
	if app != nil {
		return app
	}
	app = &Application{}

	return app
}

// 解析命令行参数
func parseFlag() bool {

	serverConfig := config.ServerConfig{}

	// 服务器配置
	flag.StringVar(&serverConfig.SystemName, "s", "dwc", "System name")
	flag.StringVar(&serverConfig.ID, "i", "connector-1", "Server id")
	flag.StringVar(&serverConfig.Kind, "k", "connector", "Server kind")
	flag.StringVar(&serverConfig.Host, "H", "127.0.0.1", "Server host")
	flag.StringVar(&serverConfig.Port, "p", "3110", "server port")
	flag.BoolVar(&serverConfig.IsConnector, "c", true, "Client port")
	flag.StringVar(&serverConfig.Token, "t", "ksYNdrAo", "System token")

	// Zookeeper 配置
	zkConfig := config.ZkConfig{}
	flag.StringVar(&zkConfig.Host, "zh", "127.0.0.1", "Zookeeper host")
	flag.StringVar(&zkConfig.Port, "zp", "2181", "Zookeeper port")

	// 是否参看启动参数帮助
	var help bool
	flag.BoolVar(&help, "h", false, "Help")

	// 解析
	flag.Parse()
	if help {
		// 打印帮助详情
		flag.Usage()
		return !help
	}

	// 打印命令行参数
	data, err := json.Marshal(serverConfig)
	if err != nil {
		panic(err)
	}
	fmt.Println("Server config: \n\t", string(data))

	// 保存配置
	config.SetServerConfig(&serverConfig)
	config.SetZkConfig(&zkConfig)

	return !help
}
