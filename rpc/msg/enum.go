package msg

// KindEnum 消息类型枚举
var KindEnum = struct {
	Sys             int
	Handler         int
	RPC             int
	SysResponse     int
	HandlerResponse int
	RPCResponse     int
}{
	Sys:             1,     // 系统rpc
	Handler:         2,     // 客户端调用 Handler
	RPC:             3,     //  RPC
	SysResponse:     10001, // 系统消息的response
	HandlerResponse: 10002, // handler request的response
	RPCResponse:     10003, // rpc 的response
}

// TypeEnum 消息类型枚举
var TypeEnum = struct {
	// TextMessage denotes a text data me  ssage. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage int
	// BinaryMessage denotes a binary data message.
	BinaryMessage int
	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage int
	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage int
	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage int
}{
	TextMessage:   1,
	BinaryMessage: 2,
	CloseMessage:  8,
	PingMessage:   9,
	PongMessage:   10,
}

// StatusCode 消息状态码
var StatusCode = struct {
	Successful int
	Fail       int
}{
	Successful: 0,
	Fail:       200000002,
}
