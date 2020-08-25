package response

import (
	"fmt"

	"github.com/YAOHAO9/yh/rpc/message"
	"github.com/gorilla/websocket"
)

// RespCtx response context
type RespCtx struct {
	conn      *websocket.Conn
	kind      int
	requestID int
	handler   string
	Data      interface{}
	Session   *message.Session
}

// GenRespCtx 创建一个response上下文
func GenRespCtx(conn *websocket.Conn, rpcMsg *message.RPCMessage) *RespCtx {
	return &RespCtx{
		conn:      conn,
		kind:      rpcMsg.Kind,
		requestID: rpcMsg.RequestID,
		handler:   rpcMsg.Handler,
		Data:      rpcMsg.Data,
		Session:   rpcMsg.Session,
	}
}

// GetHandler 消息发送失败
func (respCtx RespCtx) GetHandler() string {
	return respCtx.handler
}

// SendFailMessage 消息发送失败
func (respCtx RespCtx) SendFailMessage(data interface{}) {
	// Notify的消息，不通知成功
	if respCtx.requestID == 0 {
		return
	}

	rpcResp := message.RPCResp{
		Kind:      respCtx.kind + 10000,
		RequestID: respCtx.requestID,
		Code:      message.StatusCode.Fail,
		Data:      data,
	}

	err := respCtx.conn.WriteMessage(message.TypeEnum.TextMessage, rpcResp.ToBytes())
	if err != nil {
		fmt.Println(err)
	}
}

// SendSuccessfulMessage 消息发送成功
func (respCtx RespCtx) SendSuccessfulMessage(data interface{}) {

	// Notify的消息，不通知成功
	if respCtx.requestID == 0 {
		return
	}

	rpcResp := message.RPCResp{
		Kind:      respCtx.kind + 10000,
		RequestID: respCtx.requestID,
		Code:      message.StatusCode.Successful,
		Data:      data,
	}

	err := respCtx.conn.WriteMessage(message.TypeEnum.TextMessage, rpcResp.ToBytes())
	if err != nil {
		fmt.Println(err)
	}
}
