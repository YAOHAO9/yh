package application

import (
	"github.com/YAOHAO9/yh/rpc/client/clientmanager"
	"github.com/YAOHAO9/yh/rpc/message"
	"github.com/YAOHAO9/yh/rpc/router"
)

type rpcManager struct {
}

// ToServer Rpc到指定的Server
func (rpc rpcManager) ToServer(serverID string, session *message.Session, msg *message.ClientMessage, f func(rpcResp *message.RPCResp)) {

	rpcClient := clientmanager.GetClientByID(serverID)
	if f == nil {
		rpcClient.SendRPCNotify(session, msg)
	} else {
		rpcClient.SendRPCRequest(session, msg, f)
	}
}

// ToServer Rpc到指定的Server
func (rpc rpcManager) ByKind(kind string, session *message.Session, msg *message.ClientMessage, f func(rpcResp *message.RPCResp)) {

	// 根据类型转发
	rpcClient := clientmanager.GetClientByRouter(router.Info{
		ServerKind: kind,
		Handler:    msg.Handler,
		Session:    *session,
	})

	if f == nil {
		rpcClient.SendRPCNotify(session, msg)
	} else {
		rpcClient.SendRPCRequest(session, msg, f)
	}
}

type notify struct{}

// ToServer Rpc到指定的Server
func (n notify) ToServer(serverID string, session *message.Session, handler string, data interface{}) {

	rpcClient := clientmanager.GetClientByID(serverID)
	rpcClient.SendRPCNotify(session, &message.ClientMessage{
		Handler: handler,
		Data:    data,
	})
}

// ByKind Rpc到指定的Server
func (n notify) ByKind(kind string, session *message.Session, handler string, data interface{}) {

	// 根据类型转发
	rpcClient := clientmanager.GetClientByRouter(router.Info{
		ServerKind: kind,
		Handler:    handler,
		Session:    *session,
	})

	rpcClient.SendRPCNotify(session, &message.ClientMessage{
		Handler: handler,
		Data:    data,
	})
}

type request struct{}

// ToServer Rpc到指定的Server
func (req request) ToServer(serverID string, session *message.Session, handler string, data interface{}, f func(rpcResp *message.RPCResp)) {
	rpcClient := clientmanager.GetClientByID(serverID)
	rpcClient.SendRPCRequest(session, &message.ClientMessage{
		Handler: handler,
		Data:    data,
	}, f)
}

// ByKind Rpc到指定的Server
func (req request) ByKind(kind string, session *message.Session, handler string, data interface{}, f func(rpcResp *message.RPCResp)) {

	// 根据类型转发
	rpcClient := clientmanager.GetClientByRouter(router.Info{
		ServerKind: kind,
		Handler:    handler,
		Session:    *session,
	})

	rpcClient.SendRPCRequest(session, &message.ClientMessage{
		Handler: handler,
		Data:    data,
	}, f)
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
