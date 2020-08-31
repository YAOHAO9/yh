package application

import (
	"github.com/YAOHAO9/yh/rpc/client"
	"github.com/YAOHAO9/yh/rpc/message"
	"github.com/YAOHAO9/yh/rpc/router"
)

// RegisterRouter 注册路由
func (app Application) RegisterRouter(serverKind string, route func(rpcMsg *message.RPCMsg, clients []*client.RPCClient) *client.RPCClient) {
	router.Manager.Register(serverKind, route)
}
