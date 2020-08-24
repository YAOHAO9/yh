package channelfactory

import "github.com/YAOHAO9/yh/rpc/channel"

var channelStore = make(map[string]*channel.Channel)

// CreateChannel 创建一个channel
func CreateChannel(channelID string) *channel.Channel {
	channelInstance, ok := channelStore[channelID]
	if ok {
		return channelInstance
	}
	channelIns := &channel.Channel{}
	channelStore[channelID] = channelIns
	return channelIns
}
