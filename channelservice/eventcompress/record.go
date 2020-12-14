package eventcompress

var eventToCode = make(map[string]byte)

var codeToEvent = make(map[byte]string)

var events = make([]string, 0, 10)

// AddEventRecord 添加需要压缩的客户端监听的事件
func AddEventRecord(eventName string) {
	if _, exist := eventToCode[eventName]; !exist {
		code := byte(len(eventToCode) + 1)
		eventToCode[eventName] = code
		codeToEvent[code] = eventName
		events = append(events, eventName)
	}
}

// GetEventByCode 获取真实Event
// func GetEventByCode(code byte) string {
// 	if value, exist := codeToEvent[code]; exist {
// 		return value
// 	}
// 	return ""
// }

// // GetCodeByEvent 获取真实Event对应的Code
func GetCodeByEvent(eventName string) byte {
	if value, exist := eventToCode[eventName]; exist {
		return value
	}
	return 0
}

// GetEvents 获取Events切片
func GetEvents() []string {
	return events
}
