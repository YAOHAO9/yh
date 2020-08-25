package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/YAOHAO9/yh/application/config"
	"github.com/YAOHAO9/yh/rpc/connector"
	"github.com/YAOHAO9/yh/rpc/handler/clienthandler"
	"github.com/YAOHAO9/yh/rpc/handler/rpchandler"
	"github.com/YAOHAO9/yh/rpc/handler/syshandler"
	"github.com/YAOHAO9/yh/rpc/msg"
	"github.com/YAOHAO9/yh/rpc/response"
	"github.com/YAOHAO9/yh/rpc/zookeeper"

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
		rc := &msg.RPCMessage{}
		err = json.Unmarshal(data, rc)

		respCtx := &response.RespCtx{
			Conn:   conn,
			RPCMsg: rc,
		}

		if err != nil {
			respCtx.SendFailMessage("无效的消息类型")
			continue
		}

		if rc.Kind == msg.KindEnum.Sys {
			syshandler.Manager.Exec(respCtx)
		} else if rc.Kind == msg.KindEnum.RPC {
			rpchandler.Manager.Exec(respCtx)
		} else if rc.Kind == msg.KindEnum.Handler {
			clienthandler.Manager.Exec(respCtx)
		}
	}
}

// 注册到zookeeper
func registToZk() {
	zookeeper.Start()
}
