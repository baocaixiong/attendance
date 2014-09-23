package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	CONTEXT_RENDERED = "context_rendered"
	CONTEXT_END      = "context_end"
	CONTEXT_SEND     = "context_send"

	FLASH_ERROR   = "error"
	FLASH_SUCCESS = "success"
	FLASH_WARNING = "warning"
)

var currentDir = func() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}()

var logger = log.New(os.Stdout, "", log.LstdFlags)
