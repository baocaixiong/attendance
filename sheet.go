package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"time"
)

type attendanceSheet struct {
	date time.Time
	d    *department
	s    *xlsx.Sheet
}

func newAttendanceSheet(date time.Time, sheet *xlsx.Sheet, d *department) *attendanceSheet {
	as := new(attendanceSheet)
	as.date = date
	as.d = d
	as.s = sheet

	as.init()
	return as
}

func (as *attendanceSheet) init() {
	// 写入初始化表头
	monthNumner := getMonthDays(as.date) //本月总天数
	xlsxFileCols := monthNumner + 16

	for i := 0; i < xlsxFileCols; i++ {
		if i == 0 {
			r := as.s.AddRow()
			c := r.AddCell()
			c.Value = "姓名"
			c.SetStyle(as.warningStyle())
		} else if i == 1 {
			fmt.Println(i)
		} else {

			for _, row := range as.d.rows {
				fmt.Println(row)
			}
		}
	}
}

func (as *attendanceSheet) warningStyle() *xlsx.Style {
	s := &xlsx.Style{
		Fill: xlsx.Fill{
			BgColor: "0xf32fdc",
			FgColor: "0x000000",
		},
	}
	return s
}

func print1(s string) {
	fmt.Println(s)
}
