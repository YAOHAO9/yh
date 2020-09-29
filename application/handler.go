package application

import (
	"github.com/YAOHAO9/yh/connector"
	"github.com/YAOHAO9/yh/rpc/context"
	"github.com/YAOHAO9/yh/rpc/handler"
	"github.com/YAOHAO9/yh/rpc/handler/rpchandler"
)

// RegisterHandler 注册Handler
func (app Application) RegisterHandler(name string, f func(rpcCtx *context.RPCCtx) *handler.Resp) {
	rpchandler.Manager.Register(connector.HandlerPrefix+name, f)
}

// RegisterRPCHandler 注册RPC Handler
func (app Application) RegisterRPCHandler(name string, f func(rpcCtx *context.RPCCtx) *handler.Resp) {
	rpchandler.Manager.Register(connector.HandlerPrefix+name, f)
}
