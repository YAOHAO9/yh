package channelfactory

import (
	"sync"

	"github.com/YAOHAO9/yh/channel"
)

var mutex sync.Mutex
var channelStore = make(map[string]*channel.Channel)

// CreateChannel 创建一个channel
func CreateChannel(channelID string) *channel.Channel {

	mutex.Lock()
	defer mutex.Unlock()

	channelInstance, ok := channelStore[channelID]
	if ok {
		return channelInstance
	}
	channelIns := &channel.Channel{}
	channelStore[channelID] = channelIns

	return channelIns
}
