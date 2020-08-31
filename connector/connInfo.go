package connector

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/YAOHAO9/yh/application/config"
	"github.com/YAOHAO9/yh/rpc/client/clientmanager"
	"github.com/YAOHAO9/yh/rpc/message"
	"github.com/YAOHAO9/yh/rpc/session"
	"github.com/gorilla/websocket"
)

var mutex sync.Mutex

// ConnInfo 用户连接信息
type ConnInfo struct {
	uid  string
	conn *websocket.Conn
	data map[string]interface{}
}

// Get 从session中查找一个值
func (connInfo ConnInfo) Get(key string) interface{} {
	return connInfo.data[key]
}

// Set 往session中设置一个键值对
func (connInfo ConnInfo) Set(key string, v interface{}) {
	connInfo.data[key] = v
}

// 回复request
func (connInfo ConnInfo) response(requestID int, code int, data interface{}) {
	clientMsgResp := ClientMsgResp{
		RequestID: requestID,
		Code:      code,
		Data:      data,
	}
	mutex.Lock()
	err := connInfo.conn.WriteMessage(message.TypeEnum.TextMessage, clientMsgResp.ToBytes())
	mutex.Unlock()
	if err != nil {
		fmt.Println(err)
	}
}

// 主动推送消息
func (connInfo ConnInfo) notify(route string, data interface{}) {

	notify := ClientNotify{
		Route: route,
		Data:  data,
	}

	mutex.Lock()
	err := connInfo.conn.WriteMessage(message.TypeEnum.TextMessage, notify.ToBytes())
	mutex.Unlock()
	if err != nil {
		fmt.Println(err)
	}
}

// StartReceiveMsg 开始接收消息
func (connInfo ConnInfo) StartReceiveMsg() {
	uid := connInfo.uid
	conn := connInfo.conn
	// 开始接收消息
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			conn.CloseHandler()(0, err.Error())
			break
		}
		// 解析消息
		clientMessage := &ClientMsg{}
		err = json.Unmarshal(data, clientMessage)

		if err != nil {
			connInfo.response(clientMessage.RequestID, message.StatusCode.Fail, "消息解析失败，请发送json消息")
			continue
		}

		if clientMessage.Route == "" {
			connInfo.response(clientMessage.RequestID, message.StatusCode.Fail, "Route不能为空")
			continue
		}

		handlerInfos := strings.Split(clientMessage.Route, ".")

		serverKind := handlerInfos[0] // 解析出服务器类型
		handler := handlerInfos[1]    // 真正的handler

		session := &session.Session{
			UID:  uid,
			CID:  config.GetServerConfig().ID,
			Data: connInfo.data,
		}

		fmt.Println(connInfo.data)

		rpcMsg := &message.RPCMsg{
			Kind:    message.KindEnum.Handler,
			Handler: handler,
			Data:    clientMessage.Data,
			Session: session,
		}
		// 获取RPCCLint
		rpcClient := clientmanager.GetClientByRouter(serverKind, rpcMsg)

		if rpcClient == nil {
			tip := fmt.Sprint("找不到任何", serverKind, "服务器", ", Route: ", clientMessage.Route)
			connInfo.response(clientMessage.RequestID, message.StatusCode.Fail, tip)
			continue
		}

		if clientMessage.RequestID == 0 {
			// 转发Notify
			rpcClient.SendRPCNotify(session, rpcMsg)
		} else {
			// 转发Request
			rpcClient.SendRPCRequest(session, rpcMsg, func(rpcResp *message.RPCResp) {
				connInfo.response(clientMessage.RequestID, rpcResp.Code, rpcResp.Data)
			})
		}
	}
}
