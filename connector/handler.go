package connector

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/YAOHAO9/pine/rpc/client/clientmanager"
	"github.com/YAOHAO9/pine/rpc/context"
	"github.com/YAOHAO9/pine/rpc/handler/clienthandler"
	"github.com/YAOHAO9/pine/rpc/handler/serverhandler"
	"github.com/YAOHAO9/pine/rpc/message"
	"github.com/YAOHAO9/pine/rpc/session"
	"github.com/YAOHAO9/pine/service/compressservice"
	"github.com/YAOHAO9/pine/util"
	"github.com/sirupsen/logrus"
)

// HandlerMap 系统PRC枚举
var HandlerMap = struct {
	PushMessage   string
	UpdateSession string
	FetchProto    string
	RouterRecords string
	GetSession    string
	Kick          string
}{
	PushMessage:   "__PushMessage__",
	UpdateSession: "__UpdateSession__",
	FetchProto:    "__FetchProto__",
	RouterRecords: "__RouterRecords__",
	GetSession:    "__GetSession__",
	Kick:          "__Kick__",
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
	serverhandler.Manager.Register(HandlerMap.UpdateSession, func(rpcCtx *context.RPCCtx, data map[string]string) {
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
	serverhandler.Manager.Register(HandlerMap.PushMessage, func(rpcCtx *context.RPCCtx, data *message.PineMsg) {
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

	// 获取proto file
	clienthandler.Manager.Register(HandlerMap.FetchProto, func(rpcCtx *context.RPCCtx, hash string) {
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

		bytes, err := json.Marshal(result)
		if err != nil {
			logrus.Error(err)
			return
		}

		rpcCtx.SendMsg(bytes)
	})

	// 获取路由记录
	serverhandler.Manager.Register(HandlerMap.RouterRecords, func(rpcCtx *context.RPCCtx, hash []string) {
		logrus.Warn(hash)
	})

	// 获取Session
	serverhandler.Manager.Register(HandlerMap.GetSession, func(rpcCtx *context.RPCCtx, data struct {
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

	serverhandler.Manager.Register(HandlerMap.Kick, func(rpcCtx *context.RPCCtx, data struct {
		UID  string
		Data []byte
	}) {
		connection := GetConnection(data.UID)
		if connection == nil {
			return
		}
		connection.Kick(data.Data)
	})
}
