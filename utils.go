package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

var currentDir = func() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}()

var logger = log.New(os.Stdout, "", log.LstdFlags)

func getMonthDays(date time.Time) int {
	t := time.Date(date.Year(), date.Month(), 32, 0, 0, 0, 0, time.UTC)
	return 32 - t.Day()
}

func makeTmpDir() string {
	tmpDir := filepath.Join(currentDir, "tmp")
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		os.MkdirAll(tmpDir, 0777)
	}

	return tmpDir
}
