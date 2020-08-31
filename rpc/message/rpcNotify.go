package message

import "encoding/json"

// RPCNotify 通过RPC主动推送给客户端的通知
type RPCNotify struct {
	Route string
	Data  interface{}
}

// ToBytes To []byte
func (m RPCNotify) ToBytes() (data []byte) {
	data, _ = json.Marshal(m)
	return
}
