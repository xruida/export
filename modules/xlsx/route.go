// Copyright 2018 by xruida.com, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package xlsx

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"baliance.com/gooxml/color"
	"baliance.com/gooxml/common"
	"baliance.com/gooxml/measurement"
	"baliance.com/gooxml/schema/soo/sml"
	"baliance.com/gooxml/spreadsheet"
	"github.com/issue9/unique"
	"github.com/issue9/web"
	"github.com/xruida/export/common/result"
	// "github.com/jung-kurt/gofpdf/internal/example"
)

const dur = 5 * time.Minute

var readers = &sync.Map{}

type reader struct {
	filename string
	*bytes.Reader
	created time.Time
}

var cacheControl = fmt.Sprintf("max-age=%d, must-revalidate", int(dur.Seconds()))

func previewXlsx(w http.ResponseWriter, r *http.Request) {
	ctx := web.NewContext(w, r)

	no, err := ctx.ParamString("no")
	if err != nil {
		ctx.NewResult(result.BadRequestInvalidParam).Add("no", err.Error()).Render()
		return
	}

	rr, found := readers.Load(no)
	if !found {
		ctx.Exit(http.StatusGone)
		return
	}

	buf := rr.(*reader)

	name := url.QueryEscape(buf.filename)

	dis := "attachment; filename=\"" + name + "\""

	ctx.ServeContent(buf.Reader, "text.xlsx", map[string]string{
		"Pragma":              "public",
		"Cache-Control":       cacheControl,
		"Content-Disposition": dis,
		"Content-type":        "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	})
}

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

	type image struct {
		URL      string    `orm:"name(url)" json:"url"`
		Position []float64 `orm:"name(position)" json:"position"`
	}

	data := &struct {
		FileName string    `orm:"name(filename)" json:"filename"`
		Image    []*image  `orm:"name(image)" json:"image"`
		Row      []float64 `orm:"name(row)" json:"row"`
		Format   []*Wps    `orm:"name(format)" json:"format"`
	}{}

	if !ctx.Read(data) {
		return
	}

	ss := spreadsheet.New()
	sheet := ss.AddSheet()

	dwng := ss.AddDrawing()
	sheet.SetDrawing(dwng)

	for _, v := range data.Image {
		res, err := http.Get(v.URL)
		if err != nil {
			ctx.Error(http.StatusBadRequest, err)
			return
		}

		err = os.MkdirAll("./upload", os.ModePerm)
		if os.IsNotExist(err) {
			ctx.Error(http.StatusInternalServerError, err)
			return
		}

		dir := "./upload/" + unique.Number().String() + ".jpg"

		t, err := os.OpenFile(dir, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			ctx.Error(http.StatusInternalServerError, err)
			return
		}

		defer t.Close()

		io.Copy(t, res.Body)
		go func() {
			select {
			case <-time.After(5 * time.Minute):
				err := os.Remove(dir)
				if err != nil {
					ctx.Error(http.StatusInternalServerError, err)
					return
				}
			}

			return
		}()

		img, err := common.ImageFromFile(dir)
		if err != nil {
			ctx.Error(http.StatusInternalServerError, err)
			return
		}

		iref, err := ss.AddImage(img)
		if err != nil {
			ctx.Error(http.StatusInternalServerError, err)
			return
		}

		anc := dwng.AddImage(iref, spreadsheet.AnchorTypeAbsolute)

		anc.SetColOffset(measurement.Distance(v.Position[0]) * measurement.Point)
		anc.SetRowOffset(measurement.Distance(v.Position[1]) * measurement.Point)

		anc.SetWidth(measurement.Distance(v.Position[2]) * measurement.Point)
		anc.SetHeight(iref.RelativeHeight(measurement.Distance(v.Position[3]) * measurement.Point))
	}

	for i := 0; i < len(data.Row); i++ {
		row := sheet.AddRow()
		row.SetHeight(measurement.Distance(data.Row[i] * measurement.Inch))
	}

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

		if len(v.Enjambment) != 0 {

			enjambment := transformation(v.Enjambment[0]) + strconv.Itoa(v.Enjambment[1])

			sheet.AddMergedCells(column, enjambment)

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
				bAll := ss.StyleSheet.AddBorder()
				centered.SetBorder(bAll)
				bAll.SetLeft(sml.ST_BorderStyleThin, color.Black)
				bAll.SetRight(sml.ST_BorderStyleThin, color.Black)
				bAll.SetTop(sml.ST_BorderStyleThin, color.Black)
				bAll.SetBottom(sml.ST_BorderStyleThin, color.Black)
			}

		}
	}

	ww := new(bytes.Buffer)
	ss.Save(ww)

	if err := ss.Validate(); err != nil {
		ctx.Error(http.StatusInternalServerError, err)
		return
	}

	buf := new(bytes.Buffer)
	if err := ss.Save(buf); err != nil {
		ctx.Error(http.StatusInternalServerError, err)
		return
	}

	no := unique.Date().String()
	readers.Store(no, &reader{
		filename: data.FileName,
		Reader:   bytes.NewReader(buf.Bytes()),
		created:  time.Now(),
	})

	url, err := web.Mux().URL("/oxml/xlsx/{no}", map[string]string{"no": no})
	if err != nil {
		ctx.Error(http.StatusInternalServerError, err)
		return
	}

	ctx.Render(http.StatusCreated, map[string]interface{}{
		"Location": url,
	}, nil)
}
