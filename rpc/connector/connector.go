package connector

import (
	"encoding/json"
	"fmt"
	"github.com/YAOHAO9/yh/application/config"
	"github.com/YAOHAO9/yh/rpc/client"
	"github.com/YAOHAO9/yh/rpc/client/clientmanager"
	"github.com/YAOHAO9/yh/rpc/msg"
	"github.com/YAOHAO9/yh/rpc/router"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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
		err := conn.WriteMessage(msg.TypeEnum.TextMessage, []byte("认证失败"))
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
			sendFailMessage(conn, msg.KindEnum.Handler, cm.RequestID, "消息解析失败，请发送json消息")
			continue
		}

		fmt.Println(cm.Data)

		if cm.Handler == "" {
			sendFailMessage(conn, msg.KindEnum.Handler, cm.RequestID, "Hanler不能为空")
			continue
		}

		handlerInfos := strings.Split(cm.Handler, ".")
		serverKind := handlerInfos[0] // 解析出服务器类型
		cm.Handler = handlerInfos[1]  // 真正的handler

		session := &msg.Session{
			UID:  id,
			CID:  config.GetServerConfig().ID,
			Data: connInfo.data,
		}

		// 获取RPCCLint
		var rpcClient *client.RPCClient
		if cm.ServerID != "" {
			// 转发到指定的后端服务器
			rpcClient = clientmanager.GetClientByID(cm.ServerID)
		} else {
			// 根据类型转发
			rpcClient = clientmanager.GetClientByRouter(router.Info{
				ServerKind: serverKind,
				Handler:    cm.Handler,
				Session:    *session,
			})
		}

		if rpcClient == nil {

			tip := ""
			if cm.ServerID == "" {
				tip = fmt.Sprint("找不到任何", serverKind, "服务器", ", Handler: ", cm.Handler)
			} else {
				tip = fmt.Sprint("服务器: ", cm.ServerID, "不存在")
			}

			sendFailMessage(conn, msg.KindEnum.Handler, cm.RequestID, tip)
			continue
		}

		if cm.RequestID == 0 {
			// 转发Notify
			rpcClient.ForwardHandlerNotify(session, cm)
		} else {
			// 转发Request
			rpcClient.ForwardHandlerRequest(session, cm, func(rpcResp *msg.RPCResp) {

				clientResp := msg.ClientResp{
					RequestID: rpcResp.RequestID,
					Code:      rpcResp.Code,
					Data:      rpcResp.Data,
				}

				bytes, err := json.Marshal(clientResp)
				if err != nil {
					fmt.Println("Invalid message")
				} else {
					conn.WriteMessage(msg.TypeEnum.TextMessage, bytes)
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
		RequestID: index,
		Code:      msg.StatusCode.Fail,
		Data:      data,
	}

	err := respConn.WriteMessage(msg.TypeEnum.TextMessage, rpcResp.ToBytes())
	if err != nil {
		fmt.Println(err)
	}
}
