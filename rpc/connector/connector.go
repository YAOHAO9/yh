package connector

import (
	"encoding/json"
	"fmt"
	"net/http"
	"trial/rpc/client"
	"trial/rpc/client/clientmanager"
	"trial/rpc/config"
	"trial/rpc/msg"
	"trial/rpc/msg/msgkind"
	"trial/rpc/msg/msgtype"

	"github.com/gorilla/websocket"
)

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
	connInfo := &ConnInfo{id: 1, conn: conn}
	ConnMap[id] = connInfo

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
			sendFailMessage(conn, msgkind.Handler, cm.Index, "消息解析失败，请发送json消息")
			continue
		}
		fmt.Println(cm.Data)
		if cm.Handler == "" {
			sendFailMessage(conn, msgkind.Handler, cm.Index, "Hanler不能为空")
			continue
		}

		// 获取RPCCLint
		var rpcClient *client.RPCClient
		if cm.ServerID != "" {
			rpcClient = clientmanager.GetClientByID(cm.ServerID)
		} else {
			rpcClient = clientmanager.GetRandClientByKind(cm.Kind)
		}

		if rpcClient == nil {
			sendFailMessage(conn, msgkind.Handler, cm.Index, fmt.Sprint("服务器不存在, Kind: ", cm.Kind, ", ServerID: ", cm.ServerID))
			continue
		}

		session := &msg.Session{
			UID:  id,
			CID:  config.GetServerConfig().ID,
			Data: connInfo.data,
		}

		if cm.Index == 0 {
			// 转发Notify
			rpcClient.SendHandlerNotify(session, cm)
		} else {
			// 转发Request
			rpcClient.SendHandlerRequest(session, cm, func(data interface{}) {
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

// sendFailMessage 消息发送失败
func sendFailMessage(respConn *websocket.Conn, Kind int, index int, data interface{}) {
	fmt.Println(data)
	// Notify的消息，不通知成功
	if index == 0 {
		return
	}

	rpcResp := msg.ClientResp{
		Index: index,
		Code:  msg.StatusCode().Fail,
		Data:  data,
	}

	err := respConn.WriteMessage(msgtype.TextMessage, rpcResp.ToBytes())
	if err != nil {
		fmt.Println(err)
	}
}
