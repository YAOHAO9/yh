package channel

import (
	"sync"

	"github.com/YAOHAO9/yh/connector"
	"github.com/YAOHAO9/yh/rpc"
	"github.com/YAOHAO9/yh/rpc/message"
	"github.com/YAOHAO9/yh/rpc/session"
	"github.com/YAOHAO9/yh/util/beeku"
)

var lock sync.RWMutex

// Channel ChannelService
type Channel map[string]*session.Session

// PushMessage 推送消息给所有人
func (channel Channel) PushMessage(route string, data interface{}) {

	lock.RLock()
	defer lock.RUnlock()

	for _, session := range channel {
		PushMessageBySession(session, route, data)
	}
}

// PushMessageToOthers 推送消息给其他人
func (channel Channel) PushMessageToOthers(uids []string, route string, data interface{}) {

	lock.RLock()
	defer lock.RUnlock()

	for uid := range channel {
		if beeku.InSlice(uid, uids) == -1 {
			channel.PushMessageToUser(uid, route, data)
		}
	}
}

// PushMessageToUsers 推送消息给指定玩家
func (channel Channel) PushMessageToUsers(uids []string, route string, data interface{}) {

	for _, uid := range uids {
		channel.PushMessageToUser(uid, route, data)
	}

}

// PushMessageToUser 推送消息给指定玩家
func (channel Channel) PushMessageToUser(uid string, route string, data interface{}) {

	lock.RLock()
	defer lock.RUnlock()

	session, ok := channel[uid]
	if !ok {
		return
	}

	PushMessageBySession(session, route, data)

}

// Add 推送消息给其他人
func (channel Channel) Add(uid string, session *session.Session) {

	lock.Lock()
	defer lock.Unlock()

	channel[uid] = session
}

// PushMessageBySession 通过session推送消息
func PushMessageBySession(session *session.Session, route string, data interface{}) {
	notify := message.RPCNotify{
		Route: route,
		Data:  data,
	}
	rpc.Notify.ToServer(session.CID, session, connector.HandlerMap.PushMessage, notify)
}
