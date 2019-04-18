// Copyright 2018 by xruida.com, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/issue9/logs"
	"github.com/issue9/web"

	"github.com/xruida/export/common/result"
	"github.com/xruida/export/common/vars"
	"github.com/xruida/export/modules/doc"
	"github.com/xruida/export/modules/xlsx"
)

func main() {
	h := flag.Bool("h", false, "显示帮助信息")
	v := flag.Bool("v", false, "显示版本号")
	c := flag.String("c", vars.AppConfig, "指定路径")
	flag.Parse()

	switch {
	case *h:
		flag.PrintDefaults()
		return
	case *v:
		fmt.Printf("%s:%s build with %s\n", vars.Name, vars.Version(), runtime.Version())
		fmt.Println("common hash:", vars.CommitHash())
		return
	}

	if err := web.Classic(*c); err != nil {
		panic(err)
	}

	result.Init()

	initModules()
	web.InitModules("")
	web.Fatal(2, web.Serve())
}

func initModules() {
	logs.Trace("开始初始化模块")
	xlsx.Init()
	doc.Init()
	logs.Trace("初始化模块完成")
}
