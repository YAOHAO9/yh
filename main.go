package main

import (
	"math/rand"

	_ "net/http/pprof"

	"github.com/YAOHAO9/pine/application"
	"github.com/YAOHAO9/pine/application/config"
	"github.com/YAOHAO9/pine/channel/channelfactory"
	"github.com/YAOHAO9/pine/rpc/client"
	"github.com/YAOHAO9/pine/rpc/context"
	"github.com/YAOHAO9/pine/rpc/handler"
	"github.com/YAOHAO9/pine/rpc/message"
	"github.com/sirupsen/logrus"
)

func main() {
	app := application.CreateApp()

	app.RegisteHandler("haha", func(rpcCtx *context.RPCCtx) *handler.Resp {
		channel := channelfactory.CreateChannel("test") // 创建channel
		channel.Add(rpcCtx.Session.CID, rpcCtx.Session)
		// logrus.Trace("这是一个有意思的log", "啊时代发生的符合")
		// logrus.Debug("这是一个有意思的log", "啊时代发生的符合")
		// logrus.Info("这是一个有意思的log", "啊时代发生的符合")
		// logrus.Warn("这是一个有意思的log", "啊时代发生的符合")

		// // logrus.WithFields(logrus.Fields{
		// // 	"name": "example",
		// // }).Panic("嘿嘿嘿", true)

		// // if rand.Int() < 1 {
		// // 	logrus.Fatal("这是一个有意思的log", "啊时代发生的符合")
		// // } else {
		// // 	logrus.Panic("panic")
		// // }

		// rpcCtx.Session.Set("RequestID", rpcCtx.Data.(map[string]interface{})["RequestID"])
		// application.UpdateSession(rpcCtx.Session, "RequestID")
		// channel.PushMessageToOthers([]string{}, "test", "哈哈哈哈哈")
		return &handler.Resp{
			Data: "asldkfasdklfjs",
		}
	})

	app.RegisteRemoter("rpc", func(rpcCtx *context.RPCCtx) *handler.Resp {

		return &handler.Resp{
			Data: config.GetServerConfig().ID + ": 收到Rpc消息",
		}
	})

	app.RegisteHandlerBeforeFilter(func(rpcCtx *context.RPCCtx) (next bool) {
		logrus.Info("BeforeFilter => requestId: ", rpcCtx.GetRequestID())
		return true
	})

	app.RegisteHandlerAfterFilter(func(rpcResp *message.RPCResp) (next bool) {
		logrus.Info("AfterFilter hangler:", rpcResp.Handler, " RequestId: ", rpcResp.RequestID)
		return true
	})

	app.RegisteRouter("ddz", func(rpcMsg *message.RPCMsg, clients []*client.RPCClient) *client.RPCClient {
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
