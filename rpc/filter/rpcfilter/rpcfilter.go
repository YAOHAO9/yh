package rpcfilter

import "trial/rpc/filter"

var beforeRPCFilterMap = make(filter.BeforeFilterSlice, 0)
var rpcAfterFilterMapOfRequest = make(filter.AfterFilterSlice, 0)

// BeforeFilterManager return beforeRPCFilterMap
func BeforeFilterManager() filter.BeforeFilterSlice {
	return beforeRPCFilterMap
}

// AfterFilterManagerOfRequest return rpcAfterFilterMapOfRequest
func AfterFilterManagerOfRequest() filter.AfterFilterSlice {
	return rpcAfterFilterMapOfRequest
}
