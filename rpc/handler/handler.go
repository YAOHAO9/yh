package handler

var handlerMap = make(Map)

// Manager return HandlerMap
func Manager() Map {
	return handlerMap
}
