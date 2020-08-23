package application

import (
	"github.com/YAOHAO9/yh/rpc/client/clientmanager"
	"github.com/YAOHAO9/yh/rpc/msg"
	"github.com/YAOHAO9/yh/rpc/router"
)

type rpcManager struct {
}

// ToServer Rpc到指定的Server
func (rpc rpcManager) ToServer(serverID string, session *msg.Session, message *msg.ClientMessage, f func(rpcResp *msg.RPCResp)) {

	rpcClient := clientmanager.GetClientByID(serverID)
	if f == nil {
		rpcClient.SendRPCNotify(session, message)
	} else {
		rpcClient.SendRPCRequest(session, message, f)
	}
}

// ToServer Rpc到指定的Server
func (rpc rpcManager) ByKind(kind string, session *msg.Session, message *msg.ClientMessage, f func(rpcResp *msg.RPCResp)) {

	// 根据类型转发
	rpcClient := clientmanager.GetClientByRouter(router.Info{
		ServerKind: kind,
		Handler:    message.Handler,
		Session:    *session,
	})

	if f == nil {
		rpcClient.SendRPCNotify(session, message)
	} else {
		rpcClient.SendRPCRequest(session, message, f)
	}
}

type notify struct{}

// ToServer Rpc到指定的Server
func (n notify) ToServer(serverID string, session *msg.Session, message *msg.ClientMessage) {
	rpcClient := clientmanager.GetClientByID(serverID)
	rpcClient.SendRPCNotify(session, message)
}

// ByKind Rpc到指定的Server
func (n notify) ByKind(kind string, session *msg.Session, message *msg.ClientMessage) {

	// 根据类型转发
	rpcClient := clientmanager.GetClientByRouter(router.Info{
		ServerKind: kind,
		Handler:    message.Handler,
		Session:    *session,
	})

	rpcClient.SendRPCNotify(session, message)
}

type request struct{}

// ToServer Rpc到指定的Server
func (req request) ToServer(serverID string, session *msg.Session, message *msg.ClientMessage, f func(rpcResp *msg.RPCResp)) {
	rpcClient := clientmanager.GetClientByID(serverID)
	rpcClient.SendRPCRequest(session, message, f)
}

// ByKind Rpc到指定的Server
func (req request) ByKind(kind string, session *msg.Session, message *msg.ClientMessage, f func(rpcResp *msg.RPCResp)) {

	// 根据类型转发
	rpcClient := clientmanager.GetClientByRouter(router.Info{
		ServerKind: kind,
		Handler:    message.Handler,
		Session:    *session,
	})

	rpcClient.SendRPCRequest(session, message, f)
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
