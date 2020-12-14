package context

import (
	"fmt"
	"sync"

	"github.com/YAOHAO9/pine/rpc/handler/routercompress"
	"github.com/YAOHAO9/pine/rpc/message"
	"github.com/YAOHAO9/pine/rpc/session"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

var sendMsgMutex sync.Mutex
var requestIDRwMutex sync.RWMutex

// RPCCtx response context
type RPCCtx struct {
	conn      *websocket.Conn
	requestID int32
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
		handlerName := routercompress.GetHandlerByCode(bytes[0])
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

// GetRequestID 获取请求的GetRequestID
func (rpcCtx *RPCCtx) GetRequestID() int32 {

	requestIDRwMutex.RLock()
	defer requestIDRwMutex.RUnlock()

	return rpcCtx.requestID
}

// SendMsg 消息发送失败
func (rpcCtx *RPCCtx) SendMsg(data []byte) {

	requestIDRwMutex.Lock()
	defer requestIDRwMutex.Unlock()

	// Notify的消息，不通知成功
	if rpcCtx.requestID == 0 {
		if data == nil {
			return
		}
		logrus.Warn(fmt.Sprintf("NotifyHandler(%s)不需要回复消息", rpcCtx.handler))
		return
	}
	// 重复回复
	if rpcCtx.requestID == -1 {
		logrus.Warn(fmt.Sprintf("Handler(%s)请勿重复回复消息", rpcCtx.handler))
		return
	}
	requestID := rpcCtx.requestID
	// 标记为已回复消息
	rpcCtx.requestID = -1
	// response
	rpcResp := &message.PineMessage{
		Route:     rpcCtx.handler,
		RequestID: &requestID,
		Data:      data,
	}

	bytes, err := proto.Marshal(rpcResp)
	if err != nil {
		logrus.Error(err)
		return
	}

	sendMsgMutex.Lock()
	defer sendMsgMutex.Unlock()
	err = rpcCtx.conn.WriteMessage(message.TypeEnum.BinaryMessage, bytes)
	if err != nil {
		logrus.Error(err)
	}

}

// ToString 格式化消息
func (rpcCtx RPCCtx) ToString() string {
	return fmt.Sprintf("RPC RequestID: %d, Data: %+v", rpcCtx.requestID, rpcCtx.RawData)
}
