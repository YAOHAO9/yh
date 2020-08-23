package channel

import (
	"github.com/YAOHAO9/yh/application"
	"github.com/YAOHAO9/yh/rpc/msg"
)

// Channel ChannelService
type Channel map[string]*msg.Session

// PushMessageToUser 推送消息给指定玩家
func (channel Channel) PushMessageToUser(uid string, data interface{}) {
	session, ok := channel[uid]
	if !ok {
		return
	}

	message := &msg.ClientMessage{
		Handler: "",
		Data:    data,
	}

	application.RPC.Notify.ToServer(session.CID, session, message)
}

// PushMessage 推送消息给所有人
func (channel Channel) PushMessage(uids []string, data interface{}) {
	//
}

// PushMessageToOthers 推送消息给其他人
func (channel Channel) PushMessageToOthers(uids []string, data interface{}) {
	//
}

// Add 推送消息给其他人
func (channel Channel) Add(uid string, session *msg.Session) {
	channel[uid] = session
}
