package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func exportExcel(outfile string, modelMap map[string]Data) {
	f := excelize.NewFile()
	// Create a new sheet.
	index := f.NewSheet("Sheet")
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)

	titles := []string{
		"型号",
		"手机名",
		"厂商",
		"CPU名",
		"异性屏",
		"分辨率",
		"上市日期",
	}

	for i, v := range titles {
		//log.Printf("%c%d:%v", int('A')+i, 1, v)
		f.SetCellValue("Sheet", fmt.Sprintf("%c%d", int('A')+i, 1), v)
	}

	line := 2
	for model, item := range modelMap {
		f.SetCellValue("Sheet", fmt.Sprintf("A%d", line), model)
		f.SetCellValue("Sheet", fmt.Sprintf("B%d", line), item.Name)
		f.SetCellValue("Sheet", fmt.Sprintf("C%d", line), item.Manu)
		//f.SetCellValue("Sheet", fmt.Sprintf("D%d", line), name.Cpuname)
		//f.SetCellValue("Sheet", fmt.Sprintf("E%d", line), name.Notch)
		//f.SetCellValue("Sheet", fmt.Sprintf("F%d", line), name.ScreenRatio)
		//f.SetCellValue("Sheet", fmt.Sprintf("G%d", line), name.Date)
		line++
	}

	// Save spreadsheet by the given path.
	if err := f.SaveAs(outfile); err != nil {
		fmt.Println(err)
	}
}