package msg

import "encoding/json"

// ========================================================
// Message
// ========================================================

// ClientMessage 客户端发过来的消息的基本格式
type ClientMessage struct {
	ServerID string
	Kind     string
	Handler  string
	Index    int
	Data     interface{}
}

// ToBytes To []byte
func (m ClientMessage) ToBytes() (data []byte) {
	data, _ = json.Marshal(m)
	return
}

// ClientResp client message response
type ClientResp struct {
	Index int
	Data  interface{}
}

// ToBytes To []byte
func (m ClientResp) ToBytes() (data []byte) {
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

// RPCMessage 转发消息结构
type RPCMessage struct {
	Kind    int
	Index   int
	Handler string
	Data    interface{}
	Session *Session
}

// ToBytes To []byte
func (m RPCMessage) ToBytes() (data []byte) {
	data, _ = json.Marshal(m)
	return
}

// ========================================================
// ResponseMessage
// ========================================================

// RPCResp 服务端推送的消息
type RPCResp struct {
	Kind  int
	Index int
	Code  int
	Event string
	Data  interface{}
}

// ToBytes To []byte
func (m RPCResp) ToBytes() (data []byte) {
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
