package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

			day, _ := strconv.Atoi(filepath.Base(x.fileName))
			r.Day = day

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
	Day        int    `json:"day"`
}

type _rows struct {
	Rows []*_row
}

func (rs *_rows) addRow(r *_row) {
	rs.Rows = append(rs.Rows, r)
}

type _attendance struct {
	parentDir string
	date      time.Time
	fileList  map[string]*_rows
}

func newAttendance(parentDir string, date time.Time) *_attendance {
	a := new(_attendance)
	a.date = date
	a.parentDir = parentDir
	a.fileList = make(map[string]*_rows)
	return a
}

func (a *_attendance) getXlsx() (string, error) {
	err := a.parseAttendanceContent()
	if err != nil {
		return "", err
	}

	ds := a.sortedByDepartment()

	zipFile, err := a.writeResult(ds)

	if err != nil {
		return "", err
	}

	return zipFile, nil
}

func (a *_attendance) parseAttendanceContent() error {

	jsonPath := a.getJsonPath()
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("从未上传该月份文件或者目录错误: %s", err))
	}

	filepath.Walk(a.getJsonPath(), func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".json") {
			return nil
		}

		f, fpOpenErr := os.Open(path)
		if fpOpenErr != nil {
			logger.Fatalln("Can not open file %v", fpOpenErr)
			return fpOpenErr
		}

		defer f.Close()

		bReader := bufio.NewReader(f)
		out := []byte{}
		for {
			buffer := make([]byte, 1024)
			readCount, readErr := bReader.Read(buffer)
			if readErr == io.EOF {
				break
			} else {
				out = append(out, buffer[:readCount]...)
			}
		}

		rs := new(_rows)

		json.Unmarshal(out, rs)
		a.fileList[strings.Split(filepath.Base(path), ".")[0]] = rs

		return err
	})

	return nil
}

func (a *_attendance) sortedByDepartment() *departments {
	ds := new(departments)
	ds.d = make(map[string]*department)

	for _, rs := range a.fileList {
		for _, r := range rs.Rows {
			var d *department
			if !ds.has(r.Department) {
				d = new(department)
				d.rows = make(map[string]*_row)
				d.name = r.Department
				ds.add(d)
			} else {
				d = ds.get(r.Department)
			}
			d.add(r)
		}
	}

	return ds
}

func (a *_attendance) writeResult(ds *departments) (string, error) {
	files := make(map[string]map[string]string)

	for _, d := range ds.d {
		file := xlsx.NewFile()
		fileName := d.name + "-" + strconv.Itoa(int(a.date.Month())) + "月.xlsx"
		sheet := file.AddSheet(fileName) // add sheet
		newAttendanceSheet(a.date, sheet)
		parts, err := file.MarshallParts()
		if err != nil {
			return "", err
		}
		files[fileName] = parts
	}

	resultPath := filepath.Join(makeTmpDir(), time.Now().Local().Format("2006年01月02号150405")+".zip")

	target, err := os.Create(resultPath)

	if err != nil {
		return "", errors.New(fmt.Sprintf("创建考勤文件出错, %s", err))
	}

	w := zip.NewWriter(target)

	for name, parts := range files {
		f, err := w.Create(name)
		if err != nil {
			return "", err
		}

		b := new(bytes.Buffer)
		zw := zip.NewWriter(b)
		for partName, part := range parts {
			var writer io.Writer
			writer, err = zw.Create(partName)
			if err != nil {
				return "", err
			}
			_, err = writer.Write([]byte(part))
			if err != nil {
				return "", err
			}
		}
		zw.Close()

		_, err = f.Write(b.Bytes())
		if err != nil {
			return "", err
		}
	}

	err = w.Close()
	if err != nil {
		return "", err
	}

	return resultPath, nil
}

func (a *_attendance) defaultRows() {

}

func (a *_attendance) getJsonPath() string {
	return filepath.Join(a.parentDir, strconv.Itoa(a.date.Year()),
		strconv.Itoa(int(a.date.Month())))
}

type department struct {
	rows map[string]*_row
	name string
}

func (d *department) has(name string) bool {
	_, ok := d.rows[name]
	return ok
}

func (d *department) add(r *_row) {
	d.rows[r.Name] = r
}

type departments struct {
	d map[string]*department
}

func (ds *departments) has(name string) bool {
	_, ok := ds.d[name]
	return ok
}

func (ds *departments) get(name string) *department {
	return ds.d[name]
}

func (ds *departments) add(d *department) {
	ds.d[d.name] = d
}
