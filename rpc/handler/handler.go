package handler

import (
	"fmt"

	"github.com/YAOHAO9/yh/rpc/context"
	"github.com/YAOHAO9/yh/rpc/message"
	"github.com/sirupsen/logrus"
)

// Resp handler返回值
type Resp struct {
	Code int
	Data interface{}
}

// Map handler函数仓库
type Map map[string]func(rpcCtx *context.RPCCtx) (resp *Resp)

// Handler Handler
type Handler struct {
	Map Map
}

// Register handler
func (handler Handler) Register(name string, f func(rpcCtx *context.RPCCtx) (resp *Resp)) {
	handler.Map[name] = f
}

// Exec 执行handler
func (handler Handler) Exec(rpcCtx *context.RPCCtx) {

	f, ok := handler.Map[rpcCtx.GetHandler()]
	if ok {
		go func() {

			defer func() {
				// 错误处理
				if err := recover(); err != nil {
					if entry, ok := err.(*logrus.Entry); ok {
						err, _ := (&logrus.JSONFormatter{}).Format(entry)
						rpcCtx.SendMsg(string(err), message.StatusCode.Fail)
						return
					}
					rpcCtx.SendMsg(err, message.StatusCode.Fail)
				}
			}()
			// 执行handler
			resp := f(rpcCtx)

			// 回复消息
			if resp == nil {
				rpcCtx.SendMsg(nil, message.StatusCode.Successful)
				return
			}

			// 回复消息
			rpcCtx.SendMsg(resp.Data, resp.Code)
		}()
	} else {
		rpcCtx.SendMsg(fmt.Sprintf("SysHandler %v 不存在", rpcCtx.GetHandler()), message.StatusCode.Fail)
	}
}

// Manager return RPCHandler
var Manager = &Handler{Map: make(Map)}
