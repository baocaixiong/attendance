package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
)

func main() {
	excelFileName := "/Users/baocaixiong/Downloads/技术部9月.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Println(err)
	}
	for _, sheet := range xlFile.Sheets {
		for i, row := range sheet.Rows {
			fmt.Println(row, len(sheet.Rows))
			for ii, cell := range row.Cells {
				fmt.Printf("%d--%s\n", ii, cell.String())
			}
			if i == 4 {
				return
			}
		}
	}
}
