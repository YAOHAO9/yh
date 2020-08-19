package handlerfilter

import "trial/rpc/filter"

type handlerFilter struct {
	*filter.BaseFilter
}

// Manager HandlerFilter manager
var Manager = &handlerFilter{
	&filter.BaseFilter{
		Before: make(filter.BeforeFilterSlice, 0),
		After:  make(filter.AfterFilterSlice, 0),
	},
}
