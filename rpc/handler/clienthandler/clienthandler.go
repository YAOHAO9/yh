package clienthandler

import (
	"trial/rpc/handler"
)

// Handler ClientHandler
type Handler struct {
	*handler.Handler
}

// Manager return ClientHandler
var Manager = &Handler{&handler.Handler{Map: make(handler.Map)}}
