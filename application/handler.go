package application

import (
	"regexp"

	"github.com/YAOHAO9/pine/connector"
	"github.com/YAOHAO9/pine/rpc/context"
	"github.com/YAOHAO9/pine/rpc/handler"
	"github.com/sirupsen/logrus"
)

// RegisteHandler 注册Handler
func (app Application) RegisteHandler(name string, f func(rpcCtx *context.RPCCtx)) {
	reg := regexp.MustCompile(`^__`)
	handler.Manager.Register(connector.HandlerPrefix+name, func(rpcCtx *context.RPCCtx) {
		rpcCtx.SetHandler(string(reg.ReplaceAll([]byte(rpcCtx.GetHandler()), []byte(""))))
		f(rpcCtx)
	})
}

// RegisteRemoter 注册RPC Handler
func (app Application) RegisteRemoter(name string, f func(rpcCtx *context.RPCCtx)) {

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
