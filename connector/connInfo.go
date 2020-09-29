package connector

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/YAOHAO9/yh/application/config"
	"github.com/YAOHAO9/yh/connector/filter"
	"github.com/YAOHAO9/yh/connector/msg"
	"github.com/YAOHAO9/yh/rpc"
	"github.com/YAOHAO9/yh/rpc/client/clientmanager"
	"github.com/YAOHAO9/yh/rpc/context"
	"github.com/YAOHAO9/yh/rpc/message"
	"github.com/YAOHAO9/yh/rpc/session"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// HandlerPrefix  Handler 前缀
var HandlerPrefix = "__"

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
	clientMsgResp := msg.ClientResp{
		RequestID: requestID,
		Code:      code,
		Data:      data,
	}
	mutex.Lock()
	err := connInfo.conn.WriteMessage(message.TypeEnum.TextMessage, clientMsgResp.ToBytes())
	mutex.Unlock()
	if err != nil {
		logrus.Error(err)
	}
}

// 主动推送消息
func (connInfo ConnInfo) notify(route string, data interface{}) {

	notify := msg.ClientNotify{
		Route: route,
		Data:  data,
	}

	mutex.Lock()
	err := connInfo.conn.WriteMessage(message.TypeEnum.TextMessage, notify.ToBytes())
	mutex.Unlock()
	if err != nil {
		logrus.Error(err)
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
		clientMessage := &msg.ClientReq{}
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

		rpcMsg := &message.RPCMsg{
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

		rpcCtx := context.GenRespCtx(conn, rpcMsg)

		if !filter.Before.Exec(rpcCtx) {
			return
		}

		if clientMessage.RequestID == 0 {
			rpc.Notify.ToServer(rpcClient.ServerConfig.ID, session, HandlerPrefix+handler, clientMessage.Data)
		} else {
			// 转发Request
			rpc.Request.ToServer(rpcClient.ServerConfig.ID, session, HandlerPrefix+handler, clientMessage.Data, func(rpcResp *message.RPCResp) {
				filter.After.Exec(rpcResp)
				connInfo.response(clientMessage.RequestID, rpcResp.Code, rpcResp.Data)
			})
		}
	}
}
