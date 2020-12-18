package channelservice

import (
	"sync"

	"github.com/YAOHAO9/pine/application/config"
	"github.com/YAOHAO9/pine/channelservice/eventcompress"
	"github.com/YAOHAO9/pine/connector"
	"github.com/YAOHAO9/pine/rpc"
	"github.com/YAOHAO9/pine/rpc/message"
	"github.com/YAOHAO9/pine/rpc/session"
	"github.com/golang/protobuf/proto"
)

var lock sync.RWMutex

// Channel ChannelService
type Channel map[string]*session.Session

// PushMessage 推送消息给所有人
func (channel Channel) PushMessage(event string, data []byte) {

	lock.RLock()
	defer lock.RUnlock()

	for _, session := range channel {
		PushMessageBySession(session, event, data)
	}
}

// PushMessageToOthers 推送消息给其他人
func (channel Channel) PushMessageToOthers(uids []string, event string, data []byte) {

	lock.RLock()
	defer lock.RUnlock()

	for uid := range channel {
		findIndex := -1
		for index, value := range uids {
			if uid == value {
				findIndex = index
				break
			}
		}
		if findIndex == -1 {
			channel.PushMessageToUser(uid, event, data)
		}
	}
}

// PushMessageToUsers 推送消息给指定玩家
func (channel Channel) PushMessageToUsers(uids []string, event string, data []byte) {

	for _, uid := range uids {
		channel.PushMessageToUser(uid, event, data)
	}

}

// PushMessageToUser 推送消息给指定玩家
func (channel Channel) PushMessageToUser(uid string, event string, data []byte) {

	lock.RLock()
	defer lock.RUnlock()

	session, ok := channel[uid]
	if !ok {
		return
	}

	PushMessageBySession(session, event, data)

}

// Add 推送消息给其他人
func (channel Channel) Add(uid string, session *session.Session) {

	lock.Lock()
	defer lock.Unlock()

	channel[uid] = session
}

// PushMessageBySession 通过session推送消息
func PushMessageBySession(session *session.Session, event string, data []byte) {

	code := eventcompress.GetCodeByEvent(event)

	var notify *message.PineMsg
	if code != 0 {
		notify = &message.PineMsg{
			Route: string(code),
			Data:  data,
		}
	} else {
		notify = &message.PineMsg{
			Route: config.GetServerConfig().Kind + "." + event,
			Data:  data,
		}
	}

	bytes, _ := proto.Marshal(notify)
	rpcMsg := &message.RPCMsg{
		Handler: connector.HandlerMap.PushMessage,
		RawData: bytes,
		Session: session,
	}
	rpc.Notify.ToServer(session.CID, rpcMsg)
}
