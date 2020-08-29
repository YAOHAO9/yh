package message

import "encoding/json"

// RPCMessage 转发消息结构
type RPCMessage struct {
	Kind      int // message kind
	Handler   string
	Data      interface{}
	Session   *Session
	RequestID int `json:",omitempty"`
}

// ToBytes To []byte
func (m RPCMessage) ToBytes() (data []byte) {
	data, _ = json.Marshal(m)
	return
}
