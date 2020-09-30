package connector

import (
	"github.com/YAOHAO9/pine/rpc/context"
	"github.com/YAOHAO9/pine/rpc/handler"
	"github.com/sirupsen/logrus"
)

func init() {
	// 更新Session
	handler.Manager.Register(HandlerMap.UpdateSession, func(rpcCtx *context.RPCCtx) *handler.Resp {
		connInfo, ok := ConnMap[rpcCtx.Session.UID]
		if !ok {
			logrus.Error("无效的Uid", rpcCtx.Session.UID, "没有找到对应的客户端连接")
			return nil
		}

		if data, ok := rpcCtx.Data.(map[string]interface{}); ok {
			for key, value := range data {
				connInfo.data[key] = value
			}
		}

		return nil
	})

	// 推送消息
	handler.Manager.Register(HandlerMap.PushMessage, func(rpcCtx *context.RPCCtx) *handler.Resp {
		connInfo, ok := ConnMap[rpcCtx.Session.UID]
		if !ok {
			logrus.Error("无效的Uid", rpcCtx.Session.UID, "没有找到对应的客户端连接")
			return nil
		}

		if notify, ok := rpcCtx.Data.(map[string]interface{}); ok {
			connInfo.notify(notify["Route"].(string), notify["Data"])
		}

		return nil
	})

}
