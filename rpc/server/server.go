package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"trial/config"
	"trial/connector"
	"trial/rpc/msg"
	"trial/rpc/zookeeper"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// Start rpc server
func Start() {

	// 注册到zookeeper
	go registToZk()

	// 获取服务器配置
	serverConfig := config.GetServerConfig()
	// RPC server启动
	fmt.Println("Rpc server started ws://" + serverConfig.Host + ":" + serverConfig.Port)
	http.HandleFunc("/rpc", webSocketHandler)

	// 对客户端暴露的ws接口
	if serverConfig.IsConnector {
		http.HandleFunc("/", connector.WebSocketHandler)
	}
	// 开启并监听
	err := http.ListenAndServe(":"+serverConfig.Port, nil)
	fmt.Println("Rpc server start fail: ", err.Error())
}

// WebSocketHandler deal with ws request
func webSocketHandler(w http.ResponseWriter, r *http.Request) {

	// 建立连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("连接失败", err.Error())
		return
	}

	// 断开连接自动清除连接信息
	conn.SetCloseHandler(func(code int, text string) error {
		conn.Close()
		return nil
	})

	// 用户认证
	token := r.URL.Query().Get("token")

	// token校验
	if token != config.GetServerConfig().Token {
		fmt.Println("用户校验失败!!!")
		conn.CloseHandler()(0, "认证失败")
		return
	}

	// 开始接收消息
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			conn.CloseHandler()(0, err.Error())
			break
		}
		// 解析消息
		connectorMessage := &msg.Message{}
		err = json.Unmarshal(data, connectorMessage)

		if err != nil {
			// 不想写两遍，这里复用connector的SendFailMessage方法
			connector.SendFailMessage(conn, connectorMessage.Index, "无效的消息类型")
			continue
		}
		// connInfo.data <- data
		fmt.Println(config.GetServerConfig().ID, "收到消息", string(data))
		if connectorMessage.Index != 0 {
			connector.SendSuccessfulMessage(conn, connectorMessage.Index, fmt.Sprint(config.GetServerConfig().ID, "收到消息", string(data)))
		}
	}
}

// 注册到zookeeper
func registToZk() {
	time.Sleep(time.Millisecond * 100)
	zookeeper.Start()
}
