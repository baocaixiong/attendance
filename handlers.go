package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

	ctx.Request.ParseForm()

	fmt.Println(ctx.Request.Form)

	tmpDate := ctx.Request.Form.Get("date")

	date, err := time.Parse("2006-01", tmpDate)
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

	ctx.Request.ParseMultipartForm(32 << 20)
	file, handler, err := ctx.Request.FormFile("uploadfile")
	if err != nil {
		ctx.FlushMessage = fmt.Sprintf("上传错误: %s", err)
		return
	}

	fileext := filepath.Ext(handler.Filename)
	if check(fileext) == false {
		ctx.FlushMessage = "不允许的上传类型"
		return
	}

	updateDir := filepath.Join(currentDir, "data")

	filename := strconv.FormatInt(time.Now().Unix(), 10) + fileext
	fullpath := filepath.Join(updateDir, filename)

	f, _ := os.OpenFile(fullpath, os.O_CREATE|os.O_WRONLY, 0660)
	_, err = io.Copy(f, file)
	if err != nil {
		ctx.FlushMessage = fmt.Sprintf("上传失败: %s", err)
		return
	}

	ctx.FlushMessage = filepath.Join(filename + "上传完成,服务器地址:" + fullpath)

	fmt.Println(isContainHead, date, f)
}

func (h *Handler) Download(ctx *Context) {
	ctx.Render("download.html")
}

func (h *Handler) DownloadDo(ctx *Context) {

}

func check(name string) bool {
	ext := []string{".xlsx"}

	for _, v := range ext {
		if v == name {
			return true
		}
	}
	return false
}
