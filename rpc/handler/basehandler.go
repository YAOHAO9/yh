package handler

import (
	"fmt"

	"github.com/YAOHAO9/yh/rpc/context"
	"github.com/YAOHAO9/yh/rpc/message"
)

// Map handler函数仓库
type Map map[string]func(rpcCtx *context.RPCCtx)

// BaseHandler BaseHandler
type BaseHandler struct {
	Map Map
}

// Register handler
func (handler BaseHandler) Register(name string, f func(rpcCtx *context.RPCCtx)) {
	handler.Map[name] = f
}

// Exec 执行handler
func (handler BaseHandler) Exec(rpcCtx *context.RPCCtx) {

	f, ok := handler.Map[rpcCtx.GetHandler()]
	if ok {
		go func() {
			f(rpcCtx)
			if rpcCtx.GetRequestID() != -1 && rpcCtx.GetRequestID() != 0 {
				fmt.Println("No response for request", rpcCtx.ToString())
			}
		}()
	} else {
		rpcCtx.SendMsg(fmt.Sprintf("SysHandler %v 不存在", rpcCtx.GetHandler()), message.StatusCode.Fail)
	}
}

var a = BaseHandler{make(Map)}
