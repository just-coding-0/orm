// Copyright 2020 just-codeding-0 . All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package camel

import (
	"math/rand"
	"testing"
	"time"
)

/* stats
command: go test -bench=BenchmarkCamelToSnake -benchmem -cpu=1,2,4,8
goos: darwin
goarch: amd64
pkg: github.com/just-coding-0/simple-orm-framework/trick/camel
BenchmarkCamelToSnake            6788091               173 ns/op              32 B/op          1 allocs/op
BenchmarkCamelToSnake-2          7172932               160 ns/op              32 B/op          1 allocs/op
BenchmarkCamelToSnake-4          6650732               179 ns/op              32 B/op          1 allocs/op
BenchmarkCamelToSnake-8          5754320               205 ns/op              32 B/op          1 allocs/op




command: go test -bench=BenchmarkSnakeToCamel -benchmem -cpu=1,2,4,8

goos: darwin
goarch: amd64
pkg: github.com/just-coding-0/simple-orm-framework/trick/camel
BenchmarkSnakeToCamel            9272526               132 ns/op              15 B/op          1 allocs/op
BenchmarkSnakeToCamel-2          9257449               128 ns/op              15 B/op          1 allocs/op
BenchmarkSnakeToCamel-4          7981471               152 ns/op              15 B/op          1 allocs/op
BenchmarkSnakeToCamel-8          6354708               186 ns/op              15 B/op          1 allocs/op

*/

func BenchmarkCamelToSnake(b *testing.B) {
	testCase := getCamelRandStr(1000)

	rand.Seed(time.Now().Unix())
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			CamelToSnake(testCase[rand.Int()%1000])
		}
	})
}

func BenchmarkSnakeToCamel(b *testing.B) {
	testCase := getCamelSnakeStr(1000)

	rand.Seed(time.Now().Unix())
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			SnakeToCamel(testCase[rand.Int()%1000])
		}
	})
}

var camelBaseByte = []byte{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
}

func getCamelRandStr(count int) []string {
	var arr []string
	var buf = make([]byte, 0, 10)
	for i := 0; i < count; i++ {
		for len(buf) != cap(buf) {
			buf = append(buf, camelBaseByte[rand.Int()%len(camelBaseByte)])
		}

		// 这里不能指向同一个data域,所以就不能使用零拷贝了。
		arr = append(arr, string(buf))
		buf = buf[:0]
	}

	return arr
}

var snakeBaseByte = []byte{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'_',
}

func getCamelSnakeStr(count int) []string {
	var arr []string
	var buf = make([]byte, 0, 10)
	for i := 0; i < count; i++ {
		for len(buf) != cap(buf) {
			_byte := snakeBaseByte[rand.Int()%len(snakeBaseByte)]
			for len(buf)+1 == cap(buf) && _byte == '_' {
				_byte = snakeBaseByte[rand.Int()%len(snakeBaseByte)]
			}
			buf = append(buf,_byte)
		}

		// 这里不能指向同一个data域,所以就不能使用零拷贝了。
		arr = append(arr, string(buf))
		buf = buf[:0]
	}

	return arr
}
