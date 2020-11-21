package clientmanager

import (
	"math/rand"
	"time"

	"github.com/YAOHAO9/pine/application/config"
	"github.com/YAOHAO9/pine/rpc/client"
	"github.com/YAOHAO9/pine/rpc/message"
	"github.com/YAOHAO9/pine/rpc/router"
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

// GetClientByRouter 通过路由后去一个客户端的Rpc连接
func GetClientByRouter(serverKind string, rpcMsg *message.RPCMsg, routeRecord *map[string]string) (rpcClient *client.RPCClient) {

	defer func() {
		if rpcClient != nil && routeRecord != nil {
			(*routeRecord)[rpcClient.ServerConfig.Kind] = rpcClient.ServerConfig.ID
		}
	}()
	clients := GetClientsByKind(serverKind)

	if len(clients) == 0 {
		return nil
	}

	route := router.Manager.Get(serverKind)
	if route != nil {
		client := route(rpcMsg, clients)
		return client
	}

	route = router.Manager.Get("*")

	if route != nil {
		client := route(rpcMsg, clients)
		return client
	}

	if routeRecord != nil {
		if clientID, ok := (*routeRecord)[serverKind]; ok {
			client := GetClientByID(clientID)
			return client
		}
	}

	client := clients[rand.Intn(len(clients))]
	return client
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
