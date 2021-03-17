package main

import (
	"flag"
	"log"
)

var benchmark = flag.String("benchmark", "", "benchmark json file")
var out = flag.String("output", "", "export to a excel file")
var epo = flag.String("epo", "", "加载解析epo导出的excel表")
var wetest = flag.String("wetest", "", "加载解析levy导出的数据库表")

// tag->fullname
var epoAssetsMap map[string]EpoAsset

// tag->model, tag->model, fullname
var wetestGoodAsset map[string]*WetestAsset

// model, fullname
var wetestModelMap map[string]*ModelDetail

// tag->fullname
var wetestBadAsset = make(map[string]*WetestAsset, 0)

// fullname->model
var wetestGoodFullname2Model map[string]string

func showWetestAsset(assets map[string]*WetestAsset) {
	index := 1
	for tag, item := range assets {
		log.Printf("%v: [tag]: %v, [name]:%v", index, tag, item.FullName)
		index++
	}
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

	// step1: 生成两张表
	//	wetestGoodAsset: tag->(model,fullname)
	//  wetestBadAsset:  tag->fullname
	for tag, item := range epoAssetsMap {
		if wetestItem, ok := wetestGoodAsset[tag]; ok {
			wetestItem.FullName = item.Name
		} else { // not found
			wetestBadAsset[tag] = & WetestAsset{
				AssetTag: tag,
				FullName: item.Name,
				ModelDetail: ModelDetail{
					Manu: item.Manu,
				},
			}
		}
	}

	// step2: 处理epo中缺少机型的哪些资产
	// wetestGoodFullname2Model: fullname->model
	wetestGoodFullname2Model = make(map[string]string)
	for _, item := range wetestGoodAsset {
		wetestGoodFullname2Model[item.FullName] = item.Model
	}

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

	// showWetestAsset(wetestBadAsset) //312个
	nameMap := epoAsset2NameMap(wetestBadAsset)
	showEpoNameMap(nameMap)

	exportEpoNameMap("209.xlsx", nameMap)

	// step3:
	//  model->aliasname
}
