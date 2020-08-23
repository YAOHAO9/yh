package connector

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
