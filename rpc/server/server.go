package server

import (
	"fmt"
	"net/http"
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

	if oldConnInfo, ok := ConnMap[id]; ok {
		oldConnInfo.conn.CloseHandler()(0, "关闭重复连接")
	}

	connInfo := &ConnInfo{id: 1, conn: conn, data: make(chan interface{})}
	ConnMap[id] = connInfo

	// 开始接收消息
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			conn.CloseHandler()(0, err.Error())
			break
		}
		// connInfo.data <- data
		fmt.Println(string(data))
		conn.WriteMessage(msgtype.TextMessage, data)
	}
}
