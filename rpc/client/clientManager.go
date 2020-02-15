package client

import (
	"fmt"
	"math/rand"
	"time"
	"trial/rpc/config"
	"trial/rpc/msgtype"
)

var rpcClientMap = make(map[string]*RPCClient)

// GetClientByID get rpc client by id
func GetClientByID(id string) (c *RPCClient) {
	c, b := rpcClientMap[id]
	if !b {
		return nil
	}
	return
}

// GetRandClientByKind get rpc client by rand num
func GetRandClientByKind(kind string) (c *RPCClient) {

	clients := make([]*RPCClient, 0)

	for _, rpcClientInfo := range rpcClientMap {
		if rpcClientInfo.serverConfig.Kind == kind {
			clients = append(clients, rpcClientInfo)
		}
	}

	if len(clients) == 0 {
		return nil
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
	rpcClient := StartClient(serverConfig, zkSessionTimeout)
	if rpcClient != nil {
		rpcClientMap[serverConfig.ID] = rpcClient
	}
}

// SendMessageByID send message to specified server
func SendMessageByID(serverID string, data []byte) {
	client := GetClientByID(serverID).clientConn
	if client != nil {
		client.WriteMessage(msgtype.TextMessage, data)
	}
}
