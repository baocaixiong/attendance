package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Handler struct {
	tpls *Tpls
	log  *log.Logger
}

func newHandler(tpls *Tpls) *Handler {
	handler := new(Handler)
	handler.tpls = tpls
	handler.log = log.New(os.Stdout, "", log.LstdFlags)

	return handler
}

func (h *Handler) handle(w http.ResponseWriter, r *http.Request) {
	context := NewContext(w, r, h.tpls)
	h.log.Println(fmt.Sprintf("%s %s", r.Method, r.URL.String()))
	switch r.URL.String() {
	case "/":
		context.IsUpdate = true
		h.Home(context)
	case "/update/do":
		context.IsUpdate = true
		h.UpdateDo(context)
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
	ctx.Render("home")
}

func (h *Handler) NotFound(ctx *Context) {
	ctx.Response.Write([]byte("This page is not found."))
}

func (h *Handler) UpdateDo(ctx *Context) {
	if ctx.Method != "POST" {
		ctx.Redirect("/")
		ctx.End()
	}
	tmpDate := ctx.Request.Form.Get("date")

	date, err := time.Parse("2006-01", tmpDate)
	if err != nil {
		ctx.Flash("日期格式不正确")
		ctx.Status = 400
		ctx.Render("home")
	}

	tmpIsContainHead := ctx.Request.Form.Get(("isContainHead"))

	var isContainHead bool
	if tmpIsContainHead == "on" {
		isContainHead = true
	} else {
		isContainHead = false
	}

	ctx.Request.ParseMultipartForm(32 << 20)
	file, handler, err := ctx.Request.FormFile("attendanceFile")
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(ctx.Response, "%v", "上传错误")
		return
	}

	fileext := filepath.Ext(handler.Filename)
	if check(fileext) == false {
		fmt.Fprintf(ctx.Response, "%v", "不允许的上传类型")
		return
	}

	updateDir := filepath.Join(currentDir, "data")

	filename := strconv.FormatInt(time.Now().Unix(), 10) + fileext
	f, _ := os.OpenFile(updateDir+filename, os.O_CREATE|os.O_WRONLY, 0660)
	_, err = io.Copy(f, file)
	if err != nil {
		fmt.Fprintf(ctx.Response, "%v", "上传失败")
		return
	}
	filedir, _ := filepath.Abs(updateDir + filename)
	fmt.Fprintf(ctx.Response, "%v", filename+"上传完成,服务器地址:"+filedir)

	fmt.Println(isContainHead, date, f)
}

func (h *Handler) Download(ctx *Context) {
	ctx.Render("download")
}

func (h *Handler) DownloadDo(ctx *Context) {

}

func check(name string) bool {
	ext := []string{".exe", ".js", ".png"}

	for _, v := range ext {
		if v == name {
			return false
		}
	}
	return true
}
