package router

import (
	"yh/rpc/client"
	"yh/rpc/msg"
)

// Map 存储自定义的路由
type Map map[string]func(routerInfo Info, clients []*client.RPCClient) *client.RPCClient

// Info 路由信息
type Info struct {
	ServerKind string
	Handler    string
	Session    msg.Session
}

// Register 注册一个路由函数
func (routeMap Map) Register(serverKind string, route func(routerInfo Info, clients []*client.RPCClient) *client.RPCClient) {
	routeMap[serverKind] = route
}

// Get 获取一个路由函数
func (routeMap Map) Get(serverKind string) func(routerInfo Info, clients []*client.RPCClient) *client.RPCClient {
	return routeMap[serverKind]
}

// Manager 管理router
var Manager = make(Map)
