package connector

import "encoding/json"

// ClientMsgResp client message response
type ClientMsgResp struct {
	RequestID int `json:",omitempty"`
	Code      int `json:",omitempty"`
	Data      interface{}
}

// ToBytes To []byte
func (m ClientMsgResp) ToBytes() (data []byte) {
	data, _ = json.Marshal(m)
	return
}
