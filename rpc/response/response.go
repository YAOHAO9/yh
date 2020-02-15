package response

import (
	"fmt"
	"trial/rpc/msg"
	"trial/rpc/msgtype"

	"github.com/gorilla/websocket"
)

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
