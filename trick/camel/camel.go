// Copyright 2020 just-codeding-0 . All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package camel

import (
	"bytes"
	"sync"
)

var pool sync.Pool

/*
tip:
	sync.pool 复用byte数组,本身sync.pool的锁是针对于P的。
*/

func init() {
	pool = sync.Pool{}
	pool.New = func() interface{} {
		return bytes.NewBuffer([]byte{})
	}

}

func getBuffer() *bytes.Buffer {
	buf := pool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

func putBuffer(buffer *bytes.Buffer) {
	pool.Put(buffer)
}

// 确保是使用反射的name入参
func CamelToSnake(str string) string {
	if len(str) == 0 {
		return ""
	}

	buf := getBuffer()

	for i := 0; i < len(str); i++ {
		if str[i] >= 'a' {
			buf.WriteByte(str[i])
		} else {
			// 'A' + 32 = 'a' 详细看ascii表
			if i != 0 {
				buf.WriteByte('_')
			}
			buf.WriteByte(str[i] + 32)
		}
	}

	// bytes to string alloc 1 object
	bufStr := buf.String()
	putBuffer(buf)
	return bufStr
}

func SnakeToCamel(str string) string {
	if len(str) == 0 {
		return ""
	}

	buf := getBuffer()
	for i := 0; i < len(str); i++ {
		if str[i] == '_' {
			i++
			buf.WriteByte(str[i] - 32)
		} else {
			buf.WriteByte(str[i])
		}
	}

	bufStr := buf.String()
	putBuffer(buf)
	return bufStr
}
