package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"
	"trial/rpc/config"
	"trial/rpc/filter/handlerfilter"
	"trial/rpc/filter/rpcfilter"
	"trial/rpc/msg"
	"trial/rpc/response"

	"github.com/gorilla/websocket"
)

var requestID = 1
var maxInt64Value = 1<<63 - 1

// 唯一的RequestID
func getRequestID() int {
	requestID++
	if requestID >= maxInt64Value {
		requestID = 1
	}
	return requestID
}

var requestMap = make(map[int]func(rpcResp *msg.RPCResp))

var lock sync.Mutex

// RPCClient websocket client 连接信息
type RPCClient struct {
	Conn         *websocket.Conn
	ServerConfig *config.ServerConfig
}

// SendSysNotify 发送系统RPC通知
func (client RPCClient) SendSysNotify(session *msg.Session, message *msg.ClientMessage) {

	fm := &msg.RPCMessage{
		Kind:    msg.KindEnum.Sys,
		Handler: message.Handler,
		Data:    message.Data,
		Session: session,
	}

	respCtx := &response.RespCtx{
		Conn: client.Conn,
		Fm:   fm,
	}

	// 执行 Before filter
	if handlerfilter.Manager.Before.Exec(respCtx) {
		client.Conn.WriteMessage(msg.TypeEnum.TextMessage, fm.ToBytes())
	}
}

// SendSysRequest 发送系统请求
func (client RPCClient) SendSysRequest(session *msg.Session, message *msg.ClientMessage, cb func(rpcResp *msg.RPCResp)) {

	requestID := getRequestID()
	fm := &msg.RPCMessage{
		RequestID: requestID,
		Kind:      msg.KindEnum.Sys,
		Handler:   message.Handler,
		Data:      message.Data,
		Session:   session,
	}

	respCtx := &response.RespCtx{
		Conn: client.Conn,
		Fm:   fm,
	}

	// 执行 Before filter
	if handlerfilter.Manager.Before.Exec(respCtx) {
		lock.Lock()
		requestMap[requestID] = cb
		lock.Unlock()
		client.Conn.WriteMessage(msg.TypeEnum.TextMessage, fm.ToBytes())
	}
}

// ForwardHandlerNotify 转发Handler通知
func (client RPCClient) ForwardHandlerNotify(session *msg.Session, message *msg.ClientMessage) {

	fm := &msg.RPCMessage{
		Kind:    msg.KindEnum.Handler,
		Handler: message.Handler,
		Data:    message.Data,
		Session: session,
	}

	respCtx := &response.RespCtx{
		Conn: client.Conn,
		Fm:   fm,
	}

	// 执行 Before filter
	if handlerfilter.Manager.Before.Exec(respCtx) {
		client.Conn.WriteMessage(msg.TypeEnum.TextMessage, fm.ToBytes())
	}
}

// ForwardHandlerRequest 转发Handler请求
func (client RPCClient) ForwardHandlerRequest(session *msg.Session, message *msg.ClientMessage, cb func(rpcResp *msg.RPCResp)) {

	requestID := getRequestID()
	fm := &msg.RPCMessage{
		RequestID: requestID,
		Kind:      msg.KindEnum.Handler,
		Handler:   message.Handler,
		Data:      message.Data,
		Session:   session,
	}

	respCtx := &response.RespCtx{
		Conn: client.Conn,
		Fm:   fm,
	}

	// 执行 Before filter
	if handlerfilter.Manager.Before.Exec(respCtx) {
		lock.Lock()
		requestMap[requestID] = cb
		lock.Unlock()
		client.Conn.WriteMessage(msg.TypeEnum.TextMessage, fm.ToBytes())
	}
}

// SendRPCNotify 发送RPC通知
func (client RPCClient) SendRPCNotify(session *msg.Session, message *msg.ClientMessage) {

	fm := &msg.RPCMessage{
		Kind:    msg.KindEnum.RPC,
		Handler: message.Handler,
		Data:    message.Data,
		Session: session,
	}

	respCtx := &response.RespCtx{
		Conn: client.Conn,
		Fm:   fm,
	}

	// 执行 Before RPC filter
	if rpcfilter.Manager.Before.Exec(respCtx) {
		client.Conn.WriteMessage(msg.TypeEnum.TextMessage, fm.ToBytes())
	}
}

// SendRPCRequest 发送RPC请求
func (client RPCClient) SendRPCRequest(session *msg.Session, message *msg.ClientMessage, cb func(rpcResp *msg.RPCResp)) {

	requestID := getRequestID()
	fm := &msg.RPCMessage{
		RequestID: requestID,
		Kind:      msg.KindEnum.RPC,
		Handler:   message.Handler,
		Data:      message.Data,
		Session:   session,
	}

	respCtx := &response.RespCtx{
		Conn: client.Conn,
		Fm:   fm,
	}

	// 执行 Before RPC filter
	if rpcfilter.Manager.Before.Exec(respCtx) {
		lock.Lock()
		requestMap[requestID] = cb
		lock.Unlock()
		client.Conn.WriteMessage(msg.TypeEnum.TextMessage, fm.ToBytes())
	}
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
	fmt.Println("连接到", serverConfig.ID, "成功！！！")
	// 接收消息
	go func() {
		for {
			_, data, err := clientConn.ReadMessage()
			// 掉线检查
			if err != nil {
				clientConn.Close()
				clientConn.CloseHandler()(0, "")
				closeFunc(serverConfig.ID)
				fmt.Println("服务", serverConfig.ID, "掉线")
				break
			}
			// 解析消息
			rpcResp := &msg.RPCResp{}
			err = json.Unmarshal(data, rpcResp)
			if err != nil {
				fmt.Println("Rpc request's response body parse fail")
				continue
			}

			// 如果是request消息，则调用回调函数
			if rpcResp.RequestID != 0 {
				lock.Lock()
				requestFunc, ok := requestMap[rpcResp.RequestID]
				if ok {
					delete(requestMap, rpcResp.RequestID)
					// 执行 After RPC filter
					if rpcResp.Kind == msg.KindEnum.RPC {
						rpcfilter.Manager.After.Exec(rpcResp)
					} else if rpcResp.Kind == msg.KindEnum.Handler {
						handlerfilter.Manager.After.Exec(rpcResp)
					}
					requestFunc(rpcResp)
				}
				lock.Unlock()
				continue
			}
			fmt.Println("Notify消息:", string(data))
		}
	}()

	return &RPCClient{
		Conn:         clientConn,
		ServerConfig: serverConfig,
	}
}
