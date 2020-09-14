package main

import (
	"math/rand"

	_ "net/http/pprof"

	"github.com/YAOHAO9/yh/application"
	"github.com/YAOHAO9/yh/application/config"
	"github.com/YAOHAO9/yh/channel/channelfactory"
	"github.com/YAOHAO9/yh/rpc/client"
	"github.com/YAOHAO9/yh/rpc/context"
	"github.com/YAOHAO9/yh/rpc/handler"
	"github.com/YAOHAO9/yh/rpc/message"
	"github.com/sirupsen/logrus"
)

func main() {
	app := application.CreateApp()

	app.RegisterHandler("haha", func(rpcCtx *context.RPCCtx) *handler.Resp {
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

	app.RegisterRPCHandler("rpc", func(rpcCtx *context.RPCCtx) *handler.Resp {

		return &handler.Resp{
			Data: config.GetServerConfig().ID + ": 收到Rpc消息",
		}
	})

	app.RegisterHandlerBeforeFilter(func(rpcCtx *context.RPCCtx) (next bool) {
		logrus.Info("BeforeFilter", rpcCtx.Data.(map[string]interface{})["RequestID"].(float64)+100)
		return true
	})

	app.RegisterRPCAfterFilter(func(rpcResp *message.RPCResp) (next bool) {
		logrus.Info("AfterFilter", rpcResp.Data)
		return true
	})

	app.RegisterHandlerAfterFilter(func(rpcResp *message.RPCResp) (next bool) {
		logrus.Info("AfterFilter hangler: ", rpcResp.Handler)
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
