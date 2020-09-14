package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/YAOHAO9/yh/application/config"
	"github.com/YAOHAO9/yh/rpc/context"
	"github.com/YAOHAO9/yh/rpc/filter/handlerfilter"
	"github.com/YAOHAO9/yh/rpc/filter/rpcfilter"
	"github.com/YAOHAO9/yh/rpc/message"
	"github.com/YAOHAO9/yh/rpc/session"
	"github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

var requestID = 1
var maxInt64Value = 1<<63 - 1
var requestIDMutex sync.Mutex

// 唯一的RequestID
func genRequestID() int {

	requestIDMutex.Lock()
	defer requestIDMutex.Unlock()

	requestID++
	if requestID >= maxInt64Value {
		requestID = 1
	}
	return requestID
}

var requestMap = make(map[int]func(rpcResp *message.RPCResp))

var requestMapLock sync.RWMutex
var websocketWriteLock sync.Mutex

// RPCClient websocket client 连接信息
type RPCClient struct {
	Conn         *websocket.Conn
	ServerConfig *config.ServerConfig
}

// SendMsg 发送消息
func (client RPCClient) SendMsg(data []byte) {
	websocketWriteLock.Lock()
	client.Conn.WriteMessage(message.TypeEnum.TextMessage, data)
	websocketWriteLock.Unlock()
}

// SendRPCNotify 发送RPC通知
func (client RPCClient) SendRPCNotify(session *session.Session, rpcMsg *message.RPCMsg) {

	rpcCtx := context.GenRespCtx(client.Conn, rpcMsg)

	// 执行 Before Handler filter
	if rpcMsg.Kind == message.KindEnum.Handler {
		if !handlerfilter.Manager.Before.Exec(rpcCtx) {
			return
		}
	}

	// 执行 Before RPC filter
	if rpcMsg.Kind == message.KindEnum.RPC {
		if !rpcfilter.Manager.Before.Exec(rpcCtx) {
			return
		}
	}

	client.SendMsg(rpcMsg.ToBytes())
}

// SendRPCRequest 发送RPC请求
func (client RPCClient) SendRPCRequest(session *session.Session, rpcMsg *message.RPCMsg, cb func(rpcResp *message.RPCResp)) {

	rpcMsg.RequestID = genRequestID()

	rpcCtx := context.GenRespCtx(client.Conn, rpcMsg)

	// 执行 Before Handler filter
	if rpcMsg.Kind == message.KindEnum.Handler {
		if !handlerfilter.Manager.Before.Exec(rpcCtx) {
			return
		}
	}

	// 执行 Before RPC filter
	if rpcMsg.Kind == message.KindEnum.RPC {
		if !rpcfilter.Manager.Before.Exec(rpcCtx) {
			return
		}
	}

	requestMapLock.Lock()
	requestMap[requestID] = cb
	requestMapLock.Unlock()
	client.SendMsg(rpcMsg.ToBytes())
}

// StartClient websocket client
func StartClient(serverConfig *config.ServerConfig, zkSessionTimeout time.Duration, closeFunc func(id string)) *RPCClient {

	// Dialer
	dialer := websocket.Dialer{}
	urlString := url.URL{
		Scheme:   "ws",
		Host:     fmt.Sprint(serverConfig.Host, ":", serverConfig.Port),
		Path:     "/rpc",
		RawQuery: fmt.Sprint("token=", serverConfig.Token),
	}

	var e error
	// 当前尝试次数
	tryTimes := 0
	// 最大尝试次数
	maxTryTimes := int(50 + zkSessionTimeout/100/time.Millisecond)

	// 尝试建立连接
	var clientConn *websocket.Conn
	for tryTimes = 0; tryTimes < maxTryTimes; tryTimes++ {

		clientConn, _, e = dialer.Dial(urlString.String(), nil)
		if e == nil {
			break
		}
		// 报错则休眠100毫秒
		time.Sleep(time.Millisecond * 100)
	}

	if tryTimes >= maxTryTimes {
		// 操过最大尝试次数则报错
		panic(fmt.Sprint("Cannot create connection with ", serverConfig.ID))
	}

	// 如果超过最大尝试次数，任然有错则报错
	if e != nil {
		panic(e)
	}

	// 连接成功！！！
	logrus.Info("连接到", serverConfig.ID, "成功！！！")
	// 接收消息
	go func() {
		for {
			_, data, err := clientConn.ReadMessage()
			// 掉线检查
			if err != nil {
				clientConn.Close()
				clientConn.CloseHandler()(0, "")
				closeFunc(serverConfig.ID)
				logrus.Warn("服务", serverConfig.ID, "掉线")
				break
			}
			// 解析消息
			rpcResp := &message.RPCResp{}
			err = json.Unmarshal(data, rpcResp)
			if err != nil {
				logrus.Error("Rpc request's response body parse fail", err)
				continue
			}

			// Notify消息，不应有回调信息
			if rpcResp.RequestID == 0 {
				logrus.Error("Notify消息，不应有回调信息")
			}

			// 执行 After RPC filter
			if rpcResp.Kind == message.KindEnum.RPCResponse && !rpcfilter.Manager.After.Exec(rpcResp) {
				continue
			}

			// 执行 After Handler filter
			if rpcResp.Kind == message.KindEnum.HandlerResponse && !handlerfilter.Manager.After.Exec(rpcResp) {
				continue
			}

			// 执行回调函数
			requestMapLock.RLock()
			requestFunc, ok := requestMap[rpcResp.RequestID]
			if ok {
				delete(requestMap, rpcResp.RequestID)
				requestFunc(rpcResp)
			}
			requestMapLock.RUnlock()
		}
	}()

	return &RPCClient{
		Conn:         clientConn,
		ServerConfig: serverConfig,
	}
}
