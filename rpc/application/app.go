package application

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"time"
	"trial/rpc/config"
	"trial/rpc/filter"
	"trial/rpc/filter/rpcfilter"
	"trial/rpc/handler"
	"trial/rpc/handler/rpchandler"
	"trial/rpc/msg"
	"trial/rpc/response"
	RpcServer "trial/rpc/server"
)

// Application app
type Application struct{}

// RegisterHandler 注册Handler
func (app Application) RegisterHandler(name string, f func(respCtx *response.RespCtx)) {
	handler.Manager.Register(name, f)
}

// RegisterRPCHandler 注册Handler
func (app Application) RegisterRPCHandler(name string, f func(respCtx *response.RespCtx)) {
	rpchandler.Manager.Register(name, f)
}

// RegisterHandlerBeforeFilter 注册before filter of handler
func (app Application) RegisterHandlerBeforeFilter(f func(respCtx *response.RespCtx) (next bool)) {
	filter.BeforeFilterManager().Register(f)
}

// RegisterHandlerAfterFilter 注册after filter of handler request
func (app Application) RegisterHandlerAfterFilter(f func(rm *msg.RPCResp) (next bool)) {
	filter.AfterFilterManager().Register(f)
}

// RegisterRPCBeforeFilter 注册before filter of rpc
func (app Application) RegisterRPCBeforeFilter(f func(respCtx *response.RespCtx) (next bool)) {
	rpcfilter.BeforeFilterManager().Register(f)
}

// RegisterRPCAfterFilter 注册after filter of rpc request
func (app Application) RegisterRPCAfterFilter(f func(rm *msg.RPCResp) (next bool)) {
	rpcfilter.AfterFilterManager().Register(f)
}

func init() {
	rand.Seed(time.Now().UnixNano())
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

// Start start application
func Start() {
	RpcServer.Start()
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
	flag.BoolVar(&serverConfig.IsConnector, "c", false, "Client port")
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
