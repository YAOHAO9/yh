package main

import (
	"errors"
	_ "net/http/pprof"

	"github.com/YAOHAO9/pine/application"
	"github.com/YAOHAO9/pine/application/config"
	"github.com/YAOHAO9/pine/rpc/client"
	"github.com/YAOHAO9/pine/rpc/context"
	"github.com/YAOHAO9/pine/rpc/handler"
	"github.com/YAOHAO9/pine/rpc/message"
)

func main() {
	app := application.CreateApp()

	app.AsConnector(func(uid string, token string, sessionData map[string]interface{}) error {

		if uid == "" || token == "" {
			return errors.New("Invalid token")
		}
		sessionData[token] = token

		return nil
	})

	app.RegisteHandler("handler", func(rpcCtx *context.RPCCtx) *handler.Resp {

		return &handler.Resp{
			Data: config.GetServerConfig().ID + ": 收到Handler消息",
		}
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
		return nil
	})

	app.Start()
}
