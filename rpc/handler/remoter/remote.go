package remoter

import "github.com/YAOHAO9/pine/rpc/handler"

// Remoter rpc
type Remoter struct {
	*handler.Handler
}

// Manager return RPCHandler
var Manager = &Remoter{
	Handler: &handler.Handler{
		Map: make(handler.Map),
	},
}
