package connector

import (
	"fmt"

	"github.com/YAOHAO9/yh/rpc"
	"github.com/YAOHAO9/yh/rpc/context"
	"github.com/YAOHAO9/yh/rpc/handler"
	"github.com/YAOHAO9/yh/rpc/handler/rpchandler"
)

func init() {
	// 更新Session
	rpchandler.Manager.Register(rpc.SysRPCEnum.UpdateSession, func(rpcCtx *context.RPCCtx) *handler.Resp {
		connInfo, ok := ConnMap[rpcCtx.Session.UID]
		if !ok {
			fmt.Println("无效的Uid", rpcCtx.Session.UID, "没有找到对应的客户端连接")
			return nil
		}

		if data, ok := rpcCtx.Data.(map[string]interface{}); ok {
			for key, value := range data {
				connInfo.data[key] = value
			}
		}

		return nil
	})

	// 推送消息
	rpchandler.Manager.Register(rpc.SysRPCEnum.PushMessage, func(rpcCtx *context.RPCCtx) *handler.Resp {
		connInfo, ok := ConnMap[rpcCtx.Session.UID]
		if !ok {
			fmt.Println("无效的Uid", rpcCtx.Session.UID, "没有找到对应的客户端连接")
			return nil
		}

		if notify, ok := rpcCtx.Data.(map[string]interface{}); ok {
			connInfo.notify(notify["Route"].(string), notify["Data"])
		}

		return nil
	})

}
