package application

import (
	"github.com/YAOHAO9/yh/rpc/context"
	"github.com/YAOHAO9/yh/rpc/filter/handlerfilter"
	"github.com/YAOHAO9/yh/rpc/filter/rpcfilter"
	"github.com/YAOHAO9/yh/rpc/message"
)

// RegisterHandlerBeforeFilter 注册before filter
func (app Application) RegisterHandlerBeforeFilter(f func(rpcCtx *context.RPCCtx) (next bool)) {
	handlerfilter.Manager.Before.Register(f)
}

// RegisterHandlerAfterFilter 注册after filter
func (app Application) RegisterHandlerAfterFilter(f func(rpcResp *message.RPCResp) (next bool)) {
	handlerfilter.Manager.After.Register(f)
}

// RegisterRPCBeforeFilter 注册before filter of rpc
func (app Application) RegisterRPCBeforeFilter(f func(rpcCtx *context.RPCCtx) (next bool)) {
	rpcfilter.Manager.Before.Register(f)
}

// RegisterRPCAfterFilter 注册after filter of rpc request
func (app Application) RegisterRPCAfterFilter(f func(rpcResp *message.RPCResp) (next bool)) {
	rpcfilter.Manager.After.Register(f)
}
