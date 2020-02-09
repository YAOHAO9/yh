package client

import (
	"fmt"
	"net/url"
	"time"
	"trial/ws/msgtype"

	"github.com/gorilla/websocket"
)

var clientConn *websocket.Conn

// Start websocket client
func Start() {

	// Dialer
	dialer := websocket.Dialer{}
	urlString := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws", RawQuery: "id=456&token=哈哈哈"}

	// 建立连接
	conn, _, e := dialer.Dial(urlString.String(), nil)
	if e != nil {
		fmt.Println(e.Error())
		return
	}

	// 连接成功！！！
	fmt.Println("连接成功！！！")
	clientConn = conn

	// 接收消息
	go func() {
		for {
			_, data, err := clientConn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println(string(data))
		}
	}()
	Send([]byte("我来也"))
	time.Sleep(time.Hour)
}

// Send message
func Send(msg []byte) {
	clientConn.WriteMessage(msgtype.TextMessage, msg)
}
