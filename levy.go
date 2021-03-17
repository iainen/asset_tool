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

type ModelDetail struct {
	Model     string
	Product   string
	Brand     string
	Manu      string
	AliasName string   // 短名称，比如Mate 10
}

// 资产编号(assetid)-机型(model)-(product)-品牌(brand)-厂商(manu)-pos
type WetestAsset struct {
	ModelDetail
	AssetTag string
	Serial   string
	IMEI     string
	Pc       string

	EpoBrand string
	FullName string
}

// 资产编号-> 手机信息表
func loadWetestGoodExcel2Map(inputFile string) (map[string]*WetestAsset, map[string]*ModelDetail, error) {
	f, err := excelize.OpenFile(inputFile)
	if err != nil {
		log.Fatalf("fail to open file: %v", err)
	}

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)

	titleRow := rows[0]
	titleMap := detainTitles(WeTestTitles, titleRow)

	var modelMap = make(map[string]*ModelDetail)
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

		if item, ok := outMap[tag]; !ok {
			outMap[tag] = &WetestAsset{
				ModelDetail: ModelDetail{
					Model: model,
					Product: product,
					Brand: brand,
					Manu: manu,
				},
				AssetTag: tag,
				Serial: serial,
				IMEI: imei,
				Pc: pc,
			}
		} else {
			log.Fatalf("error: duplicated asset tag:%v", item.AssetTag)
		}

		if _, ok := modelMap[model]; !ok {
			modelMap[model] = & ModelDetail{
				Model: model,
				Product: product,
				Brand: brand,
				Manu: manu,
			}
		}
	}

	return outMap, modelMap, nil
}

func showWetestAsset(assets map[string]*WetestAsset) {
	index := 1
	for tag, item := range assets {
		log.Printf("%v: [tag]: %v, [name]:%v", index, tag, item.FullName)
		index++
	}
}
