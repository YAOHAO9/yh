package connector

import (
	"github.com/YAOHAO9/pine/rpc/client/clientmanager"
	"github.com/YAOHAO9/pine/rpc/context"
	"github.com/YAOHAO9/pine/rpc/handler/serverhandler"
	"github.com/YAOHAO9/pine/rpc/message"
	"github.com/YAOHAO9/pine/rpc/session"
	"github.com/YAOHAO9/pine/service/compressservice"
	"github.com/YAOHAO9/pine/util"
	"github.com/sirupsen/logrus"
)

// SysHandlerMap 系统PRC枚举
var SysHandlerMap = struct {
	PushMessage   string
	UpdateSession string
	RouterRecords string
	GetSession    string
	Kick          string
	BroadCast     string
}{
	PushMessage:   "__PushMessage__",
	UpdateSession: "__UpdateSession__",
	RouterRecords: "__RouterRecords__",
	GetSession:    "__GetSession__",
	Kick:          "__Kick__",
	BroadCast:     "__BroadCast__",
}

func init() {

	// 更新Session
	serverhandler.Manager.Register(SysHandlerMap.UpdateSession, func(rpcCtx *context.RPCCtx, data map[string]string) {
		if rpcCtx.Session == nil {
			logrus.Error("Session 为 nil")
			return
		}
		connection := GetConnection(rpcCtx.Session.UID)
		if connection == nil {
			logrus.Warn("无效的UID(", rpcCtx.Session.UID, ")没有找到对应的客户端连接")
			return
		}

		for key, value := range data {
			connection.data[key] = value
		}

		if rpcCtx.GetRequestID() > 0 {
			rpcCtx.SendMsg([]byte("1"))
		}

	})

	// 推送消息
	serverhandler.Manager.Register(SysHandlerMap.PushMessage, func(rpcCtx *context.RPCCtx, data *message.PineMsg) {
		connection := GetConnection(rpcCtx.Session.UID)
		if connection == nil {
			logrus.Warn("无效的UID(", rpcCtx.Session.UID, ")没有找到对应的客户端连接")
			return
		}

		if len([]byte(data.Route)) == 1 {
			client := clientmanager.GetClientByID(rpcCtx.From)
			if client != nil {
				code := compressservice.Server.GetCodeByKind(client.ServerConfig.Kind)
				data.Route = string([]byte{code}) + data.Route
			}
		}

		connection.notify(data)

		if rpcCtx.GetRequestID() > 0 {
			rpcCtx.SendMsg([]byte("1"))
		}
	})

	// 获取路由记录
	serverhandler.Manager.Register(SysHandlerMap.RouterRecords, func(rpcCtx *context.RPCCtx, hash []string) {
		logrus.Warn(hash)
	})

	// 获取Session
	serverhandler.Manager.Register(SysHandlerMap.GetSession, func(rpcCtx *context.RPCCtx, data struct {
		CID string
		UID string
	}) {
		connection := GetConnection(data.UID)
		var session *session.Session
		if connection == nil {
			rpcCtx.SendMsg([]byte{})
		} else {
			session = connection.GetSession()
			rpcCtx.SendMsg(util.ToBytes(session))
		}

	})

	// 踢下线
	serverhandler.Manager.Register(SysHandlerMap.Kick, func(rpcCtx *context.RPCCtx, data []byte) {
		connection := GetConnection(rpcCtx.Session.UID)
		if connection == nil {
			return
		}
		connection.Kick(data)
	})

	// 广播
	serverhandler.Manager.Register(SysHandlerMap.BroadCast, func(rpcCtx *context.RPCCtx, notify *message.PineMsg) {
		for _, connection := range connStore {
			connection.notify(notify)
		}
	})

}
