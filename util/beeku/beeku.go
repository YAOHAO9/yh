package beeku

import (
	"math/rand"
	"time"
)

type reducetype func(interface{}) interface{}
type filtertype func(interface{}) bool

// SliceRandList 随机切片
func SliceRandList(min, max int) []int {
	if max < min {
		min, max = max, min
	}
	length := max - min + 1
	t0 := time.Now()
	rand.Seed(int64(t0.Nanosecond()))
	list := rand.Perm(length)
	for index := range list {
		list[index] += min
	}
	return list
}

// SliceMerge 合并切片
func SliceMerge(slice1, slice2 []interface{}) (c []interface{}) {
	c = append(slice1, slice2...)
	return
}

// InSlice 判断切片是否包含某个值
func InSlice(val interface{}, slice ...interface{}) int {
	for index, v := range slice {
		if v == val {
			return index
		}
	}
	return -1
}

// SliceReduce 切片Reduce
func SliceReduce(slice []interface{}, a reducetype) (dslice []interface{}) {
	for _, v := range slice {
		dslice = append(dslice, a(v))
	}
	return
}

// SliceRand 随机返回切片中的一个值
func SliceRand(a []interface{}) (b interface{}) {
	randnum := rand.Intn(len(a))
	b = a[randnum]
	return
}

// SliceSum 切片求和
func SliceSum(intslice []int64) (sum int64) {
	for _, v := range intslice {
		sum += v
	}
	return
}

// SliceFilter 过滤掉切片中不同类型的item
func SliceFilter(slice []interface{}, validType filtertype) (ftslice []interface{}) {
	for _, v := range slice {
		if validType(v) {
			ftslice = append(ftslice, v)
		}
	}
	return
}

// SliceDiff 返回slice1中有slice2中没有的字段
func SliceDiff(slice1, slice2 []interface{}) (diffslice []interface{}) {
	for _, v := range slice1 {
		if InSlice(v, slice2) == -1 {
			diffslice = append(diffslice, v)
		}
	}
	return
}

// SliceIntersect 返回两个切片中共有的item
func SliceIntersect(slice1, slice2 []interface{}) (diffslice []interface{}) {
	for _, v := range slice1 {
		if InSlice(v, slice2) != -1 {
			diffslice = append(diffslice, v)
		}
	}
	return
}

// SliceRange 生成等差数列切片
func SliceRange(start, end, step int64) (intslice []int64) {
	for i := start; i <= end; i += step {
		intslice = append(intslice, i)
	}
	return
}

// SlicePadStart 填充切片
func SlicePadStart(slice []interface{}, size int, val interface{}) []interface{} {
	if size <= len(slice) {
		return slice
	}
	leftLength := (size - len(slice))
	for i := 0; i < leftLength; i++ {
		newSlice := []interface{}{val}
		slice = append(newSlice, slice...)
	}
	return slice
}

// SlicePadEnd 填充切片
func SlicePadEnd(slice []interface{}, size int, val interface{}) []interface{} {
	if size <= len(slice) {
		return slice
	}
	leftLength := (size - len(slice))
	for i := 0; i < leftLength; i++ {
		slice = append(slice, val)
	}
	return slice
}

// SliceUnique 去除重复值
func SliceUnique(slice []interface{}) (uniqueslice []interface{}) {
	for _, v := range slice {
		if InSlice(v, uniqueslice) == -1 {
			uniqueslice = append(uniqueslice, v)
		}
	}
	return
}

// SliceShuffle 切片乱序
func SliceShuffle(slice []interface{}) []interface{} {
	for i := len(slice) - 1; i > 0; i-- {
		a := rand.Intn(i)
		slice[a], slice[i] = slice[i], slice[a]
	}
	return slice
}
