package connector

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/YAOHAO9/pine/rpc/context"
	"github.com/YAOHAO9/pine/rpc/handler"
	"github.com/YAOHAO9/pine/rpc/message"
	"github.com/sirupsen/logrus"
)

// HandlerMap 系统PRC枚举
var HandlerMap = struct {
	PushMessage   string
	UpdateSession string
	FetchProto    string
}{
	PushMessage:   "__PushMessage__",
	UpdateSession: "__UpdateSession__",
	FetchProto:    "__FetchProto__",
}

var serverProtoCentent []byte
var clientProtoCentent []byte

func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func init() {
	// 更新Session
	handler.Manager.Register(HandlerMap.UpdateSession, func(rpcCtx *context.RPCCtx, data map[string]string) {
		connection := GetConnection(rpcCtx.Session.UID)
		if connection == nil {
			logrus.Error("无效的Uid", rpcCtx.Session.UID, "没有找到对应的客户端连接")
		}

		for key, value := range data {
			connection.data[key] = value
		}

		if rpcCtx.GetRequestID() > 0 {
			rpcCtx.SendMsg([]byte("1"))
		}

	})

	// 推送消息
	handler.Manager.Register(HandlerMap.PushMessage, func(rpcCtx *context.RPCCtx, data *message.PineMessage) {
		connection := GetConnection(rpcCtx.Session.UID)
		if connection == nil {
			logrus.Error("无效的Uid", rpcCtx.Session.UID, "没有找到对应的客户端连接")
		}

		connection.notify(data)

		if rpcCtx.GetRequestID() > 0 {
			rpcCtx.SendMsg([]byte("1"))
		}
	})

	// 推送消息
	handler.Manager.Register(HandlerMap.FetchProto, func(rpcCtx *context.RPCCtx, hash string) {
		pwd, _ := os.Getwd()

		serverProto := path.Join(pwd, "/proto/server.proto")
		clientProto := path.Join(pwd, "/proto/client.proto")

		var result map[string]string = map[string]string{}

		if serverProtoCentent == nil && checkFileIsExist(serverProto) {
			var err error
			serverProtoCentent, err = ioutil.ReadFile(serverProto)

			if err != nil {
				logrus.Error(err)
				return
			}
		}

		if clientProtoCentent == nil && checkFileIsExist(clientProto) {
			var err error
			clientProtoCentent, err = ioutil.ReadFile(clientProto)

			if err != nil {
				logrus.Error(err)
				return
			}

		}

		result["server"] = string(serverProtoCentent)
		result["client"] = string(clientProtoCentent)

		bytes, err := json.Marshal(result)
		if err != nil {
			logrus.Error(err)
			return
		}

		rpcCtx.SendMsg(bytes)
	})

}
