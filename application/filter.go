package application

import (
	connector_filter "github.com/YAOHAO9/pine/connector/filter"
	"github.com/YAOHAO9/pine/rpc/message"
)

// RegisteHandlerBeforeFilter 注册before filter
func (app Application) RegisteHandlerBeforeFilter(f func(rpcMsg *message.RPCMsg) (error)) {
	connector_filter.Before.Register(f)
}

// RegisteHandlerAfterFilter 注册after filter
func (app Application) RegisteHandlerAfterFilter(f func(rpcResp *message.PineMsg) (error)) {
	connector_filter.After.Register(f)
}
