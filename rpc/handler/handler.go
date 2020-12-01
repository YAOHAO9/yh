package handler

import (
	"encoding/json"
	"fmt"
	"reflect"
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
type Map map[string]interface{}

// Handler Handler
type Handler struct {
	Map Map
}

// Register handler
func (handler Handler) Register(handlerName string, handlerFunc interface{}) {
	handlerType := reflect.TypeOf(handlerFunc)
	if handlerType.Kind() != reflect.Func {
		logrus.Panic("handler(" + handlerName + ")只能为函数")
		return
	}

	handlerValue := reflect.TypeOf(handlerFunc)

	if handlerValue.NumIn() != 2 {
		logrus.Panic("handler(" + handlerName + ")参数只能两个")
		return
	}

	if handlerType.In(0) != reflect.TypeOf(&context.RPCCtx{}) {
		logrus.Panic("handler(" + handlerName + ")第一个参数必须为*context.RPCCtx类型")
		return
	}

	handler.Map[handlerName] = handlerFunc
}

// Exec 执行handler
func (handler Handler) Exec(rpcCtx *context.RPCCtx) {

	handlerInterface, ok := handler.Map[rpcCtx.GetHandler()]

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

		dataInterface := reflect.New(reflect.TypeOf(handlerInterface).In(1)).Interface()

		bdata, e := json.Marshal(rpcCtx.Data)
		if e != nil {
			logrus.Panic(e)
		}

		json.Unmarshal(bdata, dataInterface)

		// 执行handler
		reflect.ValueOf(handlerInterface).Call([]reflect.Value{
			reflect.ValueOf(rpcCtx),
			reflect.ValueOf(dataInterface).Elem(),
		})

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
