package filter

var beforeFilterMap = make(BeforeFilterSlice, 0)
var afterFilterMapOfRequest = make(AfterFilterSlice, 0)

// BeforeFilterManager return beforeFilterMap
func BeforeFilterManager() BeforeFilterSlice {
	return beforeFilterMap
}

// AfterFilterManagerOfRequest return afterFilterMapOfRequest
func AfterFilterManagerOfRequest() AfterFilterSlice {
	return afterFilterMapOfRequest
}
