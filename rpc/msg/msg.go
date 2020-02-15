package msg

import "encoding/json"

// ========================================================
// Message
// ========================================================

// Message 客户端发过来的消息的基本格式
type Message struct {
	ServerID string
	Kind     string
	Handler  string
	Index    int
	Data     interface{}
}

// ToBytes To []byte
func (m Message) ToBytes() (data []byte) {
	data, _ = json.Marshal(m)
	return
}

// ========================================================
// Session
// ========================================================

// Session of connection
type Session struct {
	uid  string
	data map[string]interface{}
}

// ========================================================
// ForwardMessage
// ========================================================

// ForwardMessage 转发消息结构
type ForwardMessage struct {
	IsRPC bool
	Msg   *Message
	Session
}

// ToBytes To []byte
func (m ForwardMessage) ToBytes() (data []byte) {
	data, _ = json.Marshal(m)
	return
}

// ========================================================
// ResponseMessage
// ========================================================

// ResponseMessage 服务端推送的消息
type ResponseMessage struct {
	Index int
	Code  int
	Event string
	Data  interface{}
}

// ToBytes To []byte
func (m ResponseMessage) ToBytes() (data []byte) {
	data, _ = json.Marshal(m)
	return
}

// ========================================================
// MessageCode
// ========================================================

// MessageCode 消息状态码
type MessageCode struct {
	Successful int
	Fail       int
}

// ========================================================
// StatusCode
// ========================================================

var messageCode MessageCode

// StatusCode 消息状态码
func StatusCode() MessageCode {
	return messageCode
}

func init() {
	messageCode = MessageCode{
		Successful: 0,
		Fail:       200000002,
	}
}
