// Copyright 2018 by xruida.com, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package xlsx

import (
	"bytes"
	"net/http"
	"strconv"

	"baliance.com/gooxml/color"
	"baliance.com/gooxml/document"
	"baliance.com/gooxml/measurement"
	"baliance.com/gooxml/schema/soo/sml"
	"baliance.com/gooxml/schema/soo/wml"
	"baliance.com/gooxml/spreadsheet"
	"github.com/issue9/web"
	// "github.com/jung-kurt/gofpdf/internal/example"
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
		Width      float64 `orm:"name(width)" json:"width"`           //宽度
		Height     float64 `orm:"name(height)" json:"height"`         //高度
		Horizontal int     `orm:"name(horizontal)" json:"horizontal"` //水平对齐 2左 3中 4右
		Color      []uint8 `orm:"name(color)" json:"color"`           //字体颜色
		Top        bool    `orm:"name(top)" json:"top"`               //顶部描边
		Bottom     bool    `orm:"name(bottom)" json:"bottom"`         //底部描边
	}

	data := &struct {
		// URL    string  `orm:"name(url);len(50)" json:"url"` //单元格内的数据
		Row    []float64 `orm:"name(row)" json:"row"`
		Format []*Wps    `orm:"name(format)" json:"format"`
	}{}

	if !ctx.Read(data) {
		return
	}

	// if len(data.URL) != 0 {
	// 	image, err := http.Get(data.URL)
	// 	r := Image{}

	// 	imgDec, ifmt, err := image.Decode(image)
	// 	if err != nil {
	// 		return r, fmt.Errorf("unable to parse image: %s", err)
	// 	}

	// 	r.Format = ifmt
	// 	r.Size = imgDec.Bounds().Size()

	// 	iref, err := ss.AddImage(r)
	// 	if err != nil {
	// 		log.Fatalf("unable to add image to workbook: %s", err)
	// 	}

	// 	dwng := ss.AddDrawing()
	// 	sheet.SetDrawing(dwng)
	// }

	ss := spreadsheet.New()
	sheet := ss.AddSheet()
	ww := new(bytes.Buffer)
	ss.Save(ww)

	// row := sheet.AddRow()
	// row.AddCell()

	for r := 0; r < 5; i++ {
		row := sheet.AddRow()
		row.SetHeight(measurement.Distance(data.Row[i] * measurement.Inch))
	}

	// for _, k := range data.Line {
	// 	centered := ss.StyleSheet.AddCellStyle()
	// 	for i := 1; i <= k.Row; i++ {
	// 		cel := transformation(i) + strconv.Itoa(k.Cell)
	// 		sheet.Cell(cel).SetStyle(centered)
	// 		bAll := ss.StyleSheet.AddBorder()
	// 		centered.SetBorder(bAll)
	// 		bAll.SetTop(sml.ST_BorderStyleThin, color.Black)
	// 	}
	// }

	for _, v := range data.Format {
		column := transformation(v.Column[0]) + strconv.Itoa(v.Column[1])
		if v.Width != 1 && v.Width != 0 {
			sheet.Column(uint32(v.Column[0])).SetWidth(measurement.Distance(v.Width * measurement.Inch))
		}

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
		//设置字体
		run.SetFont("宋体")
		if len(v.Color) == 0 {
			run.SetColor(color.Black)
		} else {
			run.SetColor(color.RGB(v.Color[0], v.Color[1], v.Color[2]))
		}
		centered := ss.StyleSheet.AddCellStyle()
		centered.SetWrapped(true)
		//合并单元格
		sheet.Cell(column).SetStyle(centered)
		centered.SetHorizontalAlignment(sml.ST_HorizontalAlignment(v.Horizontal))
		centered.SetVerticalAlignment(sml.ST_VerticalAlignmentCenter)

		// if v.Top {
		// 	bAll := ss.StyleSheet.AddBorder()
		// 	centered.SetBorder(bAll)
		// 	bAll.SetTop(sml.ST_BorderStyleThin, color.Black)
		// }

		if len(v.Enjambment) != 0 {

			enjambment := transformation(v.Enjambment[0]) + strconv.Itoa(v.Enjambment[1])

			sheet.AddMergedCells(column, enjambment)

			// sheet.Cell(column).SetStyle(centered)

			if v.Top {
				for i := v.Column[0]; i <= v.Enjambment[0]; i++ {
					cel := transformation(i) + strconv.Itoa(v.Column[1])
					sheet.Cell(cel).SetStyle(centered)
					bAll := ss.StyleSheet.AddBorder()
					centered.SetBorder(bAll)
					bAll.SetTop(sml.ST_BorderStyleThin, color.Black)
				}
			}

			if v.Bottom {
				for i := v.Column[0]; i <= v.Enjambment[0]; i++ {
					cel := transformation(i) + strconv.Itoa(v.Enjambment[1])
					sheet.Cell(cel).SetStyle(centered)
					bAll := ss.StyleSheet.AddBorder()
					centered.SetBorder(bAll)
					bAll.SetBottom(sml.ST_BorderStyleThin, color.Black)
				}
			}
			if v.Ball {
				//单元格边框设置

				// add some borders to the style (ordering isn't important, we could just as
				// easily construct the cell style and then apply it to the cell)
				for i := v.Column[1]; i <= v.Enjambment[1]; i++ {
					for j := v.Column[0]; j <= v.Enjambment[0]; j++ {
						cel := transformation(j) + strconv.Itoa(i)
						sheet.Cell(cel).SetStyle(centered)
						bAll := ss.StyleSheet.AddBorder()
						centered.SetBorder(bAll)
						bAll.SetLeft(sml.ST_BorderStyleThin, color.Black)
						bAll.SetRight(sml.ST_BorderStyleThin, color.Black)
						bAll.SetTop(sml.ST_BorderStyleThin, color.Black)
						bAll.SetBottom(sml.ST_BorderStyleThin, color.Black)
					}
				}
			}

		} else if len(v.Enjambment) == 0 {

			if v.Top {
				bAll := ss.StyleSheet.AddBorder()
				centered.SetBorder(bAll)
				bAll.SetTop(sml.ST_BorderStyleThin, color.Black)
			}

			if v.Bottom {
				bAll := ss.StyleSheet.AddBorder()
				centered.SetBorder(bAll)
				bAll.SetBottom(sml.ST_BorderStyleThin, color.Black)
			}

			if v.Ball {

				// sheet.Cell(column).SetStyle(centered)
				bAll := ss.StyleSheet.AddBorder()
				centered.SetBorder(bAll)
				bAll.SetLeft(sml.ST_BorderStyleThin, color.Black)
				bAll.SetRight(sml.ST_BorderStyleThin, color.Black)
				bAll.SetTop(sml.ST_BorderStyleThin, color.Black)
				bAll.SetBottom(sml.ST_BorderStyleThin, color.Black)
			}

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

func exportDOC(w http.ResponseWriter, r *http.Request) {
	ctx := web.NewContext(w, r)

	doc := document.New()

	// First Table
	{
		table := doc.AddTable()
		// width of the page
		table.Properties().SetWidthPercent(100)
		// with thick borers
		borders := table.Properties().Borders()
		borders.SetAll(wml.ST_BorderSingle, color.Auto, 2*measurement.Point)

		row := table.AddRow()
		run := row.AddCell().AddParagraph().AddRun()
		run.AddText("Name")
		run.Properties().SetHighlight(wml.ST_HighlightColorYellow)
		row.AddCell().AddParagraph().AddRun().AddText("John Smith")
		row = table.AddRow()
		row.AddCell().AddParagraph().AddRun().AddText("Street Address")
		row.AddCell().AddParagraph().AddRun().AddText("111 Country Road")
	}

	if err := doc.Validate(); err != nil {
		ctx.Error(http.StatusInternalServerError, err)
		return
	}

	doc.SaveToFile("tables.docx")

}
