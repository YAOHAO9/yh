package message

import "encoding/json"

// Notify 主动推送的客户端通知
type Notify struct {
	Route string
	Data  interface{}
}

// ToBytes To []byte
func (m Notify) ToBytes() (data []byte) {
	data, _ = json.Marshal(m)
	return
}
