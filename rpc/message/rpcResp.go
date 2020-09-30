package message

import "encoding/json"

// RPCResp 服务端推送的消息
type RPCResp struct {
	Handler   string
	RequestID int
	Code      int
	Data      interface{} `json:",omitempty"`
}

// ToBytes To []byte
func (m RPCResp) ToBytes() (data []byte) {
	data, _ = json.Marshal(m)
	return
}
