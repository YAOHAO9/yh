package connector

import (
	"sync"

	"github.com/YAOHAO9/pine/application/config"
	"github.com/YAOHAO9/pine/logger"
	"github.com/YAOHAO9/pine/rpc/message"
	"github.com/YAOHAO9/pine/rpc/session"
	"github.com/YAOHAO9/pine/serializer"
)

// ConnInfo 用户连接信息
type ConnInfo struct {
	uid            string
	conn           ConnectionInterface
	data           map[string]string
	routeRecord    map[string]string
	compressRecord map[string]bool
	mutex          sync.Mutex
}

// Get 从session中查找一个值
func (connInfo *ConnInfo) Get(key string) string {
	return connInfo.data[key]
}

// Set 往session中设置一个键值对
func (connInfo *ConnInfo) Set(key string, v string) {
	connInfo.data[key] = v
}

// 回复request
func (connInfo *ConnInfo) response(pineMsg *message.PineMsg) {
	connInfo.mutex.Lock()
	defer connInfo.mutex.Unlock()
	err := connInfo.conn.SendMsg(serializer.ToBytes(pineMsg))

	if err != nil {
		logger.Error(err)
	}
}

// 主动推送消息
func (connInfo *ConnInfo) notify(notify *message.PineMsg) {

	connInfo.mutex.Lock()
	defer connInfo.mutex.Unlock()

	err := connInfo.conn.SendMsg(serializer.ToBytes(notify))

	if err != nil {
		logger.Error(err)
	}
}

// GetSession 获取session
func (connInfo *ConnInfo) GetSession() *session.Session {
	session := &session.Session{
		UID:  connInfo.uid,
		CID:  config.GetServerConfig().ID,
		Data: connInfo.data,
	}
	return session
}
