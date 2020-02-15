package rpchandler

import "trial/rpc/handler"

var rpchandlerMap = make(handler.Map)

// Manager return HandlerMap
func Manager() handler.Map {
	return rpchandlerMap
}
