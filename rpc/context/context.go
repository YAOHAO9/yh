package context

import (
	"fmt"
	"sync"

	"github.com/YAOHAO9/yh/rpc/message"
	"github.com/YAOHAO9/yh/rpc/session"
	"github.com/gorilla/websocket"
)

var mutex sync.Mutex

// RPCCtx response context
type RPCCtx struct {
	conn      *websocket.Conn
	kind      int
	requestID int
	handler   string
	Data      interface{}
	Session   *session.Session
}

// GenRespCtx 创建一个response上下文
func GenRespCtx(conn *websocket.Conn, rpcMsg *message.RPCMessage) *RPCCtx {
	return &RPCCtx{
		conn:      conn,
		kind:      rpcMsg.Kind,
		requestID: rpcMsg.RequestID,
		handler:   rpcMsg.Handler,
		Data:      rpcMsg.Data,
		Session:   rpcMsg.Session,
	}
}

// GetHandler 获取请求的Handler
func (rpcCtx RPCCtx) GetHandler() string {
	return rpcCtx.handler
}

// SendMsg 发送消息
func (rpcCtx RPCCtx) SendMsg(data []byte) {
	mutex.Lock()
	err := rpcCtx.conn.WriteMessage(message.TypeEnum.TextMessage, data)
	if err != nil {
		fmt.Println(err)
	}
	mutex.Unlock()
}

// SendFailMessage 消息发送失败
func (rpcCtx RPCCtx) SendFailMessage(data interface{}) {
	// Notify的消息，不通知成功
	if rpcCtx.requestID == 0 {
		return
	}

	rpcResp := message.RPCResp{
		Kind:      rpcCtx.kind + 10000,
		RequestID: rpcCtx.requestID,
		Code:      message.StatusCode.Fail,
		Data:      data,
	}

	rpcCtx.SendMsg(rpcResp.ToBytes())
}

// SendSuccessfulMessage 消息发送成功
func (rpcCtx RPCCtx) SendSuccessfulMessage(data interface{}) {

	// Notify的消息，不通知成功
	if rpcCtx.requestID == 0 {
		return
	}

	rpcResp := message.RPCResp{
		Kind:      rpcCtx.kind + 10000,
		RequestID: rpcCtx.requestID,
		Code:      message.StatusCode.Successful,
		Data:      data,
	}

	rpcCtx.SendMsg(rpcResp.ToBytes())
}
