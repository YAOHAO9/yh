package handler

import (
	"fmt"
	"trial/rpc/response"
)

// Map map
type Map map[string]func(respCtx *response.RespCtx)

// Register handler
func (handlerMap Map) Register(name string, f func(respCtx *response.RespCtx)) {
	handlerMap[name] = f
}

// Exec 执行handler
func (handlerMap Map) Exec(respCtx *response.RespCtx) {

	f, ok := handlerMap[respCtx.Fm.Handler]
	if ok {
		f(respCtx)
	} else {
		respCtx.SendFailMessage(fmt.Sprintf("Handler %v 不存在", respCtx.Fm.Handler))
	}
}
