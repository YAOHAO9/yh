package main

import (
	"fmt"
	"math/rand"

	"github.com/YAOHAO9/yh/application"
	"github.com/YAOHAO9/yh/application/config"
	"github.com/YAOHAO9/yh/channel/channelfactory"
	"github.com/YAOHAO9/yh/rpc/client"
	"github.com/YAOHAO9/yh/rpc/context"
	"github.com/YAOHAO9/yh/rpc/handler"
	"github.com/YAOHAO9/yh/rpc/message"
)

func main() {
	app := application.CreateApp()

	app.RegisterHandler("handler", func(rpcCtx *context.RPCCtx) *handler.Resp {
		channel := channelfactory.CreateChannel("test") // 创建channel
		channel.Add(rpcCtx.Session.CID, rpcCtx.Session)

		fmt.Println("RequestID of session", rpcCtx.Session.Get("RequestID"), "RequestID", rpcCtx.Data.(map[string]interface{})["RequestID"])
		rpcCtx.Session.Set("RequestID", rpcCtx.Data.(map[string]interface{})["RequestID"])
		application.UpdateSession(rpcCtx.Session, "RequestID")
		channel.PushMessageToOthers([]string{}, "test", "哈哈哈哈哈")
		return nil
	})

	app.RegisterRPCHandler("rpc", func(rpcCtx *context.RPCCtx) *handler.Resp {

		return &handler.Resp{
			Data: config.GetServerConfig().ID + ": 收到Rpc消息",
		}
	})

	app.RegisterRPCAfterFilter(func(rpcResp *message.RPCResp) (next bool) {
		return true
	})

	app.RegisterHandlerAfterFilter(func(rpcResp *message.RPCResp) (next bool) {
		return true
	})

	app.RegisterRouter("ddz", func(rpcMsg *message.RPCMsg, clients []*client.RPCClient) *client.RPCClient {
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
