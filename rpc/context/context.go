package context

import (
	"fmt"
	"sync"

	"github.com/YAOHAO9/pine/rpc/message"
	"github.com/YAOHAO9/pine/rpc/session"
	"github.com/YAOHAO9/pine/service/compressservice"
	"github.com/YAOHAO9/pine/util"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var sendMsgMutex sync.Mutex

// RPCCtx response context
type RPCCtx struct {
	conn      *websocket.Conn
	requestID *int32
	handler   string
	From      string
	RawData   []byte `json:",omitempty"`
	Session   *session.Session
}

// GenRespCtx 创建一个response上下文
func GenRespCtx(conn *websocket.Conn, rpcMsg *message.RPCMsg) *RPCCtx {
	return &RPCCtx{
		conn:      conn,
		requestID: rpcMsg.RequestID,
		handler:   rpcMsg.Handler,
		From:      rpcMsg.From,
		RawData:   rpcMsg.RawData,
		Session:   rpcMsg.Session,
	}
}

// GetHandler 获取请求的Handler
func (rpcCtx *RPCCtx) GetHandler() string {
	bytes := []byte(rpcCtx.handler)
	if len(bytes) == 1 {
		handlerName := compressservice.Handler.GetHandlerByCode(bytes[0])
		if handlerName != "" {
			return handlerName
		}
	}
	return rpcCtx.handler
}

// SetHandler 设置Handler
func (rpcCtx *RPCCtx) SetHandler(handler string) {
	rpcCtx.handler = handler
}

// GetRequestID 获取请求的RequestID
func (rpcCtx *RPCCtx) GetRequestID() int32 {

	if rpcCtx.requestID == nil {
		return 0
	}

	return *rpcCtx.requestID
}

// SetRequestID 设置请求的RequestID
func (rpcCtx *RPCCtx) SetRequestID(id int32) {
	rpcCtx.requestID = &id
}

// SendMsg 消息发送失败
func (rpcCtx *RPCCtx) SendMsg(data []byte) {

	requestID := rpcCtx.GetRequestID()

	// Notify的消息，不通知成功
	if requestID == 0 {
		if data == nil {
			return
		}
		logrus.Error(fmt.Sprintf("NotifyHandler(%s)不需要回复消息", rpcCtx.handler))
		return
	}
	// 重复回复
	if requestID == -1 {
		logrus.Warn(fmt.Sprintf("Handler(%s)请勿重复回复消息", rpcCtx.handler))
		return
	}
	// 标记为已回复消息
	*rpcCtx.requestID = -1
	// response
	rpcResp := &message.PineMsg{
		Route:     rpcCtx.handler,
		RequestID: &requestID,
		Data:      data,
	}

	sendMsgMutex.Lock()
	defer sendMsgMutex.Unlock()
	err := rpcCtx.conn.WriteMessage(message.TypeEnum.BinaryMessage, util.ToBytes(rpcResp))
	if err != nil {
		logrus.Error(err)
	}

}

// ToString 格式化消息
func (rpcCtx *RPCCtx) ToString() string {
	return fmt.Sprintf("Handler %s, RequestID: %d, Data: %+v", rpcCtx.GetHandler(), rpcCtx.GetRequestID(), rpcCtx.RawData)
}
