// Copyright 2018 by xruida.com, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"context"
	"time"

	"github.com/issue9/web"
)

// ModuleName 模块名称
const ModuleName = "doc"

// Init 初始化信息
func Init() {
	m := web.NewModule(ModuleName, "导出 oxml 的 doc 服务")

	m.AddService(clearBuf, "清除缓存的内容")

	m.PostFunc("/oxml/docx", exportDoc).
		GetFunc("/oxml/docx/preview/{no}", previewDoc)
}

func clearBuf(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Minute) // 一分钟清除一次

	for {
		select {
		case <-ctx.Done():
			return context.Canceled
		case now := <-ticker.C:
			readers.Range(func(k, v interface{}) bool {
				r := v.(*reader)
				if r.created.Add(dur).After(now) {
					readers.Delete(k)
				}

				return true
			})
		}
	}
}
