package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type _xlsx struct {
	file          multipart.File
	fileHeader    *multipart.FileHeader
	date          time.Time
	isContainHead bool
	parentDir     string
	fileName      string
}

func newXlsx(parentDir string, file multipart.File,
	fileHeader *multipart.FileHeader, date time.Time,
	isContainHead bool) (*_xlsx, error) {
	x := &_xlsx{
		file:          file,
		fileHeader:    fileHeader,
		date:          date,
		parentDir:     parentDir,
		isContainHead: isContainHead,
	}
	x.fileName = x.getFileName()

	if !x.checkExt() {
		return nil, errors.New("不允许的上传文件类型")
	}

	return x, nil
}

func (x *_xlsx) save() error {
	err := x.saveOriginFile()
	if err != nil {
		return err
	}
	err = x.saveAsJsonFile()
	if err != nil {
		return err
	}

	return nil
}

func (x *_xlsx) saveOriginFile() error {
	filename := x.fileName + x.getFileExt()

	fullpath := filepath.Join(x.parentDir, "origin", filename)

	basedir := filepath.Dir(fullpath)

	if _, err := os.Stat(basedir); os.IsNotExist(err) {
		os.MkdirAll(basedir, 0770)
	}

	f, _ := os.OpenFile(fullpath, os.O_CREATE|os.O_WRONLY, 0660)
	_, err := io.Copy(f, x.file)
	if err != nil {
		return errors.New(fmt.Sprintf("上传失败: %s", err))
	}

	return nil
}

func (x *_xlsx) getFileExt() string {
	return filepath.Ext(x.fileHeader.Filename)
}

func (x *_xlsx) checkExt() bool {
	extensions := []string{".xlsx"}
	ext := x.getFileExt()

	for _, v := range extensions {
		if v == ext {
			return true
		}
	}
	return false
}

func (x *_xlsx) getFileName() string {
	date := x.date
	return filepath.Join(strconv.Itoa(date.Year()),
		strconv.Itoa(int(date.Month())), strconv.Itoa(date.Day()))
}

func (x *_xlsx) saveAsJsonFile() error {
	xlFile, err := xlsx.OpenFile(filepath.Join(x.parentDir, "origin", x.fileName) + x.getFileExt())

	if err != nil {
		return errors.New(fmt.Sprintf("解析文件出错: %s", err))
	}

	rs := new(_rows)

	for _, sheet := range xlFile.Sheets {
		if len(sheet.Rows) < 2 {
			continue
		}
		for index, row := range sheet.Rows {
			if index == 0 || len(row.Cells) < 2 || row.Cells[0].String() == "" {
				continue
			}
			r := &_row{
				Name:       row.Cells[0].String(),
				Department: row.Cells[1].String(),
				Begin:      row.Cells[2].String(),
				End:        row.Cells[3].String(),
			}

			if row.Cells[4].String() != "" {
				r.End = row.Cells[4].String()
			}

			rs.addRow(r)
		}
	}
	js, err := json.Marshal(rs)

	if err != nil {
		return errors.New(fmt.Sprintf("解析文件到Json出错: %s", err))
	}

	filename := x.fileName + ".json"

	fullpath := filepath.Join(x.parentDir, "data", filename)

	basedir := filepath.Dir(fullpath)

	if _, err := os.Stat(basedir); os.IsNotExist(err) {
		os.MkdirAll(basedir, 0770)
	}

	f, _ := os.OpenFile(fullpath, os.O_CREATE|os.O_WRONLY, 0660)
	f.WriteString(string(js))
	if err != nil {
		return errors.New(fmt.Sprintf("写入文件出错: %s", err))
	}

	return nil
}

type _row struct {
	Name       string `json:"name"`
	Department string `json:"department"`
	Begin      string `json:"begin"`
	End        string `json:"end"`
}

type _rows struct {
	Rows []*_row
}

func (rs *_rows) addRow(r *_row) {
	rs.Rows = append(rs.Rows, r)
}
