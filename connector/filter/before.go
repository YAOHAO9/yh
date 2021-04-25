package connector_filter

import "github.com/YAOHAO9/pine/rpc/message"

// BeforeFilterSlice map
type BeforeFilterSlice []func(rpcMsg *message.RPCMsg) error

// Register filter
func (slice *BeforeFilterSlice) Register(f func(rpcMsg *message.RPCMsg) error) {
	*slice = append(*slice, f)
}

// Exec filter(返回true标识继续往下执行)
func (slice BeforeFilterSlice) Exec(rpcMsg *message.RPCMsg) error {
	for _, f := range slice {
		err := f(rpcMsg)
		if err != nil {
			return err
		}
	}
	return nil
}

// Before filter
var Before = make(BeforeFilterSlice, 0)
