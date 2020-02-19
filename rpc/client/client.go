package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"
	"trial/rpc/config"
	"trial/rpc/filter"
	"trial/rpc/filter/rpcfilter"
	"trial/rpc/msg"
	"trial/rpc/msgtype"

	"github.com/gorilla/websocket"
)

var requestMap = make(map[int]func(rm *msg.ResponseMessage))
var lock sync.Mutex

// RPCClient websocket client 连接信息
type RPCClient struct {
	Conn         *websocket.Conn
	ServerConfig *config.ServerConfig
}

// SendHandlerNotify send handler message
func (client RPCClient) SendHandlerNotify(session *msg.Session, message *msg.Message) {
	fm := &msg.ForwardMessage{IsRPC: false, Msg: message, Session: session}
	// 执行 Before filter
	if filter.BeforeFilterManager().Exec(client.Conn, fm) {
		client.Conn.WriteMessage(msgtype.TextMessage, fm.ToBytes())
	}
}

// SendRPCNotify send rpc message
func (client RPCClient) SendRPCNotify(session *msg.Session, message *msg.Message) {
	fm := &msg.ForwardMessage{IsRPC: true, Msg: message, Session: session}
	// 执行 Before RPC filter
	if rpcfilter.BeforeFilterManager().Exec(client.Conn, fm) {
		client.Conn.WriteMessage(msgtype.TextMessage, fm.ToBytes())
	}

}

// SendHandlerRequest send handler message
func (client RPCClient) SendHandlerRequest(session *msg.Session, msgIndex int, message *msg.Message, cb func(data *msg.ResponseMessage)) {
	fm := &msg.ForwardMessage{IsRPC: false, Msg: message, Session: session}
	// 执行 Before filter
	if filter.BeforeFilterManager().Exec(client.Conn, fm) {
		client.Conn.WriteMessage(msgtype.TextMessage, fm.ToBytes())
		lock.Lock()
		requestMap[msgIndex] = cb
		lock.Unlock()
	}
}

// SendRPCRequest send rpc message
func (client RPCClient) SendRPCRequest(session *msg.Session, msgIndex int, message *msg.Message, cb func(data *msg.ResponseMessage)) {
	fm := &msg.ForwardMessage{IsRPC: true, Msg: message, Session: session}
	// 执行 Before RPC filter
	if rpcfilter.BeforeFilterManager().Exec(client.Conn, fm) {
		client.Conn.WriteMessage(msgtype.TextMessage, fm.ToBytes())
		lock.Lock()
		requestMap[msgIndex] = cb
		lock.Unlock()
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
			responseMessage := &msg.ResponseMessage{}
			err = json.Unmarshal(data, responseMessage)
			if err != nil {
				fmt.Println("Rpc request's response body parse fail")
				continue
			}

			// 如果是request消息，则调用回调函数
			if responseMessage.Index != 0 {
				lock.Lock()
				requestFunc, ok := requestMap[responseMessage.Index]
				if ok {
					delete(requestMap, responseMessage.Index)
					if responseMessage.IsRPC {
						rpcfilter.AfterFilterManager().Exec(responseMessage)
					} else {
						filter.AfterFilterManager().Exec(responseMessage)
					}
					requestFunc(responseMessage)
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
