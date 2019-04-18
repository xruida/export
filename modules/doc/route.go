// Copyright 2018 by xruida.com, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"

	"baliance.com/gooxml/color"
	"baliance.com/gooxml/document"
	"baliance.com/gooxml/schema/soo/wml"
	"github.com/issue9/web"
)

func exportDOC(w http.ResponseWriter, r *http.Request) {
	ctx := web.NewContext(w, r)
	type word struct {
		Key  string `orm:"name(key)" json:"key"`
		Name string `orm:"name(name)" json:"name"`
	}

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
				if strings.Replace(r.Text(), " ", "", -1) == "{{"+v.Key+"}}" {
					st := strings.SplitAfter(r.Text(), "{{"+v.Key+"}}")
					if st[0] != "{{"+v.Key+"}}" {
						r.ClearContent()
						r.AddText(strings.Split(r.Text(), "{{"+v.Key+"}}")[0])

						r.AddText("__" + v.Name + "__")
						r.Properties().SetUnderline(wml.ST_UnderlineWords, color.Black)
					} else {
						r.ClearContent()
						r.AddText("__" + v.Name + "__")
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

	ww := new(bytes.Buffer)
	doc.Save(ww)

	buf := new(bytes.Buffer)
	if err := doc.Save(buf); err != nil {
		ctx.Error(http.StatusInternalServerError, err)
		return
	}

	reader := bytes.NewReader(buf.Bytes())

	ctx.ServeContent(reader, "text.doc", map[string]string{
		"Pragma":              "public",
		"Cache-Control":       "must-revalidate",
		"Content-Disposition": "attachment; filename=file1.doc",
		"Content-type":        "application/msword",
	})
}
