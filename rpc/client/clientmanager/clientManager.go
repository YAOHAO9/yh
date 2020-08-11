package clientmanager

import (
	"fmt"
	"math/rand"
	"time"
	"trial/rpc/client"
	"trial/rpc/client/router"
	"trial/rpc/config"
	"trial/rpc/msg"
)

var rpcClientMap = make(map[string]*client.RPCClient)

// GetClientByID get rpc client by id
func GetClientByID(id string) (c *client.RPCClient) {
	c, b := rpcClientMap[id]
	if !b {
		return nil
	}
	return
}

// GetClientByRouter get rpc client by rand num
func GetClientByRouter(routerInfo router.Info) (c *client.RPCClient) {

	clients := make([]*client.RPCClient, 0)

	for _, rpcClientInfo := range rpcClientMap {
		if rpcClientInfo.ServerConfig.Kind == routerInfo.ServerKind {
			clients = append(clients, rpcClientInfo)
		}
	}

	if len(clients) == 0 {
		return nil
	}

	route := router.Manager.Get(routerInfo.ServerKind)
	if route != nil {
		return route(routerInfo, clients)
	}

	return clients[rand.Intn(len(clients))]
}

// DelClientByID get rpc client by id
func DelClientByID(id string) {
	delete(rpcClientMap, id)
	return
}

// CreateClient create client
func CreateClient(serverConfig *config.ServerConfig, zkSessionTimeout time.Duration) {
	defer func() {
		data := recover()
		if data != nil {
			fmt.Println("Recover data:", data)
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

// SendMessageByID send message to specified server
func SendMessageByID(serverID string, data []byte) {
	client := GetClientByID(serverID).Conn
	if client != nil {
		client.WriteMessage(msg.TypeEnum.TextMessage, data)
	}
}
