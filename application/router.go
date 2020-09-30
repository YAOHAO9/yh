package application

import (
	"github.com/YAOHAO9/pine/rpc/client"
	"github.com/YAOHAO9/pine/rpc/message"
	"github.com/YAOHAO9/pine/rpc/router"
)

// RegisterRouter 注册路由
func (app Application) RegisterRouter(serverKind string, route func(rpcMsg *message.RPCMsg, clients []*client.RPCClient) *client.RPCClient) {
	router.Manager.Register(serverKind, route)
}
