package connector

import (
	"trial/rpc/handler/syshandler"
	"trial/rpc/response"
)

func init() {
	syshandler.Manager.Register("updateSession", func(respCtx *response.RespCtx) {
		// connector.GetConnInfo()

	})
	syshandler.Manager.Register("pushMessage", func(respCtx *response.RespCtx) {

	})
}
