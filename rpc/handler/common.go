package handler

import (
	"fmt"
	"trial/rpc/msg"
	"trial/rpc/msg/msgkind"
	"trial/rpc/response"

	"github.com/gorilla/websocket"
)

// Map map
type Map map[string]func(respConn *websocket.Conn, fm *msg.ForwardMessage)

// Register handler
func (handlerMap Map) Register(name string, f func(respConn *websocket.Conn, fm *msg.ForwardMessage)) {
	handlerMap[name] = f
}

// Exec 执行handler
func (handlerMap Map) Exec(respConn *websocket.Conn, fm *msg.ForwardMessage) {

	f, ok := handlerMap[fm.Msg.Handler]
	if ok {
		f(respConn, fm)
	} else {
		response.SendFailMessage(respConn, msgkind.HANDLER, fm.Msg.Index, fmt.Sprintf("Handler %v 不存在", fm.Msg.Handler))
	}
}
