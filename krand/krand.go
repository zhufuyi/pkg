package krand

import (
	"math/rand"
	"time"
)

const (
	R_NUM   = 1 // 纯数字
	R_UPPER = 2 // 大写字母
	R_LOWER = 4 // 小写字母
	R_All   = 7 // 数字、大小写字母
)

var (
	r        = rand.New(rand.NewSource(time.Now().UnixNano()))
	RefSlice = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	kinds    = [][]byte{RefSlice[0:10], RefSlice[10:36], RefSlice[0:36], RefSlice[36:62], RefSlice[36:], RefSlice[10:62], RefSlice[0:62]}
)

// KRand 生成多种类型的任意长度的随机字符串，如果参数size为空，默认长度为6
// example：KRand(R_ALL), KRand(R_ALL, 16), KRand(R_NUM|R_LOWER, 16)
func KRand(kind int, size ...int) []byte {
	if kind > 7 || kind < 1 {
		kind = R_All
	}

	length := 0
	if len(size) == 0 {
		length = 6 // 默认长度
	} else {
		length = size[0] // 只有第0个值有效，忽略其它值
		if length < 1 {
			length = 6 // 默认长度
		}
	}

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = kinds[kind-1][r.Intn(len(kinds[kind-1]))]
	}

	return result
}

// RandInt 生成指定范围大小随机数，兼容RandInt()，RandInt(max)，RandInt(min, max)，RandInt(max, min)4种方式，注：随机数包括min和max
func RandInt(rangeSize ...int) int {
	switch len(rangeSize) {
	case 0:
		return r.Intn(101) // 默认0~100
	case 1:
		return r.Intn(rangeSize[0] + 1)
	default:
		if rangeSize[0] > rangeSize[1] {
			rangeSize[0], rangeSize[1] = rangeSize[1], rangeSize[0]
		}
		return r.Intn(rangeSize[1]-rangeSize[0]+1) + rangeSize[0]
	}
}
