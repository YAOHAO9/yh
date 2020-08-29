package application

import (
	"github.com/YAOHAO9/yh/rpc/client/clientmanager"
	"github.com/YAOHAO9/yh/rpc/message"
	"github.com/YAOHAO9/yh/rpc/router"
)

type notify struct{}

// ToServer Rpc到指定的Server
func (n notify) ToServer(serverID string, session *message.Session, handler string, data interface{}) {

	rpcClient := clientmanager.GetClientByID(serverID)

	rpcMsg := &message.RPCMessage{
		Kind:    message.KindEnum.RPC,
		Handler: handler,
		Data:    data,
		Session: session,
	}

	rpcClient.SendRPCNotify(session, rpcMsg)
}

// ByKind Rpc到指定的Server
func (n notify) ByKind(kind string, session *message.Session, handler string, data interface{}) {

	// 根据类型转发
	rpcClient := clientmanager.GetClientByRouter(router.Info{
		ServerKind: kind,
		Handler:    handler,
		Session:    *session,
	})

	rpcMsg := &message.RPCMessage{
		Kind:    message.KindEnum.RPC,
		Handler: handler,
		Data:    data,
		Session: session,
	}

	rpcClient.SendRPCNotify(session, rpcMsg)
}

type request struct{}

// ToServer Rpc到指定的Server
func (req request) ToServer(serverID string, session *message.Session, handler string, data interface{}, f func(rpcResp *message.RPCResp)) {
	rpcClient := clientmanager.GetClientByID(serverID)
	rpcMsg := &message.RPCMessage{
		Kind:    message.KindEnum.RPC,
		Handler: handler,
		Data:    data,
		Session: session,
	}
	rpcClient.SendRPCRequest(session, rpcMsg, f)
}

// ByKind Rpc到指定的Server
func (req request) ByKind(kind string, session *message.Session, handler string, data interface{}, f func(rpcResp *message.RPCResp)) {

	// 根据类型转发
	rpcClient := clientmanager.GetClientByRouter(router.Info{
		ServerKind: kind,
		Handler:    handler,
		Session:    *session,
	})
	rpcMsg := &message.RPCMessage{
		Kind:    message.KindEnum.RPC,
		Handler: handler,
		Data:    data,
		Session: session,
	}
	rpcClient.SendRPCRequest(session, rpcMsg, f)
}

type rpc struct {
	Notify  notify
	Request request
}

// RPC 实例
var RPC = rpc{
	Notify:  notify{},
	Request: request{},
}
