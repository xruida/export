// Copyright 2018 by xruida.com, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"io"

	"github.com/issue9/web"
)

// ModuleName 模块名称
const ModuleName = "doc"

var readers = map[string]io.ReadSeeker{}

// Init 初始化信息
func Init() {
	m := web.NewModule(ModuleName, "导出 oxml 的 doc 服务")

	m.PostFunc("/oxml/docx", exportDoc).
		GetFunc("/oxml/docx/preview/{id}", previewDoc)
}
