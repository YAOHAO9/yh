package filter

import (
	"github.com/YAOHAO9/yh/rpc/message"
	"github.com/YAOHAO9/yh/rpc/response"
)

//===================================================
// Before
//===================================================

// BeforeFilterSlice map
type BeforeFilterSlice []func(respCtx *response.RespCtx) (next bool)

// Register filter
func (slice *BeforeFilterSlice) Register(f func(respCtx *response.RespCtx) (next bool)) {
	*slice = append(*slice, f)
}

// Exec filter
func (slice BeforeFilterSlice) Exec(respCtx *response.RespCtx) (next bool) {
	for _, f := range slice {
		next = f(respCtx)
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
type AfterFilterSlice []func(rm *message.RPCResp) (next bool)

// Register filter
func (slice *AfterFilterSlice) Register(f func(rm *message.RPCResp) (next bool)) {
	*slice = append(*slice, f)
}

// Exec filter
func (slice AfterFilterSlice) Exec(rm *message.RPCResp) (next bool) {
	for _, f := range slice {
		next = f(rm)
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
