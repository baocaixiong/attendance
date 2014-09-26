package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"time"
)

type attendanceSheet struct {
	date time.Time
}

func newAttendanceSheet(date time.Time, sheet *xlsx.Sheet) *attendanceSheet {
	as := new(attendanceSheet)
	as.date = date

	as.init()
	return as
}

func (as *attendanceSheet) init() {
	// 写入初始化表头
	monthNumner := getMonthDays(as.date) //本月总天数
	xlsxFileCols := monthNumner + 16

	dataArray := make([][]string, xlsxFileCols)

	for row, col := range dataArray {
		col = make([]string, 100)
		fmt.Println(row, col)
	}
}

func (as *attendanceSheet) writeWeek() {

}
