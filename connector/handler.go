package connector

import (
	"github.com/YAOHAO9/pine/proto/proto/custom"
	"github.com/YAOHAO9/pine/rpc/context"
	"github.com/YAOHAO9/pine/rpc/handler"
	"github.com/YAOHAO9/pine/rpc/message"
	"github.com/sirupsen/logrus"
)

// HandlerMap 系统PRC枚举
var HandlerMap = struct {
	PushMessage   string
	UpdateSession string
}{
	PushMessage:   "__PushMessage__",
	UpdateSession: "__UpdateSession__",
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
			rpcCtx.SendMsg(1, message.StatusCode.Successful)
		}

	})

	// 推送消息
	handler.Manager.Register(HandlerMap.PushMessage, func(rpcCtx *context.RPCCtx, data *custom.Request) {
		connection := GetConnection(rpcCtx.Session.UID)
		if connection == nil {
			logrus.Error("无效的Uid", rpcCtx.Session.UID, "没有找到对应的客户端连接")
		}

		connection.notify(data)

		if rpcCtx.GetRequestID() > 0 {
			rpcCtx.SendMsg(1, message.StatusCode.Successful)
		}
	})

}
