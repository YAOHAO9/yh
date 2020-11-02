package connector

import (
	"net/http"

	"github.com/YAOHAO9/pine/rpc/message"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

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
		logrus.Error("连接失败", err.Error())
		return
	}

	// 断开连接自动清除连接信息
	uid := r.URL.Query().Get("id")
	conn.SetCloseHandler(func(code int, text string) error {

		DelConnection(uid)

		conn.Close()
		logrus.Warn(text)
		return nil
	})

	// Token
	token := r.URL.Query().Get("token")

	// 认证
	data, err := authFunc(uid, token)
	if err != nil || uid == "" {
		err := conn.WriteMessage(message.TypeEnum.TextMessage, []byte("认证失败"))
		if err != nil {
			logrus.Warn("发送认证失败消息失败: ", err.Error())
		}
		conn.CloseHandler()(0, "认证失败")
		return
	}

	// 防止重复连接
	if oldConnection := GetConnection(uid); oldConnection != nil {
		oldConnection.conn.CloseHandler()(0, "关闭重复连接")
	}

	// 保存连接信息
	connection := &Connection{
		uid:         uid,
		conn:        conn,
		data:        data,
		routeRecord: make(map[string]string),
	}

	SaveConnection(connection)

	connection.StartReceiveMsg()

}
