package filter

var beforeFilterMap = make(BeforeFilterSlice, 0)
var afterFilterMap = make(AfterFilterSlice, 0)

// BeforeFilterManager return beforeFilterMap
func BeforeFilterManager() *BeforeFilterSlice {
	return &beforeFilterMap
}

// AfterFilterManager return afterFilterMap
func AfterFilterManager() *AfterFilterSlice {
	return &afterFilterMap
}
