package application

import (
	"fmt"

	"github.com/YAOHAO9/yh/rpc"
	"github.com/YAOHAO9/yh/rpc/session"
)

// UpdateSession 注册路由
func UpdateSession(session *session.Session, keys ...string) {

	// 更新session中所有的数据
	if len(keys) == 0 {
		RPC.Notify.ToServer(session.CID, session, rpc.SysRPCEnum.UpdateSession, session.Data)
		return
	}

	// 根据需要更新指定的数据
	data := make(map[string]interface{})
	for _, key := range keys {
		if value, ok := session.Data[key]; ok {
			data[key] = value
		}
	}

	if len(data) == 0 {
		fmt.Println("Update session failed. No such data")
		return
	}

	RPC.Notify.ToServer(session.CID, session, rpc.SysRPCEnum.UpdateSession, data)

}
