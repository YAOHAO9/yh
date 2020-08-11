package main

import (
	"trial/rpc/application"
	"trial/rpc/config"
	"trial/rpc/msg"
	"trial/rpc/response"
	RpcServer "trial/rpc/server"
)

func main() {
	app := application.CreateApp()
	register(app)
	// 启动RPC server
	RpcServer.Start()
}

func register(app *application.Application) {

	app.RegisterHandler("handler", func(respCtx *response.RespCtx) {
		respCtx.SendSuccessfulMessage(config.GetServerConfig().ID + ": 收到Handler消息")
	})

	app.RegisterRPCHandler("rpc", func(respCtx *response.RespCtx) {
		respCtx.SendSuccessfulMessage(config.GetServerConfig().ID + ": 收到Rpc消息")
	})

	app.RegisterRPCAfterFilter(func(rm *msg.RPCResp) (next bool) {
		rm.RequestID -= 1000
		return true
	})

	app.RegisterHandlerAfterFilter(func(rm *msg.RPCResp) (next bool) {
		rm.RequestID += 1000
		return true
	})
}
