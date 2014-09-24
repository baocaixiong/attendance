package main

import (
	"fmt"
	"github.com/baocaixiong/attendance/tpls"
	"net/http"
	"os"
)

func main() {
	tpls.Parse()

	handler := newHandler(tpls.T)

	http.HandleFunc("/", handler.handle)

	logger.Println("开启服务, 请在浏览器访问 http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
