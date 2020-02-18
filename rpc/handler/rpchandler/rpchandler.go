package rpchandler

import (
	"fmt"
	"trial/rpc/handler"
	"trial/rpc/msg"
	"trial/rpc/response"

	"github.com/gorilla/websocket"
)

// Map of rpc
type Map handler.Map

var rpchandlerMap = make(handler.Map)

// Manager return rpchandlerMap
func Manager() handler.Map {
	return rpchandlerMap
}

// Exec 执行handler
func (rpchandlerMap Map) Exec(conn *websocket.Conn, forwardMessage *msg.ForwardMessage) {

	f, ok := rpchandlerMap[forwardMessage.Msg.Handler]
	if ok {
		f(conn, forwardMessage)
	} else {
		response.SendFailMessage(conn, true, forwardMessage.Msg.Index, fmt.Sprintf("RPCHandler %v 不存在", forwardMessage.Msg.Handler))
	}
}
