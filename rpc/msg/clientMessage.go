package msg

import "encoding/json"

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
