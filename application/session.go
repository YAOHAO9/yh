package application

import (
	"encoding/json"

	"github.com/YAOHAO9/pine/connector"
	"github.com/YAOHAO9/pine/rpc"
	"github.com/YAOHAO9/pine/rpc/message"
	"github.com/YAOHAO9/pine/rpc/session"
	"github.com/sirupsen/logrus"
)

// UpdateSession 注册路由
func UpdateSession(session *session.Session, keys ...string) {

	// 更新session中所有的数据
	if len(keys) == 0 {
		bytes, _ := json.Marshal(session.Data)
		rpc.Notify.ToServer(session.CID, session, connector.HandlerMap.UpdateSession, bytes)
		return
	}

	// 根据需要更新指定的数据
	data := make(map[string]interface{})
	for _, key := range keys {
		if value, ok := session.Data[key]; ok {
			data[key] = value
		}
	}

	if len(data) == 0 {
		logrus.Error("Update session failed. Not any data")
		return
	}

	waitChan := make(chan bool, 1)
	bytes, _ := json.Marshal(data)
	rpc.Request.ToServer(session.CID, session, connector.HandlerMap.UpdateSession, bytes, func(rpcResp *message.PineMessage) {
		waitChan <- true
	})

	<-waitChan
}
