package main

import (
	"errors"
	"fmt"
	_ "net/http/pprof"
	"time"

	"github.com/YAOHAO9/pine/application"
	"github.com/YAOHAO9/pine/channel/channelfactory"
	"github.com/YAOHAO9/pine/rpc"
	"github.com/YAOHAO9/pine/rpc/client"
	"github.com/YAOHAO9/pine/rpc/context"
	"github.com/YAOHAO9/pine/rpc/handler"
	"github.com/YAOHAO9/pine/rpc/message"
)

func main() {

	// slice
	slice := []interface{}{"hello", "world", "hello", "everyone!"}
	for index, value := range slice {
		fmt.Printf("slice %d is: %s\n", index, value)
	}

	app := application.CreateApp()

	app.AsConnector(func(uid string, token string, sessionData map[string]interface{}) error {

		if uid == "" || token == "" {
			return errors.New("Invalid token")
		}
		sessionData[token] = token

		return nil
	})

	app.RegisteHandler("handler", func(rpcCtx *context.RPCCtx) *handler.Resp {

		channelInstance := channelfactory.CreateChannel("101")
		channelInstance.Add(rpcCtx.Session.UID, rpcCtx.Session)

		channelInstance.PushMessage("onMsg", "PushMessage")                                               // 推送给所有在当前channel中的玩家
		channelInstance.PushMessageToOthers([]string{rpcCtx.Session.UID}, "onMsg", "PushMessageToOthers") // 推送给除了切片内的channel中的玩家
		channelInstance.PushMessageToUser(rpcCtx.Session.UID, "onMsg", "PushMessageToUser")               // 只推送给当前玩家
		channelInstance.PushMessageToUsers([]string{rpcCtx.Session.UID}, "onMsg", "PushMessageToUsers")   // 只推送给切片的指定的玩家

		// rpc.Request.ToServer(serverID string, session *session.Session, handler string, data interface{}, f func(rpcResp *message.RPCResp))
		// rpc.Request.ByKind(serverKind string, session *session.Session, handler string, data interface{}, f func(rpcResp *message.RPCResp))
		// rpc.Notify.ToServer(serverID string, session *session.Session, handler string, data interface{})
		// rpc.Notify.ByKind(serverKind string, session *session.Session, handler string, data interface{})

		rpc.Request.ByKind("connector", nil, "getOneRobot", nil, func(rpcResp *message.RPCResp) {
			fmt.Println("收到Rpc的回复：", rpcResp.Data)
		})

		return nil
	})

	app.RegisteRemoter("getOneRobot", func(rpcCtx *context.RPCCtx) *handler.Resp {

		return &handler.Resp{
			Data: map[string]interface{}{
				"name": "盖伦",
				"sex":  1,
				"age":  18,
			},
		}
	})

	app.RegisteHandlerBeforeFilter(func(rpcCtx *context.RPCCtx) (next bool) {

		lastEnterRoomTimeInterface := rpcCtx.Session.Data["lastEnterRoomTime"]
		if lastEnterRoomTimeInterface != nil {
			if lastEnterRoomTime, ok := lastEnterRoomTimeInterface.(time.Time); ok && lastEnterRoomTime.Sub(time.Now()) > time.Second {
				rpcCtx.SendMsg("操作太频繁", message.StatusCode.Fail) // 返回结果
				return false                                     // 停止执行下个before filter以及hanler
			}
		}

		rpcCtx.Session.Set("lastEnterRoomTime", time.Now())
		application.UpdateSession(rpcCtx.Session, "lastEnterRoomTime")
		return true // 继续执行下个before filter直到执行handler
	})

	app.RegisteHandlerAfterFilter(func(rpcResp *message.RPCResp) (next bool) {

		// 修改pine定义的错误码
		if rpcResp.Code == message.StatusCode.Fail {
			rpcResp.Code = 400 // 自定义错误码
		} else if rpcResp.Code == message.StatusCode.Successful {
			rpcResp.Code = 200 // 自定义成功码
		}

		return true // return true继续执行下个after filter, return false停止执行下个after filter
	})

	app.RegisteRouter("ddz", func(rpcMsg *message.RPCMsg, clients []*client.RPCClient) *client.RPCClient {

		for _, clientInfo := range clients {
			if chatServerID, ok := rpcMsg.Session.Get("chatServerID").(string); ok && clientInfo.ServerConfig.ID == chatServerID {
				return clientInfo
			}
		}

		return nil //if return nil, pine will get one rpc client by random
	})

	app.Start()
}
