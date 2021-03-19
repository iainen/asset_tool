package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"log"
)

const (
	EpoAssetTag = "资产编码"
	EpoBrand    = "品牌名称"
	EpoName     = "规格型号"
)

var EpoTitles = []string{
	EpoAssetTag, EpoBrand, EpoName,
}

type EpoAsset struct {
	AssetTag string
	Brand    string
	Name     string //very long name
}

// 功能：将epm导出的资产编号excel表转换为资产字典表
func loadEpoExcel2AssetMap(inputFile string) (map[string]EpoAsset, error) {
	f, err := excelize.OpenFile(inputFile)
	if err != nil {
		log.Fatalf("fail to open file: %v", err)
	}

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)

	titleRow := rows[0]
	titleMap := detainTitles(EpoTitles, titleRow)

	var epoMap = make(map[string]EpoAsset)
	for _, row := range rows[1:] {
		tag := row[titleMap[EpoAssetTag]]
		brand := row[titleMap[EpoBrand]]
		name := row[titleMap[EpoName]]

		//log.Printf("tag:[%v], brand:[%v], name:[%v]", tag, brand, name)
		item, ok := epoMap[tag]
		if !ok {
			epoMap[tag] = EpoAsset{
				AssetTag: tag,
				Brand: brand,
				Name: name,
			}
		} else {
			log.Fatalf("error: duplicated asset tag:%v", item.AssetTag)
		}
	}

	return epoMap, nil
}

// fullname->brand
func epoAssetName2BrandMap(assets map[string]*WetestAsset) map[string]string {
	nameMap := make(map[string]string, 0)
	for _, asset := range assets {
		if _, ok := nameMap[asset.FullName]; !ok {
			nameMap[asset.FullName] = asset.EpoBrand
		}
	}

	return nameMap
}

func showEpoNameMap(nameMap map[string]string) {
	index := 1
	for name, manu := range nameMap {
		log.Printf("%v: [name]: %v, [manu]:%v", index, name, manu)
		index++
	}
}

func exportEpoNameMap(outfile string, nameMap map[string]string) {
	f := excelize.NewFile()

	sheet := "Sheet"
	// Create a new sheet.
	index := f.NewSheet(sheet)
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)

	titles := []string{
		EpoName,
		EpoBrand,
		WetestModel,
	}

	for i, v := range titles {
		f.SetCellValue(sheet, fmt.Sprintf("%c%d", int('A')+i, 1), v)
	}

	line := 2
	for name, manu := range nameMap {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", line), name)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", line), manu)
		line++
	}

	// Save spreadsheet by the given path.
	if err := f.SaveAs(outfile); err != nil {
		fmt.Println(err)
	}
}

func exportBrandMap(outfile string, brandMap map[string]string) {
	f := excelize.NewFile()

	sheet := "Sheet"
	// Create a new sheet.
	index := f.NewSheet(sheet)
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)

	titles := []string{
		WetestBrand,
		WetestManu,
	}

	for i, v := range titles {
		f.SetCellValue(sheet, fmt.Sprintf("%c%d", int('A')+i, 1), v)
	}

	line := 2
	for brand, item := range brandMap {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", line), brand)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", line), item)
		line++
	}

	// Save spreadsheet by the given path.
	if err := f.SaveAs(outfile); err != nil {
		fmt.Println(err)
	}
}

func exportManuMap(outfile string, manuMap map[string]string) {
	f := excelize.NewFile()

	sheet := "Sheet"
	// Create a new sheet.
	index := f.NewSheet(sheet)
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)

	titles := []string{
		WetestManu,
	}

	for i, v := range titles {
		f.SetCellValue(sheet, fmt.Sprintf("%c%d", int('A')+i, 1), v)
	}

	line := 2
	for item,_ := range manuMap {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", line), item)
		line++
	}

	// Save spreadsheet by the given path.
	if err := f.SaveAs(outfile); err != nil {
		fmt.Println(err)
	}
}

// fullname->model
func loadEpoNameModelExcel(inputFile string) (map[string]string, error) {
	var epoNameTitles = []string{
		EpoName, EpoBrand, WetestModel,
	}

	f, err := excelize.OpenFile(inputFile)
	if err != nil {
		log.Fatalf("fail to open file: %v", err)
	}

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)

	titleRow := rows[0]
	titleMap := detainTitles(epoNameTitles, titleRow)
	//log.Printf("len:%v, %v", len(titleRow), titleRow)
	//log.Printf("%v", titleMap)

	var nameModelMap = make(map[string]string)
	for _, row := range rows[1:] {
		name := row[titleMap[EpoName]]
		model := row[titleMap[WetestModel]]
		// log.Printf("%v-->%v", name, model)
		if _, ok := nameModelMap[name]; !ok {
			nameModelMap[name] = model
		} else {
			log.Fatalf("error: duplicated asset name:%v", name)
		}
	}

	return nameModelMap, nil
}

// model->detail
func loadEpoNewModelExcel(inputFile string, sheetIndex int, input map[string]*ModelDetail) {
	var epoNameTitles = []string{
		EpoName, WetestModel, WetestBrand, WetestManu,
	}

	f, err := excelize.OpenFile(inputFile)
	if err != nil {
		log.Fatalf("fail to open file: %v", err)
	}

	sheetName := f.GetSheetName(sheetIndex)
	rows, err := f.GetRows(sheetName)

	titleRow := rows[0]
	titleMap := detainTitles(epoNameTitles, titleRow)

	index := 0
	for _, row := range rows[1:] {
		model := row[titleMap[WetestModel]]
		brand := row[titleMap[WetestBrand]]
		manu := row[titleMap[WetestManu]]
		name := row[titleMap[EpoName]]

		// TODO: 读取第二个表时，不只为何会读空行，这里判断是否读完所有行
		if model == "" && brand == "" && manu == "" {
			break
		}

		index++
		// log.Printf("%v: %v-->%v", index, name, model)
		if _, ok := input[model]; !ok {
			input[model] = &ModelDetail{
				Model: model,
				Brand: brand,
				Manu:  manu,
				AliasName: name,
			}
		} else {
			log.Fatalf("error: duplicated asset model:%v, name:%v", model, name)
		}
	}
}
