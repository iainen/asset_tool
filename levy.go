package main

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"log"
)

const (
	WetestAssetTag="assetid"
	WetestModel="model"
	WetestProudct="product"
	WetestBrand="brand"
	WetestManu="manu"
	WetestSerial="serial"
	WetestIMEI="imei"
	WetestPc="pc"
)

var WeTestTitles = []string{
	WetestAssetTag, WetestModel, WetestProudct, WetestBrand, WetestManu, WetestSerial, WetestIMEI, WetestPc,
}

// 资产编号(assetid)-机型(model)-(product)-品牌(brand)-厂商(manu)-pos
type WetestAsset struct {
	AssetTag string
	Model    string
	Product  string
	Brand    string
	Manu     string

	Serial   string
	IMEI     string
	Pc       string

	FullName string
}

// 资产编号-> 手机信息表
func loadWetestGoodExcel2Map(inputFile string) (map[string]*WetestAsset, error) {
	f, err := excelize.OpenFile(inputFile)
	if err != nil {
		log.Fatalf("fail to open file: %v", err)
	}

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)

	titleRow := rows[0]
	titleMap := detainTitles(WeTestTitles, titleRow)
	log.Printf("len:%v, %v", len(titleRow), titleRow)
	log.Printf("%v", titleMap)

	var outMap = make(map[string]*WetestAsset)
	for _, row := range rows[1:] {
		tag := row[titleMap[WetestAssetTag]]
		model := row[titleMap[WetestModel]]
		product := row[titleMap[WetestProudct]]
		brand := row[titleMap[WetestBrand]]
		manu := row[titleMap[WetestManu]]
		serial := row[titleMap[WetestSerial]]
		imei := row[titleMap[WetestIMEI]]
		pc := row[titleMap[WetestPc]]

		//log.Printf("tag:[%v], model:[%v], prod:[%v], manu:[%v], serial:[%v], imei:[%v], pc:[%v]",
		//	tag, model, product, manu, serial, imei, pc)

		item, ok := outMap[tag]
		if !ok {
			outMap[tag] = &WetestAsset{
				AssetTag: tag,
				Model: model,
				Product: product,
				Brand: brand,
				Manu: manu,
				Serial: serial,
				IMEI: imei,
				Pc: pc,
			}
		} else {
			log.Fatalf("error: duplicated asset tag:%v", item.AssetTag)
		}
	}

	return outMap, nil
}
