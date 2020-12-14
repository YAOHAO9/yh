package routercompress

var handlerToCode = make(map[string]byte)

var codeToHandler = make(map[byte]string)

var handlers = make([]string, 0, 10)

// AddHandlerRecord 添加记录
func AddHandlerRecord(handlerName string) {
	if _, exist := handlerToCode[handlerName]; !exist {
		code := byte(len(handlerToCode) + 1)
		handlerToCode[handlerName] = code
		codeToHandler[code] = handlerName
		handlers = append(handlers, handlerName)
	}
}

// GetHandlerByCode 获取真实Handler
func GetHandlerByCode(code byte) string {
	if value, exist := codeToHandler[code]; exist {
		return value
	}
	return ""
}

// GetCodeByHandler 获取真实Handler对应的Code
func GetCodeByHandler(handlerName string) byte {
	if value, exist := handlerToCode[handlerName]; exist {
		return value
	}
	return 0
}

// GetHandlers 获取Handlers切片
func GetHandlers() []string {
	return handlers
}

// AddEventRecord 添加需要压缩的客户端监听的事件
func AddEventRecord(eventName string) {
	if _, exist := handlerToCode[eventName]; !exist {
		code := byte(len(handlerToCode) + 1)
		handlerToCode[eventName] = code
		codeToHandler[code] = eventName
		handlers = append(handlers, eventName)
	}
}
