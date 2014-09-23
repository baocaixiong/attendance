package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	tpls := newTpls()
	handler := newHandler(tpls)

	http.HandleFunc("/", handler.handle)
	http.HandleFunc("/assets/", Static)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
