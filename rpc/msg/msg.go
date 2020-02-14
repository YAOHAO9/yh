package msg

import "encoding/json"

// Message 客户端发过来的消息的基本格式
type Message struct {
	ServerID string
	Kind     string
	Handler  string
	Index    int
	Data     interface{}
}

// ResponseMessage 服务端推送的消息
type ResponseMessage struct {
	Index int
	Code  int
	Event string
	Data  interface{}
}

// ToBytes To []byte
func (rm ResponseMessage) ToBytes() (data []byte) {
	data, _ = json.Marshal(rm)
	return
}

// MessageCode 消息状态码
type MessageCode struct {
	Successful int
	Fail       int
}

var messageCode MessageCode

func init() {
	messageCode = MessageCode{
		Successful: 0,
		Fail:       200000002,
	}
}

// StatusCode 消息状态码
func StatusCode() MessageCode {
	return messageCode
}
