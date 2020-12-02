package message

import "encoding/json"

// RPCResp 服务端推送的消息
type RPCResp struct {
	Handler   string
	RequestID int32
	Code      int
	Data      interface{} `json:",omitempty"`
}

// ToBytes To []byte
func (m RPCResp) ToBytes() (bytes []byte) {
	bytes, _ = json.Marshal(m)
	return
}
