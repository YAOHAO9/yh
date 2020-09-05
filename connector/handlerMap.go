package connector

// HandlerMap 系统PRC枚举
var HandlerMap = struct {
	PushMessage   string
	UpdateSession string
}{
	PushMessage:   "__PushMessage__",
	UpdateSession: "__UpdateSession__",
}
