package syshandler

import (
	"fmt"
	"trial/rpc/handler"
	"trial/rpc/response"
)

// Map of rpc
type Map handler.Map

var syshandlerMap = make(handler.Map)

// Manager return syshandlerMap
func Manager() handler.Map {
	return syshandlerMap
}

// Exec 执行handler
func (syshandlerMap Map) Exec(respCtx *response.RespCtx) {

	f, ok := syshandlerMap[respCtx.Fm.Handler]
	if ok {
		f(respCtx)
	} else {
		respCtx.SendFailMessage(fmt.Sprintf("SysHandler %v 不存在", respCtx.Fm.Handler))
	}
}
