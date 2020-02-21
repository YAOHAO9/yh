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
	UID  string
	data map[string]interface{}
}

// Get a value from session
func (s Session) Get(key string) interface{} {
	return s.data[key]
}

// Set a value to session
func (s Session) Set(key string, v interface{}) {
	s.data[key] = v
}

// ========================================================
// ForwardMessage
// ========================================================

// ForwardMessage 转发消息结构
type ForwardMessage struct {
	Kind    int
	Index   int
	Msg     *Message
	Session *Session
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
	Kind  int
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
