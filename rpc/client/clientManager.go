package client

import (
	"github.com/gorilla/websocket"
	"trial/config"
	"trial/rpc/msgtype"
)

// RPCClientInfo websocket client 连接信息
type RPCClientInfo struct {
	clientConn   *websocket.Conn
	serverConfig *config.ServerConfig
}

var rpcClientMap = make(map[string]*RPCClientInfo)

// GetClientByID get rpc client by id
func GetClientByID(id string) (c *RPCClientInfo) {
	c, b := rpcClientMap[id]
	if !b {
		return nil
	}
	return
}

// DelClientByID get rpc client by id
func DelClientByID(id string) {
	delete(rpcClientMap, id)
	return
}

// CreateClient create client
func CreateClient(serverConfig *config.ServerConfig) {
	clientConn := Start(serverConfig)
	rpcClientMap[serverConfig.ID] = &RPCClientInfo{clientConn: clientConn, serverConfig: serverConfig}
}

// SendMessageByID send message to specified server
func SendMessageByID(serverID string, data []byte) {
	client := GetClientByID(serverID).clientConn
	if client != nil {
		client.WriteMessage(msgtype.TextMessage, data)
	}
}
