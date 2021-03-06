// Copyright 2018 by xruida.com, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
	"os"
	"io"

	"baliance.com/gooxml/color"
	"baliance.com/gooxml/common"
	"baliance.com/gooxml/schema/soo/wml"
	"baliance.com/gooxml/measurement"
	"baliance.com/gooxml/document"
	"github.com/issue9/unique"
	"github.com/issue9/web"

	"github.com/xruida/export/common/result"
)

const dur = 5 * time.Minute

var readers = &sync.Map{}

type reader struct {
	filename string
	*bytes.Reader
	created time.Time
}

var cacheControl = fmt.Sprintf("max-age=%d, must-revalidate", int(dur.Seconds()))

func previewDoc(w http.ResponseWriter, r *http.Request) {
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

	ctx.ServeContent(buf.Reader, "text.doc", map[string]string{
		"Pragma":              "public",
		"Cache-Control":       cacheControl,
		"Content-Disposition": dis,
		"Content-type":        "application/msword",
	})
}

func exportDoc(w http.ResponseWriter, r *http.Request) {
	ctx := web.NewContext(w, r)
	type word struct {
		Key  string `orm:"name(key)" json:"key"`
		Name string `orm:"name(name)" json:"name"`
		URL      string    `orm:"name(url)" json:"url"`
		Position []float64 `orm:"name(position)" json:"position"`
		Size     []float64 `orm:"name(size)" json:"size"`
	}

	// type image struct {

	// 	Page     int       `orm:"page" json:"page"`
	// }

	data := &struct {
		URL      string `orm:"name(url)" json:"url"`
		FileName string `orm:"name(filename)" json:"filename"`
		Format   []word `orm:"name(format)" json:"format"`
	}{}

	if !ctx.Read(data) {
		return
	}

	res, err := http.Get(data.URL)
	if err != nil {
		ctx.Error(http.StatusBadRequest, err)
		return
	}
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ctx.Error(http.StatusBadRequest, err)
		return
	}
	defer res.Body.Close()

	doc, err := document.Read(bytes.NewReader(bs), int64(len(bs)))
	if err != nil {
		ctx.Error(http.StatusNotFound, err)
		return
	}

	paragraphs := []document.Paragraph{}
	for _, p := range doc.Paragraphs() {
		paragraphs = append(paragraphs, p)
	}

	// This sample document uses structured document tags, which are not common
	// except for in document templates.  Normally you can just iterate over the
	// document's paragraphs.
	for _, sdt := range doc.StructuredDocumentTags() {
		for _, p := range sdt.Paragraphs() {
			paragraphs = append(paragraphs, p)
		}
	}

	for _, p := range paragraphs {
		rm := p.Runs()

		for i := 0; i < len(rm); i++ {

			if strings.ContainsAny(rm[i].Text(), "}}") && !strings.ContainsAny(rm[i].Text(), "{{") && !strings.ContainsAny(rm[i-1].Text(), "}}") && strings.ContainsAny(rm[i-1].Text(), "{{") && i > 0 {
				x := rm[i].Text()
				y := rm[i-1].Text()
				rm[i-1].ClearContent()
				rm[i].ClearContent()
				rm[i].AddText(y + x)
			} else if strings.ContainsAny(rm[i].Text(), "}}") && !strings.ContainsAny(rm[i].Text(), "{{") && !strings.ContainsAny(rm[i-2].Text(), "}}") && strings.ContainsAny(rm[i-2].Text(), "{{") && i > 1 {
				x := rm[i].Text()
				y := rm[i-1].Text()
				z := rm[i-2].Text()
				rm[i-1].ClearContent()
				rm[i-2].ClearContent()
				rm[i].ClearContent()
				rm[i].AddText(z + y + x)
			} else if strings.ContainsAny(rm[i].Text(), "}}") && !strings.ContainsAny(rm[i].Text(), "{{") && !strings.ContainsAny(rm[i-1].Text(), "{{") && !strings.ContainsAny(rm[i-2].Text(), "{{") && !strings.ContainsAny(rm[i-3].Text(), "}}") && strings.ContainsAny(rm[i-3].Text(), "{{") && i > 3 {
				x := rm[i].Text()
				y := rm[i-1].Text()
				z := rm[i-2].Text()
				l := rm[i-3].Text()
				rm[i-1].ClearContent()
				rm[i-2].ClearContent()
				rm[i-3].ClearContent()
				rm[i].ClearContent()
				rm[i].AddText(l + z + y + x)
			}

		}
	}

	for _, p := range paragraphs {
		for _, r := range p.Runs() {
			for _, v := range data.Format {
				st := strings.SplitAfter(r.Text(), "{{"+v.Key+"}}")
				if st[0] == "{{"+v.Key+"}}" && v.URL != ""{
					r.ClearContent()
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
	
					iref, err := doc.AddImage(img)
					if err != nil {
						ctx.Error(http.StatusInternalServerError, err)
						return
					}
	
					anchored, err := p.AddRun().AddDrawingAnchored(iref)
					if err != nil {
						ctx.Error(http.StatusInternalServerError, err)
						return
					}
	
					anchored.SetSize(measurement.Distance(v.Size[0])*measurement.Point,measurement.Distance(v.Size[1])*measurement.Point)
					anchored.SetOffset(measurement.Distance(v.Position[0])* measurement.Point,measurement.Distance(v.Position[1])* measurement.Point)
				}
				if strings.Replace(r.Text(), " ", "", -1) == "{{"+v.Key+"}}" {
					st := strings.SplitAfter(r.Text(), "{{"+v.Key+"}}")
					if st[0] != "{{"+v.Key+"}}" {
						r.ClearContent()
						r.AddText(strings.Split(r.Text(), "{{"+v.Key+"}}")[0])

						r.AddText("" + v.Name + "")
						r.Properties().SetUnderline(wml.ST_UnderlineWords, color.Black)
					} else {
						r.ClearContent()
						r.AddText("" + v.Name + "")
						r.Properties().SetUnderline(wml.ST_UnderlineWords, color.Black)
					}
					if len(st) > 1 {
						// rr := p.AddRun()
					    r.AddText(st[1])
					}
				}
			}
		}
	}

	for _, p := range paragraphs {
		for _, r := range p.Runs() {
			if strings.ContainsAny(r.Text(), "{{&}}") {
				st := strings.Split(r.Text(), "{{")
				if st[0] != "{{" {
					r.ClearContent()
					r.AddText(strings.Split(r.Text(), "{{")[0] + "______")
				} else {
					r.ClearContent()
					r.AddText("____")
				}
				st = strings.Split(r.Text(), "}}")
				if len(st) > 1 {
					r.AddText(st[1])
				}
			}
		}
	}

	buf := new(bytes.Buffer)
	if err := doc.Save(buf); err != nil {
		ctx.Error(http.StatusInternalServerError, err)
		return
	}

	no := unique.Date().String()
	readers.Store(no, &reader{
		filename: data.FileName,
		Reader:   bytes.NewReader(buf.Bytes()),
		created:  time.Now(),
	})

	url, err := web.Mux().URL("/oxml/docx/{no}", map[string]string{"no": no})
	if err != nil {
		ctx.Error(http.StatusInternalServerError, err)
		return
	}
	ctx.Render(http.StatusCreated, map[string]interface{}{
		"Location": url,
	}, nil)
}
