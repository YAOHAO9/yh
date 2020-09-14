package filter

import (
	"github.com/YAOHAO9/yh/rpc/context"
	"github.com/YAOHAO9/yh/rpc/message"
)

//===================================================
// Before
//===================================================

// BeforeFilterSlice map
type BeforeFilterSlice []func(rpcCtx *context.RPCCtx) (next bool)

// Register filter
func (slice *BeforeFilterSlice) Register(f func(rpcCtx *context.RPCCtx) (next bool)) {
	*slice = append(*slice, f)
}

// Exec filter(返回true标识继续往下执行)
func (slice BeforeFilterSlice) Exec(rpcCtx *context.RPCCtx) (next bool) {
	for _, f := range slice {
		next = f(rpcCtx)
		if !next {
			return false
		}
	}
	return true
}

//===================================================
// After
//===================================================

// AfterFilterSlice map
type AfterFilterSlice []func(rpcResp *message.RPCResp) (next bool)

// Register filter
func (slice *AfterFilterSlice) Register(f func(rpcResp *message.RPCResp) (next bool)) {
	*slice = append(*slice, f)
}

// Exec filter
func (slice AfterFilterSlice) Exec(rpcResp *message.RPCResp) (next bool) {
	for _, f := range slice {
		next = f(rpcResp)
		if !next {
			return false
		}
	}
	return true
}

// BaseFilter baseFilter
type BaseFilter struct {
	Before BeforeFilterSlice
	After  AfterFilterSlice
}
