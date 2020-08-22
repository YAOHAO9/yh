package rpchandler

import (
	"yh/rpc/handler"
)

// Handler RPCHandler
type Handler struct {
	*handler.BaseHandler
}

// Manager return RPCHandler
var Manager = &Handler{&handler.BaseHandler{Map: make(handler.Map)}}
