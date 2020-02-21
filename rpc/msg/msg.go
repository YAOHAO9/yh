package msg

import "encoding/json"

// ========================================================
// Message
// ========================================================

// ClientMessage 客户端发过来的消息的基本格式
type ClientMessage struct {
	ServerID string `json:",omitempty"`
	Kind     string `json:",omitempty"` // server kind
	Handler  string
	Index    int `json:",omitempty"`
	Data     interface{}
}

// ToBytes To []byte
func (m ClientMessage) ToBytes() (data []byte) {
	data, _ = json.Marshal(m)
	return
}

// ClientResp client message response
type ClientResp struct {
	Index   int    `json:",omitempty"`
	Handler string `json:",omitempty"`
	Code    int    `json:",omitempty"`
	Data    interface{}
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
	Data map[string]interface{}
}

// Get a value from session
func (s Session) Get(key string) interface{} {
	return s.Data[key]
}

// Set a value to session
func (s Session) Set(key string, v interface{}) {
	s.Data[key] = v
}

// ========================================================
// ForwardMessage
// ========================================================

// RPCMessage 转发消息结构
type RPCMessage struct {
	Kind    int `json:",omitempty"` // message kind
	Index   int `json:",omitempty"`
	Handler string
	Data    interface{}
	Session *Session `json:",omitempty"`
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
	Kind  int `json:",omitempty"` // response kind
	Index int `json:",omitempty"`
	Code  int
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
