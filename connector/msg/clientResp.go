package msg

import "encoding/json"

// ClientResp client message response
type ClientResp struct {
	RequestID int
	Code      int
	Data      interface{} `json:",omitempty"`
}

// ToBytes To []byte
func (m ClientResp) ToBytes() (data []byte) {
	data, _ = json.Marshal(m)
	return
}
