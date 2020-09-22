package application

import (
	"github.com/YAOHAO9/yh/connector/filter"
	"github.com/YAOHAO9/yh/rpc/context"
	"github.com/YAOHAO9/yh/rpc/message"
)

// RegisterHandlerBeforeFilter 注册before filter
func (app Application) RegisterHandlerBeforeFilter(f func(rpcCtx *context.RPCCtx) (next bool)) {
	filter.Before.Register(f)
}

// RegisterHandlerAfterFilter 注册after filter
func (app Application) RegisterHandlerAfterFilter(f func(rpcResp *message.RPCResp) (next bool)) {
	filter.After.Register(f)
}
