package clienthandler

import (
	"github.com/YAOHAO9/pine/rpc/handler"
	"github.com/YAOHAO9/pine/service/compressservice"
)

// ClientHandler rpc
type ClientHandler struct {
	*handler.Handler
}

// Manager return RPCHandler
var Manager = &ClientHandler{
	Handler: &handler.Handler{
		Map: make(handler.Map),
	},
}

// Register remoter
func (clienthandler *ClientHandler) Register(handlerName string, remoterFunc interface{}) {
	compressservice.Handler.AddHandlerRecord(handlerName)
	clienthandler.Handler.Register(handlerName, remoterFunc)
}
