package connstore

import "github.com/gorilla/websocket"

// ConnInfo 用户连接信息
type ConnInfo struct {
	id   int
	conn *websocket.Conn
	data map[string]interface{}
}

// Get a value from session
func (connInfo ConnInfo) Get(key string) interface{} {
	return connInfo.data[key]
}

// Set a value to session
func (connInfo ConnInfo) Set(key string, v interface{}) {
	connInfo.data[key] = v
}

// ConnMap socket connection map
var ConnMap = make(map[string]*ConnInfo)

// GetConnInfo get connect info
func GetConnInfo(uid string) *ConnInfo {
	connInfo, ok := ConnMap[uid]
	if ok {
		return connInfo
	}
	return nil
}
