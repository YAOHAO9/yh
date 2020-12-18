package sessionservice

import (
	"encoding/json"

	"github.com/YAOHAO9/pine/connector"
	"github.com/YAOHAO9/pine/rpc"
	"github.com/YAOHAO9/pine/rpc/message"
	"github.com/YAOHAO9/pine/rpc/session"
	"github.com/YAOHAO9/pine/util"
	"github.com/sirupsen/logrus"
)

// UpdateSession 注册路由
func UpdateSession(session *session.Session, keys ...string) {

	// 更新session中所有的数据
	if len(keys) == 0 {
		bytes, _ := json.Marshal(session.Data)
		rpcMsg := &message.RPCMsg{
			Session: session,
			Handler: connector.HandlerMap.UpdateSession,
			RawData: bytes,
		}
		rpc.Notify.ToServer(session.CID, rpcMsg)
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

	bytes, _ := json.Marshal(data)
	rpcMsg := &message.RPCMsg{
		Session: session,
		Handler: connector.HandlerMap.UpdateSession,
		RawData: bytes,
	}
	rpc.Notify.ToServer(session.CID, rpcMsg)

}

// GetSession 获取session
func GetSession(CID, UID string, f func(session *session.Session)) {

	data := map[string]string{
		"UID": UID,
		"CID": CID,
	}
	rpcMsg := &message.RPCMsg{
		Handler: connector.HandlerMap.GetSession,
		RawData: util.ToBytes(data),
	}
	rpc.Request.ToServer(CID, rpcMsg, f)
	return
}
