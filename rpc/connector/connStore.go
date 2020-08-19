package connector

import "github.com/gorilla/websocket"

// ConnInfo 用户连接信息
type ConnInfo struct {
	id   int
	conn *websocket.Conn
	data map[string]interface{}
}

// Get 从session中查找一个值
func (connInfo ConnInfo) Get(key string) interface{} {
	return connInfo.data[key]
}

// Set 往session中设置一个键值对
func (connInfo ConnInfo) Set(key string, v interface{}) {
	connInfo.data[key] = v
}

// ConnMap 连接信息Map
var ConnMap = make(map[string]*ConnInfo)

// GetConnInfo 获取连接信息
func GetConnInfo(uid string) *ConnInfo {
	connInfo, ok := ConnMap[uid]
	if ok {
		return connInfo
	}
	return nil
}
