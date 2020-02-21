package rpchandler

import (
	"fmt"
	"trial/rpc/handler"
	"trial/rpc/msg"
	"trial/rpc/msg/msgkind"
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
func (rpchandlerMap Map) Exec(respConn *websocket.Conn, fm *msg.ForwardMessage) {

	f, ok := rpchandlerMap[fm.Msg.Handler]
	if ok {
		f(respConn, fm)
	} else {
		response.SendFailMessage(respConn, msgkind.RPC, fm.Msg.Index, fmt.Sprintf("RPCHandler %v 不存在", fm.Msg.Handler))
	}
}
