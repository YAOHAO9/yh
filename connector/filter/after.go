package connector_filter

import "github.com/YAOHAO9/pine/rpc/message"

// AfterFilterSlice map
type AfterFilterSlice []func(rpcResp *message.PineMsg) (error)

// Register filter
func (slice *AfterFilterSlice) Register(f func(rpcResp *message.PineMsg) (error)) {
	*slice = append(*slice, f)
}

// Exec filter
func (slice AfterFilterSlice) Exec(rpcResp *message.PineMsg) (error) {
	for _, f := range slice {
		err := f(rpcResp)
		return err
	}
	return nil
}

// After filter
var After = make(AfterFilterSlice, 0)
