package msg

import "encoding/json"

// ClientResp client message response
type ClientResp struct {
	RequestID int `json:",omitempty"`
	Code      int `json:",omitempty"`
	Data      interface{}
}

// ToBytes To []byte
func (m ClientResp) ToBytes() (data []byte) {
	data, _ = json.Marshal(m)
	return
}
