package main

import (
	"math/rand"
	"time"

	"github.com/YAOHAO9/yh/application"
	"github.com/YAOHAO9/yh/application/config"
	"github.com/YAOHAO9/yh/rpc/channel/channelfactory"
	"github.com/YAOHAO9/yh/rpc/client"
	"github.com/YAOHAO9/yh/rpc/msg"
	"github.com/YAOHAO9/yh/rpc/response"
	"github.com/YAOHAO9/yh/rpc/router"
)

func main() {
	app := application.CreateApp()

	app.RegisterHandler("handler", func(respCtx *response.RespCtx) {
		respCtx.SendSuccessfulMessage(config.GetServerConfig().ID + ": 收到Handler消息")
		channel := channelfactory.CreateChannel("test")
		channel.Add(respCtx.RPCMsg.Session.CID, respCtx.RPCMsg.Session)

		go func() {
			for {
				time.Sleep(time.Second * 1)
				channel.PushMessage([]string{respCtx.RPCMsg.Session.CID}, "啊哈哈啊")
			}
		}()
	})

	app.RegisterRPCHandler("rpc", func(respCtx *response.RespCtx) {
		respCtx.SendSuccessfulMessage(config.GetServerConfig().ID + ": 收到Rpc消息")
	})

	app.RegisterRPCAfterFilter(func(rm *msg.RPCResp) (next bool) {
		// rm.RequestID -= 1000
		return true
	})

	app.RegisterHandlerAfterFilter(func(rm *msg.RPCResp) (next bool) {
		// rm.RequestID += 1000
		return true
	})

	app.RegisterRouter("ddz", func(routerInfo router.Info, clients []*client.RPCClient) *client.RPCClient {
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
