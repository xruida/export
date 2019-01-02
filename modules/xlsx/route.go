// Copyright 2018 by xruida.com, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package xlsx

import (
	"bytes"
	"net/http"
	"strconv"

	"baliance.com/gooxml/color"
	"baliance.com/gooxml/measurement"
	"baliance.com/gooxml/schema/soo/sml"
	"baliance.com/gooxml/spreadsheet"
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

	type Wps struct {
		Name       string  `orm:"name(name);len(50)" json:"name"`     //单元格内的数据
		Column     []int   `orm:"name(column)" json:"column"`         //所在单元格
		Enjambment []int   `orm:"name(enjambment)" json:"enjambment"` //合并单元格的右下坐标
		Ball       bool    `orm:"name(ball)" json:"ball"`             //单元格是否有边框
		Italic     bool    `orm:"name(italic)" json:"italic"`         //斜体
		Bold       bool    `orm:"name(bold)" json:"bold"`             //黑体
		Size       float64 `orm:"name(size)" json:"size"`             //字体大小
	}

	data := &struct {
		Format []*Wps `orm:"name(format)" json:"format"`
	}{}

	if !ctx.Read(data) {
		return
	}

	ss := spreadsheet.New()
	sheet := ss.AddSheet()
	ww := new(bytes.Buffer)
	ss.Save(ww)
	row := sheet.AddRow()
	row.AddCell()

	for _, v := range data.Format {

		column := transformation(v.Column[0]) + strconv.Itoa(v.Column[1])

		rt := sheet.Cell(column).SetRichTextString()
		run := rt.AddRun()
		//设置单元格内的数据
		run.SetText(v.Name)
		//设置字体大小
		run.SetSize(measurement.Distance(v.Size))
		//设置加粗
		run.SetBold(v.Bold)
		//设置斜体
		run.SetItalic(v.Italic)
		run.SetColor(color.Black)
		centered := ss.StyleSheet.AddCellStyle()
		//合并单元格
		if len(v.Enjambment) != 0 {

			enjambment := numTostring(v.Enjambment[0]) + strconv.Itoa(v.Enjambment[1])

			sheet.AddMergedCells(column, enjambment)

			centered.SetHorizontalAlignment(sml.ST_HorizontalAlignmentCenter)
			centered.SetVerticalAlignment(sml.ST_VerticalAlignmentCenter)
			sheet.Cell(column).SetStyle(centered)

		}
		if v.Ball {
			//单元格边框设置
			sheet.Cell(column).SetStyle(centered)

			// add some borders to the style (ordering isn't important, we could just as
			// easily construct the cell style and then apply it to the cell)
			bAll := ss.StyleSheet.AddBorder()
			centered.SetBorder(bAll)
			bAll.SetLeft(sml.ST_BorderStyleThin, color.Black)
			bAll.SetRight(sml.ST_BorderStyleThin, color.Black)
			bAll.SetTop(sml.ST_BorderStyleThin, color.Black)
			bAll.SetBottom(sml.ST_BorderStyleThin, color.Black)
		}

	}

	if err := ss.Validate(); err != nil {
		ctx.Error(http.StatusInternalServerError, err)
		return
	}

	buf := new(bytes.Buffer)
	if err := ss.Save(buf); err != nil {
		ctx.Error(http.StatusInternalServerError, err)
		return
	}

	reader := bytes.NewReader(buf.Bytes())

	ctx.ServeContent(reader, "text.xlsx", map[string]string{
		"Pragma":              "public",
		"Cache-Control":       "must-revalidate",
		"Content-Disposition": "attachment; filename=file1.xlsx",
		"Content-type":        "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	})
}
