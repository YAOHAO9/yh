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

		serverProto := path.Join(pwd, "/proto/jsonproto/server.json")
		clientProto := path.Join(pwd, "/proto/jsonproto/client.json")
		if !checkFileIsExist(serverProto) || !checkFileIsExist(clientProto) {

		}

		serverProtoCentent, err1 := ioutil.ReadFile(serverProto)
		clientProtoCentent, err2 := ioutil.ReadFile(clientProto)

		if err1 != nil || err2 != nil {
			logrus.Error(err1, err2)
			return
		}

		result := map[string]string{
			"server": string(serverProtoCentent),
			"client": string(clientProtoCentent),
		}

		bytes, err := json.Marshal(result)
		if err != nil {
			logrus.Error(err)
			return
		}

		rpcCtx.SendMsg(bytes)
	})

}
