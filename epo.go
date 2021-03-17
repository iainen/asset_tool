package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"log"
)

const (
	EpoAssetTag         = "资产编码"
	EpoManu             = "品牌名称"
	EpoName             = "规格型号"
)

var EpoTitles = []string{
	EpoAssetTag, EpoManu, EpoName,
}

type EpoAsset struct {
	AssetTag string
	Manu     string
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
	log.Printf("len:%v, %v", len(titleRow), titleRow)
	log.Printf("%v", titleMap)

	var epoMap = make(map[string]EpoAsset)
	for _, row := range rows[1:] {
		tag := row[titleMap[EpoAssetTag]]
		manu := row[titleMap[EpoManu]]
		name := row[titleMap[EpoName]]

		//log.Printf("tag:[%v], manu:[%v], name:[%v]", tag, manu, name)
		item, ok := epoMap[tag]
		if !ok {
			epoMap[tag] = EpoAsset{
				AssetTag: tag,
				Manu: manu,
				Name: name,
			}
		} else {
			log.Fatalf("error: duplicated asset tag:%v", item.AssetTag)
		}
	}

	return epoMap, nil
}

// fullname->brand
func epoAsset2NameMap(assets map[string]*WetestAsset) map[string]string {
	nameMap := make(map[string]string, 0)
	for _, asset := range assets {
		if _, ok := nameMap[asset.FullName]; !ok {
			nameMap[asset.FullName] = asset.Manu
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
	// Create a new sheet.
	index := f.NewSheet("Sheet")
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)

	titles := []string{
		EpoManu,
		EpoName,
		WetestModel,
		WetestProudct,
		WetestBrand,
		WetestManu,
	}

	for i, v := range titles {
		f.SetCellValue("Sheet", fmt.Sprintf("%c%d", int('A')+i, 1), v)
	}

	line := 2
	for name, manu := range nameMap {
		f.SetCellValue("Sheet", fmt.Sprintf("A%d", line), name)
		f.SetCellValue("Sheet", fmt.Sprintf("B%d", line), manu)
		line++
	}

	// Save spreadsheet by the given path.
	if err := f.SaveAs(outfile); err != nil {
		fmt.Println(err)
	}
}
