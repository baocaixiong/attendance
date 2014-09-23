package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"path/filepath"
	"strings"

	"net/http"
	"os"
	"path"
	"reflect"
)

const (
	CONTEXT_RENDERED = "context_rendered"
	CONTEXT_END      = "context_end"
	CONTEXT_SEND     = "context_send"
)

var currentDir = func() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}()

func Static(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(" [url]", filepath.Join(currentDir, r.URL.String()))
	http.ServeFile(w, r, filepath.Join(currentDir, r.URL.String()))
}

type Tpls struct {
	tpls       map[string]*template.Template
	tplDir     string
	mainLayout string
}

func newTpls(dir ...string) *Tpls {
	tpls := new(Tpls)

	if len(dir) == 1 {
		tpls.tplDir = dir[0]
	} else {
		tpls.tplDir = filepath.Join(currentDir, "tpls")
	}

	tpls.tpls = make(map[string]*template.Template)
	tpls.mainLayout = filepath.Join(tpls.tplDir, "main.layout")

	tpls.loadTpls()

	return tpls
}

func (h *Tpls) loadTpls() {
	filepath.Walk(h.tplDir, func(filePath string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(filePath), ".tpl") {
			t, err := template.ParseFiles(h.mainLayout, filePath)
			if err != nil {
				fmt.Println("[error]", err)
				return err
			}
			h.tpls[strings.Split(info.Name(), ".")[0]] = t

		}

		return nil
	})
}

func (t *Tpls) getTpl(name string) (*template.Template, error) {
	if tpl, ok := t.tpls[name]; ok {
		return tpl, nil
	} else {
		return nil, errors.New(fmt.Sprintf("%s template is not found!", name))
	}
}

type Context struct {
	Request    *http.Request
	Url        string
	RequsetUrl string
	Method     string
	Ext        string

	Response http.ResponseWriter
	Status   int
	Header   map[string]string
	Body     []byte

	routerParams map[string]string
	flashData    map[string]interface{}

	eventFunc map[string][]reflect.Value
	IsSend    bool
	IsEnd     bool

	tpls *Tpls

	IsUpdate   bool
	IsDownload bool
}

func NewContext(res http.ResponseWriter, req *http.Request, tpls *Tpls) *Context {
	context := new(Context)
	context.flashData = make(map[string]interface{})
	context.IsSend = false
	context.IsEnd = false

	context.Request = req
	context.Url = req.URL.Path
	context.RequsetUrl = req.RequestURI
	context.Method = req.Method
	context.Ext = path.Ext(req.URL.Path)

	context.Response = res
	context.Status = 200
	context.Header = make(map[string]string)
	context.Header["Content-Type"] = "text/html;charset=UTF-8"

	context.tpls = tpls

	context.IsUpdate = false
	context.IsDownload = false

	req.ParseForm()
	return context
}

func (ctx *Context) Param(key string) string {
	return ctx.routerParams[key]
}

func (ctx *Context) Flash(key string, v ...interface{}) interface{} {
	if len(v) == 0 {
		return ctx.flashData[key]
	}
	ctx.flashData[key] = v[0]
	return nil
}

func (ctx *Context) Input() map[string]string {
	data := make(map[string]string)
	for key, v := range ctx.Request.Form {
		data[key] = v[0]
	}
	return data
}

func (ctx *Context) Redirect(url string, status ...int) {
	code := 302
	if len(status) > 0 {
		code = status[0]
	}

	http.Redirect(ctx.Response, ctx.Request, url, code)
}

func (ctx *Context) Send() {
	if ctx.IsSend {
		return
	}
	for name, value := range ctx.Header {
		ctx.Response.Header().Set(name, value)
	}
	ctx.Response.WriteHeader(ctx.Status)
	ctx.Response.Write(ctx.Body)
	ctx.IsSend = true
}

func (ctx *Context) End() {
	if ctx.IsEnd {
		return
	}

	ctx.IsEnd = true
}

func (ctx *Context) Render(tplName string) {
	tpl, err := ctx.tpls.getTpl(tplName)
	if err != nil {
		fmt.Println("[error]", err)
		fmt.Fprintf(ctx.Response, err.Error())
		return
	}

	err = tpl.ExecuteTemplate(ctx.Response, "main", ctx)
	if err != nil {
		fmt.Println("[error]", err)
		fmt.Fprintf(ctx.Response, err.Error())
		return
	}

	ctx.End()
}

func (ctx *Context) Download(file string) {
	f, e := os.Stat(file)
	if e != nil {
		ctx.Status = 404
		return
	}
	if f.IsDir() {
		ctx.Status = 403
		return
	}

	output := ctx.Response.Header()
	output.Set("Content-Type", "application/octet-stream")
	output.Set("Content-Disposition", "attachment; filename="+path.Base(file))
	output.Set("Content-Transfer-Encoding", "binary")
	output.Set("Expires", "0")
	output.Set("Cacha-Control", "must-revalidate")
	output.Set("Parama", "public")
	http.ServeFile(ctx.Response, ctx.Request, file)
	ctx.IsSend = true
}
