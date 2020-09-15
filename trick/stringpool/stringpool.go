// Copyright 2020 just-codeding-0 . All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package stringpool

import "sync"

var arr []strPool
var lastMax int
var (
	// 二的幂
	minB = 3
	maxB = 10
)

type strPool struct {
	min  int
	max  int
	pool sync.Pool
}

func init() {

	var baseData = make([]int, 0, maxB-minB)

	for i := minB; i <= maxB; i++ { // 8 16 32 64 128 ....... 1024
		baseData = append(baseData, 1<<i)
	}

	for _, v := range baseData {
		p := strPool{
			min:  lastMax + 1,
			max:  v,
			pool: sync.Pool{},
		}

		lastMax = v
		tmp := v
		p.pool.New = func() interface{} {
			return make([]string, 0, tmp)
		}
		arr = append(arr, p)
	}

}

func GetStringSlice(length int) []string {
	if length == 0 {
		return nil
	}

	for i := 0; i < len(arr); i++ {
		if arr[i].max >= length {
			return arr[i].pool.Get().([]string)
		}
	}

	return make([]string, 0, length)
}

func PutStringSlice(str []string) {
	l := cap(str)
	if l <= lastMax { // 得在范围内,否则直接丢弃
		idx := minB
		for idx <= maxB {
			if 1<<idx == l { // 检查是否为2的幂
				str = str[:0] // reset length
				arr[idx-minB].pool.Put(str)
			}
			idx++
		}

	}

}
