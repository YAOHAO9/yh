package handler

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/YAOHAO9/pine/application/config"
	"github.com/YAOHAO9/pine/rpc/context"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
)

// Map handler函数仓库
type Map map[string]interface{}

// Handler Handler
type Handler struct {
	Map Map
}

// Register handler
func (handler *Handler) Register(handlerName string, handlerFunc interface{}) {

	handlerType := reflect.TypeOf(handlerFunc)
	if handlerType.Kind() != reflect.Func {
		logrus.Panic("handler(" + handlerName + ")只能为函数")
		return
	}

	handlerValue := reflect.TypeOf(handlerFunc)

	if handlerValue.NumIn() < 1 {
		logrus.Panic("handler(" + handlerName + ")参数不能少于1个")
		return
	}

	if handlerType.In(0) != reflect.TypeOf(&context.RPCCtx{}) {
		logrus.Panic("handler(" + handlerName + ")第一个参数必须为*context.RPCCtx类型")
		return
	}

	if handlerValue.NumIn() > 2 {
		logrus.Panic("handler(" + handlerName + ")参数不能多于2个")
		return
	}

	handler.Map[handlerName] = handlerFunc
}

// Exec 执行handler
func (handler *Handler) Exec(rpcCtx *context.RPCCtx) (exist bool) {

	defer func() {
		// 错误处理
		if err := recover(); err != nil {
			if entry, ok := err.(*logrus.Entry); ok {
				err, _ := (&logrus.JSONFormatter{}).Format(entry)
				logrus.Error(err, "\nRpcCtx：", rpcCtx.ToString())
				return
			}
			logrus.Error(err, "\nRpcCtx：", rpcCtx.ToString())
		}
	}()

	handlerInterface, exist := handler.Map[rpcCtx.GetHandler()]

	if !exist {
		return
	}

	if rpcCtx.GetRequestID() > 0 {
		go time.AfterFunc(time.Minute, func() {
			if rpcCtx.GetRequestID() != -1 {
				logrus.Error(fmt.Sprintf("(%v.%v) response timeout ", config.GetServerConfig().Kind, rpcCtx.GetHandler()))
			}
		})
	}

	handlerType := reflect.TypeOf(handlerInterface)

	if handlerType.NumIn() == 1 {
		// 执行handler
		reflect.ValueOf(handlerInterface).Call([]reflect.Value{reflect.ValueOf(rpcCtx)})
		return
	}

	paramType := handlerType.In(1)

	var dataInterface interface{}
	if paramType.Kind() == reflect.Ptr {
		dataInterface = reflect.New(paramType.Elem()).Interface()
	} else {
		dataInterface = reflect.New(paramType).Interface()
	}

	msesage, ok := dataInterface.(proto.Message)
	if ok { // protobuf

		proto.Unmarshal(rpcCtx.RawData, msesage)
		var param reflect.Value

		if paramType.Kind() == reflect.Ptr {
			param = reflect.ValueOf(msesage)
		} else {
			param = reflect.ValueOf(msesage).Elem()
		}

		// 执行handler
		reflect.ValueOf(handlerInterface).Call([]reflect.Value{
			reflect.ValueOf(rpcCtx),
			param,
		})
	} else { // json

		dataInterface = reflect.New(paramType).Interface()

		if paramType.Kind() == reflect.Slice && paramType.Elem().Kind() == reflect.Uint8 {
			// 执行handler
			reflect.ValueOf(handlerInterface).Call([]reflect.Value{
				reflect.ValueOf(rpcCtx),
				reflect.ValueOf(rpcCtx.RawData),
			})
			return
		}

		json.Unmarshal(rpcCtx.RawData, dataInterface)

		// 执行handler
		reflect.ValueOf(handlerInterface).Call([]reflect.Value{
			reflect.ValueOf(rpcCtx),
			reflect.ValueOf(dataInterface).Elem(),
		})
	}
	return
}
