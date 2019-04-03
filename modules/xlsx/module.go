// Copyright 2018 by xruida.com, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package xlsx

import (
	"github.com/issue9/web"
)

// ModuleName 模块名称
const ModuleName = "xlsx"

// Init 初始化信息
func Init() {
	m := web.NewModule(ModuleName, "导出 oxml 的 xlsx 服务")

	m.PostFunc("/oxml/xlsx", exportXLSX)
	m.PostFunc("/oxml/doc", exportDOC)
	m.PostFunc("/oxml/upload", uploadDOC)
}
