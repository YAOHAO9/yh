package application

import (
	"fmt"

	"github.com/YAOHAO9/yh/rpc/client/clientmanager"
	"github.com/YAOHAO9/yh/rpc/message"
	"github.com/YAOHAO9/yh/rpc/session"
)

type notify struct{}

// ToServer Rpc到指定的Server
func (n notify) ToServer(serverID string, session *session.Session, handler string, data interface{}) {

	rpcClient := clientmanager.GetClientByID(serverID)
	if rpcClient == nil {
		fmt.Println("Rpc 消息发送失败，没有找到对应的服务器")
		return
	}

	rpcMsg := &message.RPCMessage{
		Kind:    message.KindEnum.RPC,
		Handler: handler,
		Data:    data,
		Session: session,
	}

	rpcClient.SendRPCNotify(session, rpcMsg)
}

// ByKind Rpc到指定的Server
func (n notify) ByKind(serverKind string, session *session.Session, handler string, data interface{}) {
	rpcMsg := &message.RPCMessage{
		Kind:    message.KindEnum.RPC,
		Handler: handler,
		Data:    data,
		Session: session,
	}

	// 根据类型转发
	rpcClient := clientmanager.GetClientByRouter(serverKind, rpcMsg)
	if rpcClient == nil {
		fmt.Println("Rpc 消息发送失败，没有找到对应的服务器")
		return
	}
	rpcClient.SendRPCNotify(session, rpcMsg)
}

type request struct{}

// ToServer Rpc到指定的Server
func (req request) ToServer(serverID string, session *session.Session, handler string, data interface{}, f func(rpcResp *message.RPCResp)) {
	rpcClient := clientmanager.GetClientByID(serverID)
	if rpcClient == nil {
		fmt.Println("Rpc 消息发送失败，没有找到对应的服务器")
		return
	}
	rpcMsg := &message.RPCMessage{
		Kind:    message.KindEnum.RPC,
		Handler: handler,
		Data:    data,
		Session: session,
	}
	rpcClient.SendRPCRequest(session, rpcMsg, f)
}

// ByKind Rpc到指定的Server
func (req request) ByKind(serverKind string, session *session.Session, handler string, data interface{}, f func(rpcResp *message.RPCResp)) {
	rpcMsg := &message.RPCMessage{
		Kind:    message.KindEnum.RPC,
		Handler: handler,
		Data:    data,
		Session: session,
	}
	// 根据类型转发
	rpcClient := clientmanager.GetClientByRouter(serverKind, rpcMsg)
	if rpcClient == nil {
		fmt.Println("Rpc 消息发送失败，没有找到对应的服务器")
		return
	}
	rpcClient.SendRPCRequest(session, rpcMsg, f)
}

// RPC 实例
var RPC struct {
	Notify  notify
	Request request
}
