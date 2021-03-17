package main

import (
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
