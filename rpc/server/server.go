package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/YAOHAO9/pine/application/config"
	"github.com/YAOHAO9/pine/connector"
	"github.com/YAOHAO9/pine/rpc/context"
	"github.com/YAOHAO9/pine/rpc/handler/clienthandler"
	"github.com/YAOHAO9/pine/rpc/handler/serverhandler"
	"github.com/YAOHAO9/pine/rpc/message"
	"github.com/YAOHAO9/pine/rpc/zookeeper"
	"github.com/YAOHAO9/pine/service/compressservice"
	"github.com/YAOHAO9/pine/util"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"

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
	logrus.Info("Rpc server started ws://" + serverConfig.Host + ":" + fmt.Sprint(serverConfig.Port))
	http.HandleFunc("/rpc", webSocketHandler)

	// 对客户端暴露的ws接口
	if serverConfig.IsConnector {
		http.HandleFunc("/", connector.WebSocketHandler)
	}
	// 开启并监听
	err := http.ListenAndServe(":"+fmt.Sprint(serverConfig.Port), nil)
	logrus.Error("Rpc server start fail: ", err.Error())
}

// WebSocketHandler deal with ws request
func webSocketHandler(w http.ResponseWriter, r *http.Request) {

	// 建立连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Error("连接失败", err.Error())
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
		logrus.Error("用户校验失败!!!")
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
		rpcMsg := &message.RPCMsg{}
		err = proto.Unmarshal(data, rpcMsg)

		rpcCtx := context.GenRespCtx(conn, rpcMsg)

		if err != nil {
			logrus.Error(err)
			continue
		}

		if rpcMsg.Type == message.RemoterTypeEnum.HANDLER {
			ok := clienthandler.Manager.Exec(rpcCtx)
			if !ok {
				if rpcCtx.GetRequestID() == 0 {
					logrus.Warn(fmt.Sprintf("NotifyHandler(%v)不存在", rpcCtx.GetHandler()))
				} else {
					logrus.Warn([]byte(fmt.Sprintf("Handler(%v)不存在", rpcCtx.GetHandler())))
				}
			}

		} else if rpcMsg.Type == message.RemoterTypeEnum.REMOTER {
			ok := serverhandler.Manager.Exec(rpcCtx)
			if !ok {
				if rpcCtx.GetRequestID() == 0 {
					logrus.Warn(fmt.Sprintf("NotifyRemoter(%v)不存在", rpcCtx.GetHandler()))
				} else {
					logrus.Warn([]byte(fmt.Sprintf("Remoter(%v)不存在", rpcCtx.GetHandler())))
				}
			}
		} else {
			logrus.Panic("无效的消息类型")
		}
	}
}

// 注册到zookeeper
func registToZk() {
	zookeeper.Start()
}

func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func init() {
	var serverProtoCentent []byte
	var clientProtoCentent []byte

	// 获取proto file
	clienthandler.Manager.Register("__FetchProto__", func(rpcCtx *context.RPCCtx, hash string) {
		pwd, _ := os.Getwd()

		serverProto := path.Join(pwd, "/proto/server.proto")
		clientProto := path.Join(pwd, "/proto/client.proto")

		var result = map[string]interface{}{}

		// server proto
		if serverProtoCentent == nil && checkFileIsExist(serverProto) {
			var err error
			serverProtoCentent, err = ioutil.ReadFile(serverProto)

			if err != nil {
				logrus.Error(err)
				return
			}
		}
		result["server"] = string(serverProtoCentent)

		// client proto
		if clientProtoCentent == nil && checkFileIsExist(clientProto) {
			var err error
			clientProtoCentent, err = ioutil.ReadFile(clientProto)

			if err != nil {
				logrus.Error(err)
				return
			}

		}
		result["client"] = string(clientProtoCentent)

		// handlers
		handlers := compressservice.Handler.GetHandlers()
		result["handlers"] = handlers

		// events
		result["events"] = compressservice.Event.GetEvents()

		rpcCtx.SendMsg(util.ToBytes(result))
	})
}
