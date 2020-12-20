package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	_ "net/http/pprof"
	"strconv"
	"time"

	"github.com/YAOHAO9/pine/application"
	"github.com/YAOHAO9/pine/handlermessage"
	"github.com/YAOHAO9/pine/rpc"
	"github.com/YAOHAO9/pine/rpc/client"
	"github.com/YAOHAO9/pine/rpc/context"
	"github.com/YAOHAO9/pine/rpc/message"
	"github.com/YAOHAO9/pine/rpc/session"
	"github.com/YAOHAO9/pine/service/channelservice"
	"github.com/YAOHAO9/pine/service/compressservice"
	"github.com/YAOHAO9/pine/service/sessionservice"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
)

func main() {

	app := application.CreateApp()

	compressservice.Event.AddEventCompressRecords("onMsg1", "onMsgJSON") // 需要压缩的Event

	app.AsConnector(func(uid string, token string, sessionData map[string]string) error {

		if uid == "" || token == "" {
			return errors.New("Invalid token")
		}
		sessionData[token] = token

		return nil
	})

	app.RegisteHandler("handler", func(rpcCtx *context.RPCCtx, data *handlermessage.Handler) {

		channel := channelservice.CreateChannel("101")
		channel.Add(rpcCtx.Session.UID, rpcCtx.Session)

		bytes, _ := json.Marshal(map[string]string{
			"Name": "onMsg",
			"Data": "啊哈哈傻法师打上发发",
		})

		channel.PushMessage("onMsg1", bytes)                                       // 推送给所有在当前channel中的玩家
		channel.PushMessageToOthers([]string{rpcCtx.Session.UID}, "onMsg2", bytes) // 推送给除了切片内的channel中的玩家
		channel.PushMessageToUser(rpcCtx.Session.UID, "onMsg3", bytes)             // 只推送给当前玩家
		channel.PushMessageToUsers([]string{rpcCtx.Session.UID}, "onMsg4", bytes)  // 只推送给切片的指定的玩家

		logrus.Warn(fmt.Sprintf("%#v", data))

		rpcMsg := &message.RPCMsg{
			Handler: "getOneRobot",
			RawData: []byte{},
		}

		rpc.Request.ByKind("connector", rpcMsg, func(data []byte) {
			fmt.Println("收到Rpc的回复：", string(data))
		})

		handlerResp := &handlermessage.HandlerResp{
			Code:    1,
			Name:    "HandlerResp",
			Message: "HandlerResp Message",
		}

		rpcCtx.Session.Set("hhhaaa", "有点一次=====")
		sessionservice.UpdateSession(rpcCtx.Session)
		sessionservice.GetSession(rpcCtx.Session.CID, rpcCtx.Session.UID, func(session *session.Session) {
			logrus.Debug(session.Get("hhhaaa"))
		})

		bytes, _ = proto.Marshal(handlerResp)

		rpcCtx.SendMsg(bytes)

		sessionservice.KickBySession(rpcCtx.Session, "====啊师傅打死====")
	})

	app.RegisteHandler("handlerJSON", func(rpcCtx *context.RPCCtx, data map[string]interface{}) {

		channel := channelservice.CreateChannel("101")
		channel.Add(rpcCtx.Session.UID, rpcCtx.Session)

		bytes, _ := json.Marshal(map[string]interface{}{
			"Route":     "onMsgJSON",
			"heiheihie": "heiheihei",
		})

		channel.PushMessage("onMsgJSON", bytes)                                       // 推送给所有在当前channel中的玩家
		channel.PushMessageToOthers([]string{rpcCtx.Session.UID}, "onMsgJSON", bytes) // 推送给除了切片内的channel中的玩家
		channel.PushMessageToUser(rpcCtx.Session.UID, "onMsgJSON", bytes)             // 只推送给当前玩家
		channel.PushMessageToUsers([]string{rpcCtx.Session.UID}, "onMsgJSON", bytes)  // 只推送给切片的指定的玩家

		logrus.Warn(fmt.Sprintf("%#v", data))
		rpcMsg := &message.RPCMsg{
			Handler: "getOneRobot",
		}
		rpc.Request.ByKind("connector", rpcMsg, func(data []byte) {
			fmt.Println("收到Rpc的回复：", string(data))
		})

		bytes, _ = json.Marshal(map[string]interface{}{
			"Route":      "handlerResponseJSON",
			"hahahahah ": 122222222222222,
		})

		rpcCtx.SendMsg(bytes)
	})

	app.RegisteRemoter("getOneRobot", func(rpcCtx *context.RPCCtx, data interface{}) {
		rpcCtx.SendMsg([]byte(fmt.Sprint(rand.Int())))
	})

	app.RegisteHandlerBeforeFilter(func(rpcCtx *context.RPCCtx) (next bool) {

		if rpcCtx.GetHandler() == "enterRoom" {
			lastEnterRoomTimeInterface := rpcCtx.Session.Data["lastEnterRoomTime"]
			if lastEnterRoomTimeInterface != "" {
				timestamp, e := strconv.ParseInt(lastEnterRoomTimeInterface, 10, 64)
				if e != nil {
					logrus.Error("不能将", lastEnterRoomTimeInterface, "转换成时间戳")
				} else if time.Now().Sub(time.Unix(timestamp, 0)) < time.Second {
					rpcCtx.SendMsg([]byte("操作太频繁")) // 返回结果
					return false                    // 停止执行下个before filter以及hanler
				}
			}

			rpcCtx.Session.Set("lastEnterRoomTime", fmt.Sprint(time.Now().Unix()))
			sessionservice.UpdateSession(rpcCtx.Session, "lastEnterRoomTime")
		}
		return true // 继续执行下个before filter直到执行handler
	})

	app.RegisteHandlerAfterFilter(func(rpcResp *message.PineMsg) (next bool) {

		// 修改pine定义的错误码
		// if rpcResp.Code == message.StatusCode.Fail {
		// 	rpcResp.Code = 400 // 自定义错误码
		// } else if rpcResp.Code == message.StatusCode.Successful {
		// 	rpcResp.Code = 200 // 自定义成功码
		// }

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

func aa(a string) {
	fmt.Println(a)
}
