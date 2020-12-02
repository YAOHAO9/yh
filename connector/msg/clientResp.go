package msg

import "encoding/json"

// ClientResp client message response
type ClientResp struct {
	RequestID int32
	Code      int
	Data      interface{} `json:",omitempty"`
}

// ToBytes To []byte
func (m ClientResp) ToBytes() (bytes []byte) {
	bytes, _ = json.Marshal(m)
	return
}
