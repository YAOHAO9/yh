package clienthandler

import (
	"github.com/YAOHAO9/yh/rpc/handler"
)

// Handler ClientHandler
type Handler struct {
	*handler.BaseHandler
}

// Manager return ClientHandler
var Manager = &Handler{&handler.BaseHandler{Map: make(handler.Map)}}
