package connector

import (
	"fmt"
	"strings"
	"sync"

	"github.com/YAOHAO9/pine/application/config"
	"github.com/YAOHAO9/pine/connector/filter"
	"github.com/YAOHAO9/pine/rpc"
	"github.com/YAOHAO9/pine/rpc/client/clientmanager"
	"github.com/YAOHAO9/pine/rpc/context"
	"github.com/YAOHAO9/pine/rpc/message"
	"github.com/YAOHAO9/pine/rpc/session"
	"github.com/YAOHAO9/pine/service/compressservice"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var mutex sync.Mutex

// Connection 用户连接信息
type Connection struct {
	uid         string
	conn        *websocket.Conn
	data        map[string]string
	routeRecord map[string]string
}

// Get 从session中查找一个值
func (connection Connection) Get(key string) string {
	return connection.data[key]
}

// Set 往session中设置一个键值对
func (connection Connection) Set(key string, v string) {
	connection.data[key] = v
}

// 回复request
func (connection Connection) response(pineMsg *message.PineMsg) {
	bytes, err := proto.Marshal(pineMsg)
	if err != nil {
		logrus.Error(err)
		return
	}
	mutex.Lock()
	err = connection.conn.WriteMessage(message.TypeEnum.BinaryMessage, bytes)
	mutex.Unlock()
	if err != nil {
		logrus.Error(err)
	}
}

// 主动推送消息
func (connection Connection) notify(notify *message.PineMsg) {

	bytes, err := proto.Marshal(notify)
	if err != nil {
		logrus.Error(err)
		return
	}
	newNotify := &message.PineMsg{}
	err = proto.Unmarshal(bytes, newNotify)

	if err != nil {
		logrus.Error(err)
		return
	}

	mutex.Lock()
	err = connection.conn.WriteMessage(message.TypeEnum.BinaryMessage, bytes)
	mutex.Unlock()
	if err != nil {
		logrus.Error(err)
	}
}

// GetSession 获取session
func (connection Connection) GetSession() *session.Session {
	session := &session.Session{
		UID:  connection.uid,
		CID:  config.GetServerConfig().ID,
		Data: connection.data,
	}
	return session
}

// StartReceiveMsg 开始接收消息
func (connection Connection) StartReceiveMsg() {

	conn := connection.conn

	// 开始接收消息
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			conn.CloseHandler()(0, err.Error())
			break
		}
		// 解析消息
		clientMessage := &message.PineMsg{}

		err = proto.Unmarshal(data, clientMessage)

		if err != nil {
			clientMessageResp := &message.PineMsg{
				Route:     "__Error__",
				RequestID: clientMessage.RequestID,
				Data:      []byte(fmt.Sprint("消息解析失败,data:", data, "err:", err)),
			}
			connection.response(clientMessageResp)
			continue
		}

		if clientMessage.Route == "" {
			clientMessageResp := &message.PineMsg{
				Route:     "__Error__",
				RequestID: clientMessage.RequestID,
				Data:      []byte("Route不能为空"),
			}
			connection.response(clientMessageResp)
			continue
		}

		var serverKind string
		var handler string

		routeBytes := []byte(clientMessage.Route)

		if len(routeBytes) == 2 {
			serverKind = compressservice.Server.GetKindByCode(routeBytes[0])
			handler = string(routeBytes[1])
		} else {
			handlerInfos := strings.Split(clientMessage.Route, ".")
			serverKind = handlerInfos[0] // 解析出服务器类型
			handler = handlerInfos[1]    // 真正的handler
		}

		session := connection.GetSession()

		rpcMsg := &message.RPCMsg{
			From:      config.GetServerConfig().ID,
			Handler:   handler,
			Type:      message.RemoterTypeEnum.HANDLER,
			RequestID: clientMessage.RequestID,
			RawData:   clientMessage.Data,
			Session:   session,
		}

		// 获取RPCCLint
		rpcClient := clientmanager.GetClientByRouter(serverKind, rpcMsg, &connection.routeRecord)

		if rpcClient == nil {
			tip := fmt.Sprint("找不到任何", serverKind, "服务器", ", Route: ", clientMessage.Route)

			clientMessageResp := &message.PineMsg{
				Route:     clientMessage.Route,
				RequestID: clientMessage.RequestID,
				Data:      []byte(tip),
			}

			connection.response(clientMessageResp)
			continue
		}

		rpcCtx := context.GenRespCtx(conn, rpcMsg)

		if !filter.Before.Exec(rpcCtx) {
			continue
		}

		if *clientMessage.RequestID == 0 {
			rpc.Notify.ToServer(rpcClient.ServerConfig.ID, rpcMsg)
		} else {
			// 转发Request
			rpc.Request.ToServer(rpcClient.ServerConfig.ID, rpcMsg, func(data []byte) {

				pineMsg := &message.PineMsg{
					RequestID: clientMessage.RequestID,
					Route:     clientMessage.Route,
					Data:      data,
				}

				filter.After.Exec(pineMsg)
				connection.response(pineMsg)
			})
		}
	}
}
