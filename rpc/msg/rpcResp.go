package msg

import "encoding/json"

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
