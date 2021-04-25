package connector

import (
	"github.com/YAOHAO9/pine/application/config"
	"github.com/YAOHAO9/pine/rpc/message"
	"github.com/YAOHAO9/pine/service/compressservice"
)

var connInfoStore = make(map[string]*ConnInfo)

// SaveConnInfo 保存连接
func SaveConnInfo(connInfo *ConnInfo) {
	connInfoStore[connInfo.uid] = connInfo
}

// GetConnInfo 获取连接
func GetConnInfo(uid string) *ConnInfo {
	connInfo, ok := connInfoStore[uid]
	if ok {
		return connInfo
	}
	return nil
}

// DelConnInfo 删除连接
func DelConnInfo(uid string) {
	delete(connInfoStore, uid)
}

// DelConnInfo 删除连接
func KickByUid(uid string, data []byte) {
	connInfo := GetConnInfo(uid)
	if connInfo == nil {
		return
	}
	notify := &message.PineMsg{
		Route: string([]byte{
			compressservice.Server.GetCodeByKind(config.GetServerConfig().Kind),
			compressservice.Event.GetCodeByEvent(ConnectorHandlerMap.Kick)}),
		Data: data,
	}
	connInfo.notify(notify)
	DelConnInfo(connInfo.uid)
	connInfo.conn.Close()
}
