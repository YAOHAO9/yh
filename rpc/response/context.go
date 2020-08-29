package response

import (
	"fmt"
	"sync"

	"github.com/YAOHAO9/yh/rpc/message"
	"github.com/gorilla/websocket"
)

var mutex sync.Mutex

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

// GetHandler 获取请求的Handler
func (respCtx RespCtx) GetHandler() string {
	return respCtx.handler
}

// SendMsg 发送消息
func (respCtx RespCtx) SendMsg(data []byte) {
	mutex.Lock()
	err := respCtx.conn.WriteMessage(message.TypeEnum.TextMessage, data)
	if err != nil {
		fmt.Println(err)
	}
	mutex.Unlock()
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

	respCtx.SendMsg(rpcResp.ToBytes())
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

	respCtx.SendMsg(rpcResp.ToBytes())
}
