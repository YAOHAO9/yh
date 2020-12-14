package connector

import (
	"fmt"
	"strings"
	"sync"

	"github.com/YAOHAO9/pine/application/config"
	"github.com/YAOHAO9/pine/connector/filter"
	"github.com/YAOHAO9/pine/connector/serverdict"
	"github.com/YAOHAO9/pine/rpc"
	"github.com/YAOHAO9/pine/rpc/client/clientmanager"
	"github.com/YAOHAO9/pine/rpc/context"
	"github.com/YAOHAO9/pine/rpc/message"
	"github.com/YAOHAO9/pine/rpc/session"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// HandlerPrefix  Handler 前缀
var HandlerPrefix = ""

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
func (connection Connection) response(clientMsgResp *message.PineMessage) {
	bytes, err := proto.Marshal(clientMsgResp)
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
func (connection Connection) notify(notify *message.PineMessage) {

	bytes, err := proto.Marshal(notify)
	if err != nil {
		logrus.Error(err)
		return
	}
	newNotify := &message.PineMessage{}
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

// StartReceiveMsg 开始接收消息
func (connection Connection) StartReceiveMsg() {

	uid := connection.uid
	conn := connection.conn

	// 开始接收消息
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			conn.CloseHandler()(0, err.Error())
			break
		}
		// 解析消息
		clientMessage := &message.PineMessage{}

		err = proto.Unmarshal(data, clientMessage)

		if err != nil {
			clientMessageResp := &message.PineMessage{
				Route:     "__Error__",
				RequestID: clientMessage.RequestID,
				Data:      []byte(fmt.Sprint("消息解析失败,data:", data, "err:", err)),
			}
			connection.response(clientMessageResp)
			continue
		}

		if clientMessage.Route == "" {
			clientMessageResp := &message.PineMessage{
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
			serverKind = serverdict.GetKindByCode(routeBytes[0])
			handler = string(routeBytes[1])
		} else {
			handlerInfos := strings.Split(clientMessage.Route, ".")
			serverKind = handlerInfos[0] // 解析出服务器类型
			handler = handlerInfos[1]    // 真正的handler
		}

		session := &session.Session{
			UID:  uid,
			CID:  config.GetServerConfig().ID,
			Data: connection.data,
		}

		rpcMsg := &message.RPCMsg{
			From:      config.GetServerConfig().ID,
			Handler:   handler,
			RequestID: *clientMessage.RequestID,
			RawData:   clientMessage.Data,
			Session:   session,
		}

		// 获取RPCCLint
		rpcClient := clientmanager.GetClientByRouter(serverKind, rpcMsg, &connection.routeRecord)

		if rpcClient == nil {
			tip := fmt.Sprint("找不到任何", serverKind, "服务器", ", Route: ", clientMessage.Route)

			clientMessageResp := &message.PineMessage{
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
			rpc.Notify.ToServer(rpcClient.ServerConfig.ID, session, HandlerPrefix+handler, clientMessage.Data)
		} else {
			// 转发Request
			rpc.Request.ToServer(rpcClient.ServerConfig.ID, session, HandlerPrefix+handler, clientMessage.Data, func(rpcResp *message.PineMessage) {
				rpcResp.RequestID = clientMessage.RequestID
				rpcResp.Route = clientMessage.Route
				filter.After.Exec(rpcResp)
				connection.response(rpcResp)
			})
		}
	}
}
