package application

import (
	"github.com/YAOHAO9/pine/rpc/handler"
	"github.com/YAOHAO9/pine/rpc/handler/remoter"
)

// RegisteHandler 注册Handler
func (app Application) RegisteHandler(name string, f interface{}) {
	handler.Manager.Register(name, f)
}

// RegisteRemoter 注册RPC Handler
func (app Application) RegisteRemoter(name string, f interface{}) {
	remoter.Manager.Register(name, f)
}
