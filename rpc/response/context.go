package response

import (
	"fmt"
	"trial/rpc/msg"
	"trial/rpc/msgtype"

	"github.com/gorilla/websocket"
)

// RespCtx response context
type RespCtx struct {
	conn *websocket.Conn
	fm   *msg.ForwardMessage
}

// SendFailMessage 消息发送失败
func (rc RespCtx) SendFailMessage(data interface{}) {
	// Notify的消息，不通知成功
	if rc.fm.Index == 0 {
		return
	}

	response := msg.ResponseMessage{
		Kind:  rc.fm.Kind + 10000,
		Index: rc.fm.Index,
		Code:  msg.StatusCode().Fail,
		Data:  data,
	}

	err := rc.conn.WriteMessage(msgtype.TextMessage, response.ToBytes())
	if err != nil {
		fmt.Println(err)
	}
}

// SendSuccessfulMessage 消息发送成功
func (rc RespCtx) SendSuccessfulMessage(data interface{}) {

	// Notify的消息，不通知成功
	if rc.fm.Index == 0 {
		return
	}

	response := msg.ResponseMessage{
		Kind:  rc.fm.Kind + 10000,
		Index: rc.fm.Index,
		Code:  msg.StatusCode().Successful,
		Data:  data,
	}

	err := rc.conn.WriteMessage(msgtype.TextMessage, response.ToBytes())
	if err != nil {
		fmt.Println(err)
	}
}
