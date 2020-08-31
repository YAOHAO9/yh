package message

import (
	"encoding/json"

	"github.com/YAOHAO9/yh/rpc/session"
)

// RPCMsg 转发消息结构
type RPCMsg struct {
	Kind      int // message kind
	Handler   string
	Data      interface{}
	Session   *session.Session
	RequestID int `json:",omitempty"`
}

// ToBytes To []byte
func (m RPCMsg) ToBytes() (data []byte) {
	data, _ = json.Marshal(m)
	return
}
