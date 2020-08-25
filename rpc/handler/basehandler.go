package handler

import (
	"fmt"

	"github.com/YAOHAO9/yh/rpc/response"
)

// Map handler函数仓库
type Map map[string]func(respCtx *response.RespCtx)

// BaseHandler BaseHandler
type BaseHandler struct {
	Map Map
}

// Register handler
func (handler BaseHandler) Register(name string, f func(respCtx *response.RespCtx)) {
	handler.Map[name] = f
}

// Exec 执行handler
func (handler BaseHandler) Exec(respCtx *response.RespCtx) {

	f, ok := handler.Map[respCtx.RPCMsg.Handler]
	if ok {
		f(respCtx)
	} else {
		respCtx.SendFailMessage(fmt.Sprintf("SysHandler %v 不存在", respCtx.RPCMsg.Handler))
	}
}

var a = BaseHandler{make(Map)}
