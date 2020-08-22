package rpcfilter

import "github.com/YAOHAO9/yh/rpc/filter"

type rpcFilter struct {
	*filter.BaseFilter
}

// Manager HandlerFilter manager
var Manager = &rpcFilter{
	&filter.BaseFilter{
		Before: make(filter.BeforeFilterSlice, 0),
		After:  make(filter.AfterFilterSlice, 0),
	},
}
