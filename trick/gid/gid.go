// Copyright 2020 just-codeding-0 . All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package gid

import "syscall"

func GetGid() int {
	return syscall.Getgid()
}