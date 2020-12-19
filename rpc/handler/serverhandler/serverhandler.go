package serverhandler

import (
	"github.com/YAOHAO9/pine/rpc/handler"
)

// ServerHandler ServerHandler
type ServerHandler struct {
	*handler.Handler
}

// Manager return RPCHandler
var Manager = &ServerHandler{
	Handler: &handler.Handler{
		Map: make(handler.Map),
	},
}
