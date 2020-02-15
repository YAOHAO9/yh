package handler

import (
	"fmt"
	"trial/rpc/msg"
	"trial/rpc/response"

	"github.com/gorilla/websocket"
)

// Map map
type Map map[string]func(conn *websocket.Conn, forwardMessage *msg.ForwardMessage)

// Register handler
func (handlerMap Map) Register(name string, f func(conn *websocket.Conn, forwardMessage *msg.ForwardMessage)) {
	handlerMap[name] = f
}

// Exec 执行handler
func (handlerMap Map) Exec(conn *websocket.Conn, forwardMessage *msg.ForwardMessage) {

	f, ok := handlerMap[forwardMessage.Msg.Handler]
	if ok {
		f(conn, forwardMessage)
	} else {
		response.SendFailMessage(conn, forwardMessage.Msg.Index, fmt.Sprintf("Handler %v 不存在", forwardMessage.Msg.Handler))
	}
}

var handlerMap Map

func init() {
	handlerMap = make(Map)
}

// Manager return HandlerMap
func Manager() Map {
	return handlerMap
}
