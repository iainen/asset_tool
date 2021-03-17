package main

import (
	"flag"
	"fmt"
	"log"
)

var benchmark = flag.String("benchmark", "", "benchmark json file")
var epo = flag.String("epo", "", "加载解析epo导出的excel表")
var wetest = flag.String("wetest", "", "加载解析levy导出的数据库excel表")
var nameModel = flag.String("nameModel", "", "加载线下名称-机型excel表")
var out = flag.String("output", "", "export to a excel file")

// tag->fullname
var epoAssetsMap map[string]EpoAsset

// tag->model, tag->model, fullname
var wetestGoodAsset map[string]*WetestAsset

// model->detail
var wetestModelMap map[string]*ModelDetail

// tag->fullname
var wetestBadAsset = make(map[string]*WetestAsset, 0)

// fullname->model
var wetestGoodFullname2Model map[string]string

// brand->manu
var wetestGoodBrand map[string]string

// manu
var wetestGoodManu map[string]string

// fullname->model，线下录入的全名到model
var wetestHandFullname2Model map[string]string

// benchmark's model->alias name
var benchmarkModelMap map[string]Data

type ManuBrand struct {
	Brand     string
	Manu      string
}

func filterGoodAsset(goodAssets map[string]*WetestAsset) (
	nameMap map[string]string,
	bandMap map[string]string,
	manuMap map[string]string) {

	nameMap = make(map[string]string)
	bandMap = make(map[string]string)
	manuMap = make(map[string]string)
	for _, item := range goodAssets {
		nameMap[item.FullName] = item.Model
		bandMap[item.Brand] = item.Manu
		manuMap[item.Manu] = ""
	}

	return nameMap, bandMap, manuMap
}

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// epoAssetsMap：tag->fullname
	if *epo != "" {
		epoAssetsMap, _ = loadEpoExcel2AssetMap(*epo)
	}

	// wetestGoodAsset: tag->model
	// wetestModelMap:  model->model_detail
	if *wetest != "" {
		wetestGoodAsset, wetestModelMap, _ = loadWetestGoodExcel2Map(*wetest)
	}

	// step1: 两张表
	//	更新wetestGoodAsset: tag->(model,fullName,epoBrand)
	//  生成wetestBadAsset:  tag->(fullname, epoBrand)
	for tag, item := range epoAssetsMap {
		if wetestItem, ok := wetestGoodAsset[tag]; ok {
			wetestItem.FullName = item.Name
			wetestItem.EpoBrand = item.Brand
		} else { // not found
			wetestBadAsset[tag] = & WetestAsset{
				AssetTag: tag,
				FullName: item.Name,
				EpoBrand: item.Brand,
			}
		}
	}

	// step2: 处理epo中缺少机型的哪些资产
	// wetestGoodFullname2Model: fullname->model
	wetestGoodFullname2Model, _, _ = filterGoodAsset(wetestGoodAsset)
	for tag, item := range wetestBadAsset {
		if model, ok := wetestGoodFullname2Model[item.FullName]; ok {
			wetestGoodAsset[tag] = item
			item.Model = model

			// 查找model->detail表
			detail := wetestModelMap[model]
			item.Product = detail.Product
			item.Brand = detail.Brand
			item.Manu = detail.Manu

			delete(wetestBadAsset, tag)
		}
	}

	// 导出详细的机型信息表
	wetestGoodFullname2Model, wetestGoodBrand, wetestGoodManu = filterGoodAsset(wetestGoodAsset)
	exportEpoNameMap("name.xlsx", wetestGoodFullname2Model)
	exportBrandMap("brand.xlsx", wetestGoodBrand)
	exportManuMap("manu.xlsx", wetestGoodManu)

	// showWetestAsset(wetestBadAsset) //312个
	// nameMap := epoAsset2NameMap(wetestBadAsset)
	// showEpoNameMap(nameMap)

	// exportEpoNameMap("209.xlsx", nameMap)

	// step3: 将线下录入合并回
	loadEpoNewModelExcel(*nameModel, 1, wetestModelMap)

	wetestHandFullname2Model, _ = loadEpoNameModelExcel(*nameModel)
	newModel := make(map[string]string)
	for name, model := range wetestHandFullname2Model {
		if _, ok := wetestModelMap[model]; !ok {
			fmt.Printf("%v,%v\n", model, name)
			newModel[model] = name
		}
	}
	//exportEpoNameMap("newModel.xlsx", newModel)
	
	// step4:
	// benchmarkMap model->aliasname
	// benchmarkModelMap, _ := loadBenchmark2ModelMap(*benchmark)
}
