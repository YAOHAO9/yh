package connector

import (
	"encoding/json"
	"fmt"
	"net/http"
	"trial/rpc/client"
	"trial/rpc/client/clientmanager"
	"trial/rpc/msg"
	"trial/rpc/msgtype"
	"trial/rpc/response"

	"github.com/gorilla/websocket"
)

// ConnInfo 用户连接信息
type ConnInfo struct {
	id   int
	conn *websocket.Conn
	data chan interface{}
}

// ConnMap socket connection map
var ConnMap = make(map[string]*ConnInfo)

var upgrader = websocket.Upgrader{}

// WebSocketHandler deal with ws request
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {

	// 建立连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("连接失败", err.Error())
		return
	}

	// 断开连接自动清除连接信息
	id := r.URL.Query().Get("id")
	conn.SetCloseHandler(func(code int, text string) error {
		delete(ConnMap, id)
		conn.Close()
		fmt.Println("CloseHandler: ", text)
		return nil
	})

	// 用户认证
	token := r.URL.Query().Get("token")
	fmt.Println("Id: ", id, " Token: ", token)

	if id == "" || token == "" {
		fmt.Println("用户校验失败!!!")
		err := conn.WriteMessage(msgtype.TextMessage, []byte("认证失败"))
		if err != nil {
			fmt.Println("发送认证失败消息失败: ", err.Error())
		}
		conn.CloseHandler()(0, "认证失败")
		return
	}

	// 防止重复连接
	if oldConnInfo, ok := ConnMap[id]; ok {
		oldConnInfo.conn.CloseHandler()(0, "关闭重复连接")
	}

	// 保存连接信息
	connInfo := &ConnInfo{id: 1, conn: conn, data: make(chan interface{})}
	ConnMap[id] = connInfo

	session := &msg.Session{UID: id}
	// 开始接收消息
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			conn.CloseHandler()(0, err.Error())
			break
		}
		// 解析消息
		message := &msg.Message{}
		err = json.Unmarshal(data, message)

		if err != nil {
			response.SendFailMessage(conn, false, message.Index, "消息解析失败，请发送json消息")
			continue
		}

		if message.Handler == "" {
			response.SendFailMessage(conn, false, message.Index, "Hanler不能为空")
			continue
		}

		// 获取RPCCLint
		var connInfo *client.RPCClient
		if message.ServerID != "" {
			connInfo = clientmanager.GetClientByID(message.ServerID)
		} else {
			connInfo = clientmanager.GetRandClientByKind(message.Kind)
		}

		if connInfo == nil {
			response.SendFailMessage(conn, false, message.Index, fmt.Sprint("服务器不存在, Kind: ", message.Kind, ", ServerID: ", message.ServerID))
			continue
		}

		if message.Index == 0 {
			// 转发Notify
			connInfo.SendHandlerNotify(session, message)
		} else {
			// 转发Request
			connInfo.SendHandlerRequest(session, message.Index, message, func(rm *msg.ResponseMessage) {
				data, err := json.Marshal(rm)
				if err != nil {
					fmt.Println("Invalid message")
				} else {
					conn.WriteMessage(msgtype.TextMessage, data)
				}
			})
		}
	}
}
