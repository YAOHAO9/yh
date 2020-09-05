package handler

import (
	"fmt"
	"runtime/debug"

	"github.com/YAOHAO9/yh/rpc/context"
	"github.com/YAOHAO9/yh/rpc/message"
)

// Resp handler返回值
type Resp struct {
	Code int
	Data interface{}
}

// Map handler函数仓库
type Map map[string]func(rpcCtx *context.RPCCtx) (resp *Resp)

// BaseHandler BaseHandler
type BaseHandler struct {
	Map Map
}

// Register handler
func (handler BaseHandler) Register(name string, f func(rpcCtx *context.RPCCtx) (resp *Resp)) {
	handler.Map[name] = f
}

// Exec 执行handler
func (handler BaseHandler) Exec(rpcCtx *context.RPCCtx) {

	f, ok := handler.Map[rpcCtx.GetHandler()]
	if ok {
		go func() {

			defer func() {
				if err := recover(); err != nil {
					debug.PrintStack()
					rpcCtx.SendMsg(err, message.StatusCode.Fail)
				}
			}()

			resp := f(rpcCtx)

			if resp == nil {
				rpcCtx.SendMsg(nil, message.StatusCode.Successful)
				return
			}

			rpcCtx.SendMsg(resp.Data, resp.Code)

		}()
	} else {
		rpcCtx.SendMsg(fmt.Sprintf("SysHandler %v 不存在", rpcCtx.GetHandler()), message.StatusCode.Fail)
	}
}

var a = BaseHandler{make(Map)}
