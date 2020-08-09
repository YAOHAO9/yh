package response

import (
	"fmt"
	"trial/rpc/msg"

	"github.com/gorilla/websocket"
)

// RespCtx response context
type RespCtx struct {
	Conn *websocket.Conn
	Fm   *msg.RPCMessage
}

// SendFailMessage 消息发送失败
func (rc RespCtx) SendFailMessage(data interface{}) {
	// Notify的消息，不通知成功
	if rc.Fm.Index == 0 {
		return
	}

	rpcResp := msg.RPCResp{
		Kind:  rc.Fm.Kind + 10000,
		Index: rc.Fm.Index,
		Code:  msg.StatusCode.Fail,
		Data:  data,
	}

	err := rc.Conn.WriteMessage(msg.TypeEnum.TextMessage, rpcResp.ToBytes())
	if err != nil {
		fmt.Println(err)
	}
}

// SendSuccessfulMessage 消息发送成功
func (rc RespCtx) SendSuccessfulMessage(data interface{}) {

	// Notify的消息，不通知成功
	if rc.Fm.Index == 0 {
		return
	}

	rpcResp := msg.RPCResp{
		Kind:  rc.Fm.Kind + 10000,
		Index: rc.Fm.Index,
		Code:  msg.StatusCode.Successful,
		Data:  data,
	}

	err := rc.Conn.WriteMessage(msg.TypeEnum.TextMessage, rpcResp.ToBytes())
	if err != nil {
		fmt.Println(err)
	}
}
