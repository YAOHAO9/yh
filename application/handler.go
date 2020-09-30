package application

import (
	"regexp"

	"github.com/YAOHAO9/pine/connector"
	"github.com/YAOHAO9/pine/rpc/context"
	"github.com/YAOHAO9/pine/rpc/handler"
	"github.com/sirupsen/logrus"
)

// RegisterHandler 注册Handler
func (app Application) RegisterHandler(name string, f func(rpcCtx *context.RPCCtx) *handler.Resp) {
	reg := regexp.MustCompile(`^__`)
	handler.Manager.Register(connector.HandlerPrefix+name, func(rpcCtx *context.RPCCtx) *handler.Resp {
		rpcCtx.SetHandler(string(reg.ReplaceAll([]byte(rpcCtx.GetHandler()), []byte(""))))
		return f(rpcCtx)
	})
}

// RegisterRemote 注册RPC Handler
func (app Application) RegisterRemote(name string, f func(rpcCtx *context.RPCCtx) *handler.Resp) {

	result, err := regexp.MatchString("^__", name)

	if err != nil {
		logrus.Error(err)
		return
	}

	if result {
		logrus.Error("Remote can not start with __")
		return
	}

	handler.Manager.Register(name, f)
}
