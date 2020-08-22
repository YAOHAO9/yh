package syshandler

import (
	"github.com/YAOHAO9/yh/rpc/handler"
)

// Handler SysHandler
type Handler struct {
	*handler.BaseHandler
}

// Manager return SysHandler
var Manager = &Handler{&handler.BaseHandler{Map: make(handler.Map)}}
