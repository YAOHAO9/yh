package handler

import (
	"fmt"
	"regexp"
	"time"

	"github.com/YAOHAO9/pine/application/config"
	"github.com/YAOHAO9/pine/rpc/context"
	"github.com/YAOHAO9/pine/rpc/message"
	"github.com/sirupsen/logrus"
)

// Resp handler返回值
type Resp struct {
	Code int
	Data interface{}
}

// Map handler函数仓库
type Map map[string]func(rpcCtx *context.RPCCtx)

// Handler Handler
type Handler struct {
	Map Map
}

// Register handler
func (handler Handler) Register(name string, f func(rpcCtx *context.RPCCtx)) {
	handler.Map[name] = f
}

// Exec 执行handler
func (handler Handler) Exec(rpcCtx *context.RPCCtx) {

	f, ok := handler.Map[rpcCtx.GetHandler()]
	if ok {
		defer func() {
			// 错误处理
			if err := recover(); err != nil {
				if entry, ok := err.(*logrus.Entry); ok {
					err, _ := (&logrus.JSONFormatter{}).Format(entry)
					rpcCtx.SendMsg(fmt.Sprint(err), message.StatusCode.Fail)
					return
				}
				logrus.Error(err)
				rpcCtx.SendMsg(fmt.Sprint(err), message.StatusCode.Fail)
			}
		}()
		go time.AfterFunc(time.Minute, func() {
			if rpcCtx.GetRequestID() > 0 {
				logrus.Error(fmt.Sprintf("(%v.%v) response timeout ", config.GetServerConfig().Kind, rpcCtx.GetHandler()))
				rpcCtx.SendMsg(fmt.Sprintf("(%v.%v) response timeout ", config.GetServerConfig().Kind, rpcCtx.GetHandler()), message.StatusCode.Fail)
			}
		})
		// 执行handler
		f(rpcCtx)
	} else {
		handler := rpcCtx.GetHandler()

		reg, _ := regexp.Compile("^__")

		if reg.MatchString(handler) {

			realHandler := reg.ReplaceAll([]byte(rpcCtx.GetHandler()), []byte(""))
			rpcCtx.SetHandler(string(realHandler))
			rpcCtx.SendMsg(fmt.Sprintf("Handler %v 不存在", rpcCtx.GetHandler()), message.StatusCode.Fail)

		} else {
			rpcCtx.SendMsg(fmt.Sprintf("Remoter %v 不存在", rpcCtx.GetHandler()), message.StatusCode.Fail)
		}

	}
}

// Manager return RPCHandler
var Manager = &Handler{Map: make(Map)}
