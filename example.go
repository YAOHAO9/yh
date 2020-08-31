package main

import (
	"math/rand"

	"github.com/YAOHAO9/yh/application"
	"github.com/YAOHAO9/yh/application/config"
	"github.com/YAOHAO9/yh/rpc/channel/channelfactory"
	"github.com/YAOHAO9/yh/rpc/client"
	"github.com/YAOHAO9/yh/rpc/context"
	"github.com/YAOHAO9/yh/rpc/message"
)

func main() {
	app := application.CreateApp()

	app.RegisterHandler("handler", func(rpcCtx *context.RPCCtx) {
		rpcCtx.SendSuccessfulMessage(config.GetServerConfig().ID + ": 收到Handler消息")
		channel := channelfactory.CreateChannel("test")
		channel.Add(rpcCtx.Session.CID, rpcCtx.Session)

		rpcCtx.Session.Set("RequestID", rpcCtx.GetRequestID())
		application.UpdateSession(rpcCtx.Session, "RequestID")
	})

	app.RegisterRPCHandler("rpc", func(rpcCtx *context.RPCCtx) {
		rpcCtx.SendSuccessfulMessage(config.GetServerConfig().ID + ": 收到Rpc消息")
	})

	app.RegisterRPCAfterFilter(func(rpcResp *message.RPCResp) (next bool) {
		return true
	})

	app.RegisterHandlerAfterFilter(func(rpcResp *message.RPCResp) (next bool) {
		return true
	})

	app.RegisterRouter("ddz", func(rpcMsg *message.RPCMessage, clients []*client.RPCClient) *client.RPCClient {
		var luckClient *client.RPCClient
		for _, clientInfo := range clients {
			if clientInfo.ServerConfig.ID == "ddz-3" {
				luckClient = clientInfo
				break
			}
		}
		if luckClient != nil {
			return luckClient
		}
		return clients[rand.Intn(len(clients))]
	})

	app.Start()
}
