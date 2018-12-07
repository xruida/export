// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package vars

// 主版本号，实际版本号可能还会加上构建日期，
// 可通过 Version() 函数获取实际的版本号。
const mainVersion = "0.1.0"

var (
	version    string
	buildDate  string
	commitHash string
)

func init() {
	if len(buildDate) == 0 {
		version = mainVersion
	} else {
		version = mainVersion + "+" + buildDate
	}
}

// Version 完整的版本号
func Version() string {
	return version
}

// CommitHash Git 上最后的提交记录 hash 值。
func CommitHash() string {
	return commitHash
}
