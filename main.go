package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"trial/config"
	"trial/connector"
	"trial/zookeeper"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	// 解析命令行参数
	if parseFlag() {
		return
	}

	// 注册到zookeeper
	go regist()

	go startConnector()

	// 获取服务器配置
	serverConfig := config.GetServerConfig()

	// Handler
	fmt.Println("Rpc server started ws://" + serverConfig.Host + ":" + serverConfig.Port)
	http.HandleFunc("/", connector.WebSocketHandler)
	err := http.ListenAndServe(":"+serverConfig.Port, nil)
	fmt.Println("Rpc server start fail: ", err.Error())

}

// 注册到zookeeper
func regist() {

	time.Sleep(time.Millisecond * 100)
	zookeeper.Start()
}

func parseFlag() bool {

	serverConfig := config.ServerConfig{}

	// 服务器配置
	flag.StringVar(&serverConfig.SystemName, "s", "dwc", "System name")
	flag.StringVar(&serverConfig.ID, "i", "connector-1", "Server id")
	flag.StringVar(&serverConfig.Kind, "k", "connector", "Server kind")
	flag.StringVar(&serverConfig.Host, "H", "127.0.0.1", "Server host")
	flag.StringVar(&serverConfig.Port, "p", "3110", "server port")
	flag.StringVar(&serverConfig.ClientPort, "P", "", "Client port")
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

func startConnector() {
	// 获取服务器配置
	serverConfig := config.GetServerConfig()

	// start connector
	if serverConfig.ClientPort != "" {
		// Handler
		fmt.Println("Rpc server started ws://" + serverConfig.Host + ":" + serverConfig.ClientPort)
		http.HandleFunc("/", connector.WebSocketHandler)
		err := http.ListenAndServe(":"+serverConfig.ClientPort, nil)
		fmt.Println("Connector start fail: ", err.Error())
	}
}
