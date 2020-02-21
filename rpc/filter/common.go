package filter

import (
	"trial/rpc/msg"

	"github.com/gorilla/websocket"
)

//===================================================
// Before
//===================================================

// BeforeFilterSlice map
type BeforeFilterSlice []func(respConn *websocket.Conn, fm *msg.ForwardMessage) (next bool)

// Register filter
func (slice *BeforeFilterSlice) Register(f func(respConn *websocket.Conn, fm *msg.ForwardMessage) (next bool)) {
	*slice = append(*slice, f)
}

// Exec filter
func (slice BeforeFilterSlice) Exec(respConn *websocket.Conn, fm *msg.ForwardMessage) (next bool) {
	for _, f := range slice {
		next = f(respConn, fm)
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
type AfterFilterSlice []func(rm *msg.ResponseMessage) (next bool)

// Register filter
func (slice *AfterFilterSlice) Register(f func(rm *msg.ResponseMessage) (next bool)) {
	*slice = append(*slice, f)
}

// Exec filter
func (slice AfterFilterSlice) Exec(rm *msg.ResponseMessage) (next bool) {
	for _, f := range slice {
		next = f(rm)
		if !next {
			return false
		}
	}
	return true
}
