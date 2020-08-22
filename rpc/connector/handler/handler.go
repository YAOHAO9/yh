package connector

import (
	"github.com/YAOHAO9/yh/rpc/handler/syshandler"
	"github.com/YAOHAO9/yh/rpc/response"
)

func init() {

	syshandler.Manager.Register("updateSession", func(respCtx *response.RespCtx) {
		// connector.GetConnInfo()

	})
	syshandler.Manager.Register("pushMessage", func(respCtx *response.RespCtx) {

	})
}
