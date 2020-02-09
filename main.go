package main

import (
	"fmt"
	"net/http"
	"trial/util"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1
	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2
	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8
	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9
	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
)

// ConnInfo 用户连接信息
type ConnInfo struct {
	id   int
	conn *websocket.Conn
	data chan interface{}
}

var connMap = make(map[string]*ConnInfo)

func main() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {

		// 建立连接
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("连接失败", err.Error())
			return
		}

		// 断开连接自动清除连接信息
		id := r.URL.Query().Get("id")
		conn.SetCloseHandler(func(code int, text string) error {
			delete(connMap, id)
			return nil
		})

		// 用户认证
		token := r.URL.Query().Get("token")
		fmt.Println("Id: ", id, " Token: ", token)
		if id == "" || token == "" {
			fmt.Println("用户校验失败!!!")
			err := conn.WriteMessage(TextMessage, []byte("认证失败"))
			if err != nil {
				fmt.Println("发送认证失败消息失败: ", err.Error())
			}
			conn.Close()
			return
		}

		if oldConnInfo, ok := connMap[id]; ok {
			oldConnInfo.conn.Close()
			oldConnInfo.conn.CloseHandler()
		}

		connInfo := &ConnInfo{id: 1, conn: conn, data: make(chan interface{})}
		connMap[id] = connInfo

		// 开始接收消息
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				data, ok := util.GetValue(err, "Code")
				if ok {
					if data == 1001 {
						fmt.Println("Websocket连接断开")
						break
					}
				}

				fmt.Println("Error: ", err.Error())
				break
			}
			// connInfo.data <- data
			fmt.Println(string(data))
			conn.WriteMessage(TextMessage, []byte("哈哈哈"))
		}
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("你好世界！！！"))
	})
	port := "8080"
	fmt.Println("Server started http://localhost:" + port)
	err := http.ListenAndServe(":"+port, nil)
	fmt.Println(err.Error())
}
