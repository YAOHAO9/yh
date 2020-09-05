package connector

import "encoding/json"

// ClientMsgResp client message response
type ClientMsgResp struct {
	RequestID int
	Code      int
	Data      interface{} `json:",omitempty"`
}

// ToBytes To []byte
func (m ClientMsgResp) ToBytes() (data []byte) {
	data, _ = json.Marshal(m)
	return
}
