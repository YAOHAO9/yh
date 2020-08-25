package connector

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/YAOHAO9/yh/application/config"
	"github.com/YAOHAO9/yh/rpc/client"
	"github.com/YAOHAO9/yh/rpc/client/clientmanager"
	"github.com/YAOHAO9/yh/rpc/message"
	"github.com/YAOHAO9/yh/rpc/router"
	"github.com/gorilla/websocket"
)

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
		cm := &message.ClientMessage{}
		err = json.Unmarshal(data, cm)

		if err != nil {
			sendFailMessage(conn, message.KindEnum.Handler, cm.RequestID, "消息解析失败，请发送json消息")
			continue
		}

		if cm.Handler == "" {
			sendFailMessage(conn, message.KindEnum.Handler, cm.RequestID, "Hanler不能为空")
			continue
		}

		handlerInfos := strings.Split(cm.Handler, ".")
		serverKind := handlerInfos[0] // 解析出服务器类型
		cm.Handler = handlerInfos[1]  // 真正的handler

		session := &message.Session{
			UID:  uid,
			CID:  config.GetServerConfig().ID,
			Data: connInfo.data,
		}

		// 获取RPCCLint
		var rpcClient *client.RPCClient
		// 根据类型转发
		rpcClient = clientmanager.GetClientByRouter(router.Info{
			ServerKind: serverKind,
			Handler:    cm.Handler,
			Session:    *session,
		})

		if rpcClient == nil {

			tip := fmt.Sprint("找不到任何", serverKind, "服务器", ", Handler: ", cm.Handler)
			sendFailMessage(conn, message.KindEnum.Handler, cm.RequestID, tip)
			continue
		}

		if cm.RequestID == 0 {
			// 转发Notify
			rpcClient.ForwardHandlerNotify(session, cm)
		} else {
			// 转发Request
			rpcClient.ForwardHandlerRequest(session, cm, func(rpcResp *message.RPCResp) {

				clientResp := message.ClientResp{
					RequestID: rpcResp.RequestID,
					Code:      rpcResp.Code,
					Data:      rpcResp.Data,
				}

				bytes, err := json.Marshal(clientResp)
				if err != nil {
					fmt.Println("Invalid message")
				} else {
					conn.WriteMessage(message.TypeEnum.TextMessage, bytes)
				}
			})
		}
	}
}
