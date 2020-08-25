package message

import "encoding/json"

// RPCMessage 转发消息结构
type RPCMessage struct {
	Kind      int // message kind
	RequestID int `json:",omitempty"`
	Handler   string
	Data      interface{}
	Session   *Session
}

// ToBytes To []byte
func (m RPCMessage) ToBytes() (data []byte) {
	data, _ = json.Marshal(m)
	return
}
