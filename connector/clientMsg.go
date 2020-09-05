package connector

import "encoding/json"

// ClientMsg 客户端发过来的消息的基本格式
type ClientMsg struct {
	Route     string
	RequestID int         `json:",omitempty"`
	Data      interface{} `json:",omitempty"`
}

// ToBytes To []byte
func (m ClientMsg) ToBytes() (data []byte) {
	data, _ = json.Marshal(m)
	return
}
