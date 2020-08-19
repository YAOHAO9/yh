package syshandler

import (
	"trial/rpc/handler"
)

// Handler SysHandler
type Handler struct {
	*handler.Handler
}

// Manager return SysHandler
var Manager = &Handler{&handler.Handler{Map: make(handler.Map)}}
