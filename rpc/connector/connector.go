package connector

import (
	"encoding/json"
	"fmt"
	"net/http"
	"trial/rpc/client"
	"trial/rpc/client/clientmanager"
	"trial/rpc/msg"
	"trial/rpc/msg/msgkind"
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
		cm := &msg.ClientMessage{}
		err = json.Unmarshal(data, cm)

		if err != nil {
			response.SendFailMessage(conn, msgkind.Handler, cm.Index, "消息解析失败，请发送json消息")
			continue
		}

		if cm.Handler == "" {
			response.SendFailMessage(conn, msgkind.Handler, cm.Index, "Hanler不能为空")
			continue
		}

		// 获取RPCCLint
		var connInfo *client.RPCClient
		if cm.ServerID != "" {
			connInfo = clientmanager.GetClientByID(cm.ServerID)
		} else {
			connInfo = clientmanager.GetRandClientByKind(cm.Kind)
		}

		if connInfo == nil {
			response.SendFailMessage(conn, msgkind.Handler, cm.Index, fmt.Sprint("服务器不存在, Kind: ", cm.Kind, ", ServerID: ", cm.ServerID))
			continue
		}

		if cm.Index == 0 {
			// 转发Notify
			connInfo.SendHandlerNotify(session, cm)
		} else {
			// 转发Request
			connInfo.SendHandlerRequest(session, cm, func(data interface{}) {

				clientResp := msg.ClientResp{
					Index: cm.Index,
					Data:  data,
				}

				bytes, err := json.Marshal(clientResp)
				if err != nil {
					fmt.Println("Invalid message")
				} else {
					conn.WriteMessage(msgtype.TextMessage, bytes)
				}
			})
		}
	}
}
