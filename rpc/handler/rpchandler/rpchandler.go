package rpchandler

import (
	"fmt"
	"trial/rpc/handler"
	"trial/rpc/response"
)

// Map of rpc
type Map handler.Map

var rpchandlerMap = make(handler.Map)

// Manager return rpchandlerMap
func Manager() handler.Map {
	return rpchandlerMap
}

// Exec 执行handler
func (rpchandlerMap Map) Exec(respCtx *response.RespCtx) {

	f, ok := rpchandlerMap[respCtx.Fm.Handler]
	if ok {
		f(respCtx)
	} else {
		respCtx.SendFailMessage(fmt.Sprintf("RPCHandler %v 不存在", respCtx.Fm.Handler))
	}
}
