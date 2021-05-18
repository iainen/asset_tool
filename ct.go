package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"log"
)

// model->detail
func loadCtModelExcel(inputFile string, sheetIndex int, input map[string]*ModelDetail) {
	var epoNameTitles = []string{
		WetestModel, WetestManu, "another_name",
	}

	f, err := excelize.OpenFile(inputFile)
	if err != nil {
		log.Fatalf("fail to open file: %v", err)
	}

	sheetName := f.GetSheetName(sheetIndex)
	rows, err := f.GetRows(sheetName)

	titleRow := rows[0]
	titleMap := detainTitles(epoNameTitles, titleRow)
	// log.Printf("titile: %v, map:%v", titleRow, titleMap)

	index := 0
	for _, row := range rows[1:] {
		model := row[titleMap[WetestModel]]
		manu := row[titleMap[WetestManu]]

		var name string
		if len(row) > 2 {
			name = row[titleMap["another_name"]]
		}

		if model == "" && manu == "" {
			break
		}

		index++
		// log.Printf("%v: %v-->%v", index, name, model)
		if _, ok := input[model]; !ok {
			input[model] = &ModelDetail{
				Model: model,
				Manu:  manu,
				AliasName: name,
			}
		} else {
			log.Fatalf("error: duplicated asset model:%v, name:%v", model, name)
		}
	}
}

func exportModelDetail(outfile string, modelMap map[string]*ModelDetail2) {
	f := excelize.NewFile()
	sheet := f.GetSheetName(0)

	titles := []string{
		"厂商",
		"型号",
		"品牌",
		"别名1_ota全名",
		"别名2_bench别名",
		"别名3_终端云测别名",
	}

	for i, v := range titles {
		f.SetCellValue(sheet, fmt.Sprintf("%c%d", int('A')+i, 1), v)
	}

	line := 2
	for _, item := range modelMap {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", line), item.Manu)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", line), item.Model)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", line), item.Brand)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", line), item.AliasName1)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", line), item.AliasName2)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", line), item.AliasName3)
		line++
	}

	// Save spreadsheet by the given path.
	if err := f.SaveAs(outfile); err != nil {
		fmt.Println(err)
	}
}