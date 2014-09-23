package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"reflect"
)

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
	FlushMessage string

	eventFunc map[string][]reflect.Value
	IsSend    bool
	IsEnd     bool

	T *template.Template

	IsUpload   bool
	IsDownload bool
}

func NewContext(res http.ResponseWriter, req *http.Request, tpls *template.Template) *Context {
	context := new(Context)
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
	context.FlushMessage = ""

	context.T = tpls

	context.IsUpload = false
	context.IsDownload = false

	req.ParseForm()
	return context
}

func (ctx *Context) Param(key string) string {
	return ctx.routerParams[key]
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
	err := ctx.T.ExecuteTemplate(ctx.Response, tplName, ctx)

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
