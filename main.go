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

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
