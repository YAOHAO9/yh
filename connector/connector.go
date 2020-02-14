package connector

import (
	"encoding/json"
	"fmt"
	"net/http"
	"trial/config"
	"trial/rpc/client"
	"trial/rpc/msg"
	"trial/rpc/msgtype"

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

	// 开始接收消息
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			conn.CloseHandler()(0, err.Error())
			break
		}
		// 解析消息
		connectorMessage := &msg.Message{}
		err = json.Unmarshal(data, connectorMessage)

		if err != nil {
			SendFailMessage(conn, connectorMessage.Index, "无效的消息类型")
			continue
		}

		// 获取RPCCLint
		var connInfo *client.RPCClient
		if connectorMessage.ServerID != "" {
			connInfo = client.GetClientByID(connectorMessage.ServerID)
		} else if connectorMessage.ServerID == config.GetServerConfig().ID || connectorMessage.Kind == "connector" {
			// Connector 消息
			dealMessage(conn, connectorMessage)
		} else {
			connInfo = client.GetRandClientByKind(connectorMessage.Kind)
		}

		if connInfo == nil {
			SendFailMessage(conn, connectorMessage.Index, fmt.Sprint("服务器不存在,Kind:", connectorMessage.Kind, "ServerID:", connectorMessage.ServerID))
			continue
		}

		if connectorMessage.Index == 0 {
			// 转发Notify
			connInfo.SendRPCNotify(data)
		} else {
			// 转发Request
			connInfo.SendRPCRequest(connectorMessage.Index, data, func(data []byte) {
				conn.WriteMessage(msgtype.TextMessage, data)
			})
		}
	}
}

// SendFailMessage 消息发送失败
func SendFailMessage(conn *websocket.Conn, index int, data interface{}) {
	fmt.Println(data)
	// Notify的消息，不通知成功
	if index == 0 {
		return
	}

	response := msg.ResponseMessage{
		Index: index,
		Code:  msg.StatusCode().Fail,
		Data:  data,
	}

	err := conn.WriteMessage(msgtype.TextMessage, response.ToBytes())
	if err != nil {
		fmt.Println(err)
	}
}

// SendSuccessfulMessage 消息发送成功
func SendSuccessfulMessage(conn *websocket.Conn, index int, data interface{}) {

	// Notify的消息，不通知成功
	if index == 0 {
		return
	}

	response := msg.ResponseMessage{
		Index: index,
		Code:  msg.StatusCode().Successful,
		Data:  data,
	}

	err := conn.WriteMessage(msgtype.TextMessage, response.ToBytes())
	if err != nil {
		fmt.Println(err)
	}
}
func dealMessage(conn *websocket.Conn, connectorMessage *msg.Message) {
	data, _ := json.Marshal(connectorMessage)
	err := conn.WriteMessage(msgtype.TextMessage, data)
	if err != nil {
		fmt.Println(err)
	}
}
