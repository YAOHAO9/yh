package filter

import "github.com/YAOHAO9/pine/rpc/message"

// AfterFilterSlice map
type AfterFilterSlice []func(rpcResp *message.PineMessage) (next bool)

// Register filter
func (slice *AfterFilterSlice) Register(f func(rpcResp *message.PineMessage) (next bool)) {
	*slice = append(*slice, f)
}

// Exec filter
func (slice AfterFilterSlice) Exec(rpcResp *message.PineMessage) (next bool) {
	for _, f := range slice {
		next = f(rpcResp)
		if !next {
			return false
		}
	}
	return true
}

// After filter
var After = make(AfterFilterSlice, 0)
