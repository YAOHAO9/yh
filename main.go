package main

import (
	"errors"
	"fmt"
	_ "net/http/pprof"

	"github.com/YAOHAO9/pine/application"
	"github.com/YAOHAO9/pine/application/config"
	"github.com/YAOHAO9/pine/channel/channelfactory"
	"github.com/YAOHAO9/pine/rpc/client"
	"github.com/YAOHAO9/pine/rpc/context"
	"github.com/YAOHAO9/pine/rpc/handler"
	"github.com/YAOHAO9/pine/rpc/message"
)

func main() {

	// slice
	slice := []interface{}{"hello", "world", "hello", "everyone!"}
	for index, value := range slice {
		fmt.Printf("slice %d is: %s\n", index, value)
	}

	app := application.CreateApp()

	app.AsConnector(func(uid string, token string, sessionData map[string]interface{}) error {

		if uid == "" || token == "" {
			return errors.New("Invalid token")
		}
		sessionData[token] = token

		return nil
	})

	app.RegisteHandler("handler", func(rpcCtx *context.RPCCtx) *handler.Resp {

		channelInstance := channelfactory.CreateChannel("101")
		channelInstance.Add(rpcCtx.Session.UID, rpcCtx.Session)

		channelInstance.PushMessage("onMsg", "PushMessage")                                               // 推送给所有在当前channel中的玩家
		channelInstance.PushMessageToOthers([]string{rpcCtx.Session.UID}, "onMsg", "PushMessageToOthers") // 推送给除了切片内的channel中的玩家
		channelInstance.PushMessageToUser(rpcCtx.Session.UID, "onMsg", "PushMessageToUser")               // 只推送给当前玩家
		channelInstance.PushMessageToUsers([]string{rpcCtx.Session.UID}, "onMsg", "PushMessageToUsers")   // 只推送给切片的指定的玩家
		return nil
	})

	app.RegisteRemoter("rpc", func(rpcCtx *context.RPCCtx) *handler.Resp {

		return &handler.Resp{
			Data: config.GetServerConfig().ID + ": 收到Rpc消息",
		}
	})

	app.RegisteHandlerBeforeFilter(func(rpcCtx *context.RPCCtx) (next bool) {
		// logrus.Info("BeforeFilter => requestId: ", rpcCtx.GetRequestID())
		return true
	})

	app.RegisteHandlerAfterFilter(func(rpcResp *message.RPCResp) (next bool) {
		// logrus.Info("AfterFilter hangler:", rpcResp.Handler, " RequestId: ", rpcResp.RequestID)
		return true
	})

	app.RegisteRouter("ddz", func(rpcMsg *message.RPCMsg, clients []*client.RPCClient) *client.RPCClient {

		for _, clientInfo := range clients {
			if chatServerID, ok := rpcMsg.Session.Get("chatServerID").(string); ok && clientInfo.ServerConfig.ID == chatServerID {
				return clientInfo
			}
		}

		return nil //pine will get one rpc client by random
	})

	app.Start()
}
