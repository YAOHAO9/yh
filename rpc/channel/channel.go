package channel

import (
	"github.com/YAOHAO9/yh/application"
	"github.com/YAOHAO9/yh/rpc/msg"
	"github.com/YAOHAO9/yh/util/beeku"
)

// Channel ChannelService
type Channel map[string]*msg.Session

// PushMessageToUser 推送消息给指定玩家
func (channel Channel) PushMessageToUser(uid string, data interface{}) {
	session, ok := channel[uid]
	if !ok {
		return
	}

	application.RPC.Notify.ToServer(session.CID, "", session, data)
}

// PushMessage 推送消息给所有人
func (channel Channel) PushMessage(uids []string, data interface{}) {
	for _, uid := range uids {
		channel.PushMessageToUser(uid, data)
	}
}

// PushMessageToOthers 推送消息给其他人
func (channel Channel) PushMessageToOthers(uids []string, data interface{}) {
	for _, uid := range uids {
		if beeku.InSlice(uid, uids) == -1 {
			channel.PushMessageToUser(uid, data)
		}
	}
}

// Add 推送消息给其他人
func (channel Channel) Add(uid string, session *msg.Session) {
	channel[uid] = session
}
