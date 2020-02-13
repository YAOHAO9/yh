package client

import (
	"fmt"
	"net/url"
	"trial/config"

	"github.com/gorilla/websocket"
)

// Start websocket client
func Start(serverConfig *config.ServerConfig) (clientConn *websocket.Conn) {
	// Dialer
	dialer := websocket.Dialer{}
	urlString := url.URL{
		Scheme:   "ws",
		Host:     fmt.Sprint(serverConfig.Host, ":", serverConfig.Port),
		Path:     "/rpc",
		RawQuery: fmt.Sprint("token=", serverConfig.Token),
	}
	fmt.Println(urlString.String())
	// 建立连接
	conn, _, e := dialer.Dial(urlString.String(), nil)
	if e != nil {
		panic(e)
	}

	// 连接成功！！！
	fmt.Println("连接成功！！！")
	clientConn = conn

	// 接收消息
	go func() {
		for {
			_, data, err := clientConn.ReadMessage()
			if err != nil {
				clientConn.Close()
				clientConn.CloseHandler()(0, "")
				panic(err)
			}
			fmt.Println(string(data))
		}
	}()

	return clientConn
}

// Send message
// func Send(msg []byte) {
// 	clientConn.WriteMessage(msgtype.TextMessage, msg)
// }

// // OnClose set listener for close
// func OnClose(h func(code int, text string) error) {
// 	clientConn.SetCloseHandler(h)
// }
