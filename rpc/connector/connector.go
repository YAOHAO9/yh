package connector

import (
	"fmt"
	"net/http"

	"github.com/YAOHAO9/yh/rpc"
	"github.com/YAOHAO9/yh/rpc/context"
	"github.com/YAOHAO9/yh/rpc/handler/rpchandler"
	"github.com/YAOHAO9/yh/rpc/message"

	"github.com/gorilla/websocket"
)

func init() {

	// 更新Session
	rpchandler.Manager.Register(rpc.SysRPCEnum.UpdateSession, func(rpcCtx *context.RPCCtx) {
		connInfo, ok := ConnMap[rpcCtx.Session.UID]
		if !ok {
			fmt.Println("无效的Uid", rpcCtx.Session.UID, "没有找到对应的客户端连接")
			return
		}

		if data, ok := rpcCtx.Data.(map[string]interface{}); ok {
			for key, value := range data {
				connInfo.data[key] = value
			}
		}
	})

	// 推送消息
	rpchandler.Manager.Register(rpc.SysRPCEnum.PushMessage, func(rpcCtx *context.RPCCtx) {
		connInfo, ok := ConnMap[rpcCtx.Session.UID]
		if !ok {
			fmt.Println("无效的Uid", rpcCtx.Session.UID, "没有找到对应的客户端连接")
			return
		}

		if notify, ok := rpcCtx.Data.(map[string]interface{}); ok {
			connInfo.notify(notify["Route"].(string), notify["Data"])
		}

	})

}

var upgrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketHandler 处理ws请求
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
		err := conn.WriteMessage(message.TypeEnum.TextMessage, []byte("认证失败"))
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
	connInfo := &ConnInfo{uid: uid, conn: conn, data: make(map[string]interface{})}
	ConnMap[uid] = connInfo

	connInfo.StartReceiveMsg()

}
