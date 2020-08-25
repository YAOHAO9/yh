package message

import "encoding/json"

// RPCMessage 转发消息结构
type RPCMessage struct {
	Kind      int `json:",omitempty"` // message kind
	RequestID int `json:",omitempty"`
	Handler   string
	Data      interface{}
	Session   *Session `json:",omitempty"`
}

// ToBytes To []byte
func (m RPCMessage) ToBytes() (data []byte) {
	data, _ = json.Marshal(m)
	return
}
