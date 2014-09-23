package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
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
	tmpDate := ctx.Request.Form.Get("date")

	date, err := time.Parse("2006-01", tmpDate)
	if erro != nil {
		ctx.Status = 400
		ctx.flashData = []byte("日期格式不正确")
		ctx.Render("home")
	}

	tmpIsContainHead = ctx.Request.Form.Get(("isContainHead"))

	if tmpIsContainHead == "on" {
		isContainHead = true
	} else {
		isContainHead = false
	}
}

func (h *Handler) Download(ctx *Context) {
	ctx.Render("download")
}

func (h *Handler) DownloadDo(ctx *Context) {

}
