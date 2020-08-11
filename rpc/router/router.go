package router

import "trial/rpc/msg"

// Map 存储自定义的路由
type Map map[string]func(datas ...interface{})

// Info 路由信息
type Info struct {
	to      string
	handler string
	session msg.Session
}

func (routeMap Map) add(serverKind string, route func(datas ...interface{})) {
	routeMap[serverKind] = route
}

func (routeMap Map) get(serverKind string) func(datas ...interface{}) {
	return routeMap[serverKind]
}

// Manager 管理router
var Manager = make(Map)
