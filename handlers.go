package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

type Handler struct {
	tpls *template.Template
	log  *log.Logger
}

func newHandler(tpls *template.Template) *Handler {
	handler := new(Handler)
	handler.tpls = tpls
	handler.log = logger

	return handler
}

func (h *Handler) handle(w http.ResponseWriter, r *http.Request) {
	context := NewContext(w, r, h.tpls)
	h.log.Println(fmt.Sprintf("%s %s", r.Method, r.URL.String()))
	switch r.URL.String() {
	case "/":
		context.IsUpload = true
		h.Home(context)
	case "/upload/do":
		context.IsUpload = true
		h.UploadDo(context)
	case "/download":
		context.IsDownload = true
		h.Download(context)
	case "/download/do":
		context.IsDownload = true
		h.DownloadDo(context)
	case "/favicon.ico":
	default:
		h.NotFound(context)
	}
}

func (h *Handler) Home(ctx *Context) {
	ctx.Render("home.html")
}

func (h *Handler) NotFound(ctx *Context) {
	ctx.Response.Write([]byte("This page is not found."))
}

func (h *Handler) UploadDo(ctx *Context) {
	if ctx.Method != "POST" {
		ctx.Redirect("/")
		return
	}

	ctx.Request.ParseMultipartForm(32 << 20)

	tmpDate := ctx.Request.Form.Get("date")

	date, err := time.Parse("2006-01-02", tmpDate)
	if err != nil {
		ctx.FlushMessage = fmt.Sprintf("日期格式不正确: %s", err)
		ctx.Status = 400
		ctx.Render("home.html")
		return
	}

	tmpIsContainHead := ctx.Request.Form.Get("isContainHead")

	var isContainHead bool
	if tmpIsContainHead == "on" {
		isContainHead = true
	} else {
		isContainHead = false
	}

	file, handler, err := ctx.Request.FormFile("uploadfile")
	if err != nil {
		ctx.FlushMessage = fmt.Sprintf("上传错误: %s", err)
		ctx.Status = 400
		ctx.Render("home.html")
		return
	}

	x, err := newXlsx(currentDir, file, handler, date, isContainHead)
	if err != nil {
		ctx.FlushMessage = fmt.Sprintf("上传错误: %s", err)
		ctx.Status = 400
		ctx.Render("home.html")
		return
	}
	if err = x.save(); err != nil {
		ctx.FlushMessage = fmt.Sprintf("保存文件错误: %s", err)
		ctx.Status = 400
		ctx.Render("home.html")
		return
	}

	ctx.FlushMessage = filepath.Join("上传完成")
	ctx.Render("home.html")
}

func (h *Handler) Download(ctx *Context) {
	ctx.Render("download.html")
}

func (h *Handler) DownloadDo(ctx *Context) {
	if ctx.Method != "POST" {
		ctx.Redirect("/download")
		return
	}

	tmpDate := ctx.Request.Form.Get("date")

	date, err := time.Parse("2006-01", tmpDate)

	if err != nil {
		ctx.FlushMessage = fmt.Sprintf("日期格式不正确: %s", err)
		ctx.Status = 400
		ctx.Render("download.html")
		return
	}
	a := newAttendance(filepath.Join(currentDir, "data"), date)
	file, err := a.getXlsx()
	if err != nil {
		ctx.FlushMessage = fmt.Sprintf("处理文件出错: %s", err)
		ctx.Status = 400
		ctx.Render("download.html")
	}

	ctx.Download(file)
}
