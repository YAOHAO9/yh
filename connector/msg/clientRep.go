package msg

import "encoding/json"

// ClientReq 客户端发过来的消息的基本格式
type ClientReq struct {
	Route     string
	RequestID int    `json:",omitempty"`
	Data      []byte `json:",omitempty"`
}

// ToBytes To []byte
func (m ClientReq) ToBytes() (bytes []byte) {
	bytes, _ = json.Marshal(m)
	return
}
