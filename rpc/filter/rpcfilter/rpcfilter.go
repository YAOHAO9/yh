package rpcfilter

import "trial/rpc/filter"

var beforeRPCFilterMap = make(filter.BeforeFilterSlice, 0)
var rpcAfterFilterMap = make(filter.AfterFilterSlice, 0)

// BeforeFilterManager return beforeRPCFilterMap
func BeforeFilterManager() *filter.BeforeFilterSlice {
	return &beforeRPCFilterMap
}

// AfterFilterManager return rpcAfterFilterMap
func AfterFilterManager() *filter.AfterFilterSlice {
	return &rpcAfterFilterMap
}
