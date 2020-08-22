package connector

import (
	"yh/rpc/handler/syshandler"
	"yh/rpc/response"
)

func init() {

	syshandler.Manager.Register("updateSession", func(respCtx *response.RespCtx) {
		// connector.GetConnInfo()

	})
	syshandler.Manager.Register("pushMessage", func(respCtx *response.RespCtx) {

	})
}
