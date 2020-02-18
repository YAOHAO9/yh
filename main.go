package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"time"
	"trial/rpc/app"
	"trial/rpc/config"
	"trial/rpc/msg"
	"trial/rpc/response"
	RpcServer "trial/rpc/server"

	"github.com/gorilla/websocket"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	// 解析命令行参数
	if parseFlag() {
		return
	}
	register()
	// 启动RPC server
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
		// 答应帮助详情
		flag.Usage()
		return help
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

	return help
}

func register() {
	app.RegisterHandler("handler", func(conn *websocket.Conn, forwardMessage *msg.ForwardMessage) {
		fmt.Println("UID:", forwardMessage.Session.UID)
		response.SendSuccessfulMessage(conn, false, forwardMessage.Msg.Index, config.GetServerConfig().ID+": 收到Handler消息")
	})
	app.RegisterRPC("rpc", func(conn *websocket.Conn, forwardMessage *msg.ForwardMessage) {
		response.SendSuccessfulMessage(conn, true, forwardMessage.Msg.Index, config.GetServerConfig().ID+": 收到Rpc消息")
	})
}
