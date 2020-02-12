package client

import (
	"fmt"
	"net/url"
	"trial/rpc/msgtype"

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
				fmt.Println("ReadMessage error: ", err)
				clientConn.Close()
				clientConn.CloseHandler()(0, "")
				break
			}
			fmt.Println(string(data))
		}
	}()
	for i := 0; i < 100; i++ {
		Send([]byte(fmt.Sprint(i)))
	}

}

// Send message
func Send(msg []byte) {
	clientConn.WriteMessage(msgtype.TextMessage, msg)
}

// OnClose set listener for close
func OnClose(h func(code int, text string) error) {
	clientConn.SetCloseHandler(h)
}
