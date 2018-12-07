// Copyright 2018 by xruida.com, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package result 定义了所有的错误消息代码及对应的消息
package result

import (
	"net/http"

	"github.com/issue9/web/result"
)

// Init 初始化消息内容
func Init() error {
	return result.NewMessages(messages)
}

// 400
const (
	BadRequest = http.StatusBadRequest*100 + iota
	BadRequestInvalidParam
	BadRequestInvalidQuery
	BadRequestInvalidBody
	BadRequestInvalidHeader
)

// 401
const (
	Unauthorized = http.StatusUnauthorized*100 + iota
	UnauthorizedInvalidState
	UnauthorizedInvalidToken
	UnauthorizedInvalidPassword
	UnauthorizedInvalidUsername
	UnauthorizedNeedChangePassword
	UnauthorizedAuthTokenExpired
)

// 403
const (
	Forbidden = http.StatusForbidden*100 + iota
	ForbiddenNewOldPasswordIsEqual
	ForbiddenExistSubItem
	ForbiddenStateNotAllow
	ForbiddenIsFull
	ForbiddenNotDeleteOwn
	ForbiddenOnlyCompanyOwner
	ForbiddenExists
)

// 404
const (
	NotFoundCompany = http.StatusNotFound*100 + iota
	NotFoundDepartment
	NotFoundStaff
)

var messages = map[int]string{
	BadRequest:              "错误的消息",
	BadRequestInvalidParam:  "地址参数错误",
	BadRequestInvalidQuery:  "查询参数错误",
	BadRequestInvalidBody:   "请求内容错误",
	BadRequestInvalidHeader: "请求报头错误",

	Unauthorized:                   "账号或是密码错误",
	UnauthorizedInvalidState:       "当前状态无法登录",
	UnauthorizedInvalidToken:       "token 无效",
	UnauthorizedInvalidPassword:    "密码无效",
	UnauthorizedInvalidUsername:    "无效的账号",
	UnauthorizedNeedChangePassword: "只有修改才能进行其它操作",
	UnauthorizedAuthTokenExpired:   "验证码已经过期", // 第三方的 access token 过期

	Forbidden:                      "权限错误",
	ForbiddenNewOldPasswordIsEqual: "新旧密码不能相同",
	ForbiddenExistSubItem:          "存在子项目",
	ForbiddenStateNotAllow:         "当前状态不允许该操作",
	ForbiddenIsFull:                "子项目已满，不允许再添加",
	ForbiddenNotDeleteOwn:          "不能删除自己",
	ForbiddenOnlyCompanyOwner:      "该操作只能是企业的所有者才能执行",
	ForbiddenExists:                "已经存在",

	NotFoundCompany:    "公司不存在",
	NotFoundDepartment: "部门不存在",
	NotFoundStaff:      "用户不存在",
}
