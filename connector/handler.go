package connector

import (
	"github.com/YAOHAO9/pine/rpc/context"
	"github.com/YAOHAO9/pine/rpc/handler"
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
	handler.Manager.Register(HandlerMap.UpdateSession, func(rpcCtx *context.RPCCtx) *handler.Resp {
		connection := GetConnection(rpcCtx.Session.UID)
		if connection == nil {
			logrus.Error("无效的Uid", rpcCtx.Session.UID, "没有找到对应的客户端连接")
			return nil
		}

		if data, ok := rpcCtx.Data.(map[string]interface{}); ok {
			for key, value := range data {
				connection.data[key] = value
			}
		}

		return nil
	})

	// 推送消息
	handler.Manager.Register(HandlerMap.PushMessage, func(rpcCtx *context.RPCCtx) *handler.Resp {
		connection := GetConnection(rpcCtx.Session.UID)
		if connection == nil {
			logrus.Error("无效的Uid", rpcCtx.Session.UID, "没有找到对应的客户端连接")
			return nil
		}

		if notify, ok := rpcCtx.Data.(map[string]interface{}); ok {
			connection.notify(notify["Route"].(string), notify["Data"])
		}

		return nil
	})

}
