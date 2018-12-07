// Copyright 2018 by xruida.com, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main // import "github.com/xruida/export"

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"runtime"

	"github.com/issue9/logs"
	"github.com/issue9/web"
	"github.com/issue9/web/encoding"
	"github.com/issue9/web/encoding/gob"
	yaml "gopkg.in/yaml.v2"

	"code.xruida.com/xruida/server/common/result"
	"code.xruida.com/xruida/server/common/vars"
	"code.xruida.com/xruida/server/modules/xlsx"
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

	err := encoding.AddUnmarshals(map[string]encoding.UnmarshalFunc{
		encoding.DefaultMimeType: gob.Unmarshal,
		"application/json":       json.Unmarshal,
		"application/xml":        xml.Unmarshal,
		"text/vnd.yaml":          yaml.Unmarshal,
	})
	if err != nil {
		panic(err)
	}

	err = encoding.AddMarshals(map[string]encoding.MarshalFunc{
		encoding.DefaultMimeType: gob.Marshal,
		"application/json":       json.Marshal,
		"application/xml":        xml.Marshal,
		"text/vnd.yaml":          yaml.Marshal,
	})
	if err != nil {
		panic(err)
	}

	if err := result.Init(); err != nil {
		panic(err)
	}

	if err := web.Init(*c); err != nil {
		panic(err)
	}

	initModules()

	logs.Critical(web.Serve())
	logs.Flush()
}

func initModules() {
	logs.Trace("开始初始化模块")
	xlsx.Init()
	logs.Trace("初始化模块完成")
}
