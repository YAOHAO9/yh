package response

import (
	"fmt"
	"trial/rpc/msg"
	"trial/rpc/msgtype"

	"github.com/gorilla/websocket"
)

// SendFailMessage 消息发送失败
func SendFailMessage(respConn *websocket.Conn, Kind int, index int, data interface{}) {
	fmt.Println(data)
	// Notify的消息，不通知成功
	if index == 0 {
		return
	}

	response := msg.RPCResp{
		Kind:  Kind,
		Index: index,
		Code:  msg.StatusCode().Fail,
		Data:  data,
	}

	err := respConn.WriteMessage(msgtype.TextMessage, response.ToBytes())
	if err != nil {
		fmt.Println(err)
	}
}

// SendSuccessfulMessage 消息发送成功
func SendSuccessfulMessage(respConn *websocket.Conn, Kind int, index int, data interface{}) {

	// Notify的消息，不通知成功
	if index == 0 {
		return
	}

	response := msg.RPCResp{
		Kind:  Kind,
		Index: index,
		Code:  msg.StatusCode().Successful,
		Data:  data,
	}

	err := respConn.WriteMessage(msgtype.TextMessage, response.ToBytes())
	if err != nil {
		fmt.Println(err)
	}
}
