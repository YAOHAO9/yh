package connector

import (
	"fmt"
	"net/http"

	"github.com/YAOHAO9/yh/rpc/handler/rpchandler"
	"github.com/YAOHAO9/yh/rpc/msg"
	"github.com/YAOHAO9/yh/rpc/response"

	"github.com/gorilla/websocket"
)

func init() {

	rpchandler.Manager.Register(SysRPCEnum.UpdateSession, func(respCtx *response.RespCtx) {
		// connector.GetConnInfo()

	})

	rpchandler.Manager.Register(SysRPCEnum.PushMessage, func(respCtx *response.RespCtx) {
		connInfo, ok := ConnMap[respCtx.Fm.Session.UID]
		if !ok {
			fmt.Println("无效的Uid", respCtx.Fm.Session.UID, "没有找到对应的客户端连接")
			return
		}
		connInfo.conn.WriteMessage(msg.TypeEnum.TextMessage, respCtx.Fm.ToBytes())
	})

}

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
	uid := r.URL.Query().Get("id")
	conn.SetCloseHandler(func(code int, text string) error {
		delete(ConnMap, uid)
		conn.Close()
		fmt.Println("CloseHandler: ", text)
		return nil
	})

	// 用户认证
	token := r.URL.Query().Get("token")
	fmt.Println("Id: ", uid, " Token: ", token)

	if uid == "" || token == "" {
		fmt.Println("用户校验失败!!!")
		err := conn.WriteMessage(msg.TypeEnum.TextMessage, []byte("认证失败"))
		if err != nil {
			fmt.Println("发送认证失败消息失败: ", err.Error())
		}
		conn.CloseHandler()(0, "认证失败")
		return
	}

	// 防止重复连接
	if oldConnInfo, ok := ConnMap[uid]; ok {
		oldConnInfo.conn.CloseHandler()(0, "关闭重复连接")
	}

	// 保存连接信息
	connInfo := &ConnInfo{uid: uid, conn: conn}
	ConnMap[uid] = connInfo

	connInfo.StartReceiveMsg()

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
