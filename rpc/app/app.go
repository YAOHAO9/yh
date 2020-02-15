package app

import (
	"trial/rpc/handler"
	"trial/rpc/msg"
	"trial/rpc/rpchandler"

	"github.com/gorilla/websocket"
)

// RegisterHandler 注册Handler
func RegisterHandler(name string, f func(conn *websocket.Conn, forwardMessage *msg.ForwardMessage)) {
	handler.Manager().Register(name, f)
}

// RegisterRPC 注册Handler
func RegisterRPC(name string, f func(conn *websocket.Conn, forwardMessage *msg.ForwardMessage)) {
	rpchandler.Manager().Register(name, f)
}
