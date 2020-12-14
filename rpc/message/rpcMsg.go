package message

import (
	"encoding/json"

	"github.com/YAOHAO9/pine/rpc/session"
)

// RPCMsg 转发消息结构
type RPCMsg struct {
	From      string
	Handler   string
	RawData   []byte `json:",omitempty"`
	Session   *session.Session
	RequestID int32 `json:",omitempty"`
}

// ToBytes To []byte
func (m RPCMsg) ToBytes() (bytes []byte) {
	bytes, _ = json.Marshal(m)
	return
}
