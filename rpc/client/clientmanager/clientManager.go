package clientmanager

import (
	"math/rand"
	"time"

	"github.com/YAOHAO9/yh/application/config"
	"github.com/YAOHAO9/yh/rpc/client"
	"github.com/YAOHAO9/yh/rpc/message"
	"github.com/YAOHAO9/yh/rpc/router"
)

var rpcClientMap = make(map[string]*client.RPCClient)

// GetClientByID 通过ID获取Rpc连接客户端
func GetClientByID(id string) (c *client.RPCClient) {
	c, b := rpcClientMap[id]
	if !b {
		return nil
	}
	return
}

// GetClientsByKind 根据服务器类型获取RPC连接客户端
func GetClientsByKind(serverKind string) (c []*client.RPCClient) {

	clients := make([]*client.RPCClient, 0)

	for _, rpcClientInfo := range rpcClientMap {
		if rpcClientInfo.ServerConfig.Kind == serverKind {
			clients = append(clients, rpcClientInfo)
		}
	}
	return clients
}

// GetClientByRouter 随机获取一个Rpc连接客户端
func GetClientByRouter(serverKind string, rpcMsg *message.RPCMsg) (c *client.RPCClient) {

	clients := GetClientsByKind(serverKind)

	if len(clients) == 0 {
		return nil
	}

	route := router.Manager.Get(serverKind)
	if route != nil {
		return route(rpcMsg, clients)
	}

	route = router.Manager.Get("*")

	if route != nil {
		return route(rpcMsg, clients)
	}

	return clients[rand.Intn(len(clients))]
}

// DelClientByID 删除RPC连接客户端
func DelClientByID(id string) {
	delete(rpcClientMap, id)
	return
}

// CreateClient 创建RPC连接客户端
func CreateClient(serverConfig *config.ServerConfig, zkSessionTimeout time.Duration) {
	defer func() {
		data := recover()
		if data != nil {
			delete(rpcClientMap, serverConfig.ID)
		}
	}()
	rpcClient := client.StartClient(serverConfig, zkSessionTimeout, func(id string) {
		DelClientByID(id)
	})
	if rpcClient != nil {
		rpcClientMap[serverConfig.ID] = rpcClient
	}
}
