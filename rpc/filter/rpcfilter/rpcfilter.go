package rpcfilter

import "yh/rpc/filter"

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
