// Copyright 2018 by xruida.com, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package xlsx

import (
	"net/http"

	"github.com/issue9/web"
)

// @api GET /oxml/xlsx 导出 xlsx 内容
// @apiGroup admin
//
// @apiRequest json
// @apiHeader Authorization 提交登录凭证 accessToken
//
// @apiSuccess 200 OK
// @apiExample json
//  {
//      "admin:resources-list":"查看资源列表",
//      "admin:resources-list":"查看资源列表"
//  }
func exportXLSX(w http.ResponseWriter, r *http.Request) {
	ctx := web.NewContext(w, r)
	// TODO
	ctx.Render(http.StatusNotImplemented, nil, nil)
}
