// Copyright 2020 just-codeding-0 . All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package stringpool

import (
	"testing"
)

func BenchmarkStringPool(t *testing.B) {

	t.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			str := GetStringSlice(1024)
			PutStringSlice(str)
		}
	})
}

func TestStringPool(t *testing.T) {
	str := GetStringSlice(128)
	PutStringSlice(str)

}
