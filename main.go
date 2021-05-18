package main

import (
	"flag"
	"log"
)

var benchmark = flag.String("benchmark", "", "benchmark json file")
var epo = flag.String("epo", "", "加载解析epo导出的excel表")
var wetest = flag.String("wetest", "", "加载解析levy导出的数据库excel表")
var nameModel = flag.String("nameModel", "", "加载线下名称-机型excel表")
var ctModel = flag.String("ctModel", "", "加载线下名称-机型excel表")

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
	// 输出: tag->name表，wetestBadAsset: 没有机型信息，500个左右

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
	// 输出：tag->fullname，wetestBadAsset: 新机器，但是名称可能有重复，，312个
	// showWetestAsset(wetestBadAsset) //312个

	// 导出详细的机型信息表
	wetestGoodFullname2Model, wetestGoodBrand, wetestGoodManu = filterGoodAsset(wetestGoodAsset)
	exportEpoNameMap("name.xlsx", wetestGoodFullname2Model)
	exportBrandMap("brand.xlsx", wetestGoodBrand)
	exportManuMap("manu.xlsx", wetestGoodManu)

	// 输出：badNameMap: 没有机型信息，且去除重复名称，209个
	badNameMap := epoAssetName2BrandMap(wetestBadAsset)
	//showEpoNameMap(badNameMap)

	// 导出文件，该文件需要手动补充机型信息
	exportEpoNameMap("209.xlsx", badNameMap)

	// step3: 将线下录入的 模型-厂商-名称 合并回
	loadEpoNewModelExcel(*nameModel, 1, wetestModelMap)

	// 加载手动编辑的全名-模型名表，但是没有模型的详细信息，209个
	// 输出，模型-全名表，模型不重复，这些是新增的模型，需要手动补齐
	wetestHandFullname2Model, _ = loadEpoNameModelExcel(*nameModel)
	newModel := make(map[string]string)
	for name, model := range wetestHandFullname2Model {
		if _, ok := wetestModelMap[model]; !ok {
			log.Printf("%v,%v\n", model, name)
			newModel[model] = name
		}
	}
	if len(newModel) != 0 {
		exportEpoNameMap("newModel.xlsx", newModel)
	}

	// 对312个设备，利用名称相等，获得机型；利用机型获得机型详细信息
	for tag, item := range wetestBadAsset {
		// 匹配名称信息
		if model,ok := wetestHandFullname2Model[item.FullName]; ok {
			// 匹配机型信息机型
			if detail, ok := wetestModelMap[model]; ok {
				wetestGoodAsset[tag] = item
				// log.Printf("find : %v,%v\n", model, item.FullName)
				item.Model = model
				item.Product = detail.Product
				item.Brand = detail.Brand
				item.Manu = detail.Manu
			} else {
				log.Fatalf("error: not find model:%v, name:%v", model, tag)
			}
		} else {
			log.Fatalf("error: not find tag:%v, name: %v", tag, item.FullName)
		}
	}

	// step4: 生成所有数据的总表
	exportSnipITExcel("all.xlsx", wetestGoodAsset)

	// step5: 为总表配置别名
	// benchmarkMap model->aliasname
	benchmarkModelMap, _ = loadBenchmark2ModelMap(*benchmark)

	ctModelMap := make(map[string]*ModelDetail)
	loadCtModelExcel(*ctModel, 0, ctModelMap)

	out := threeCheckModelMap(wetestGoodAsset, benchmarkModelMap, ctModelMap)
	exportModelDetail("alias.xlsx", out)
}


type ModelDetail2 struct {
	ModelDetail
	AliasName1 string // fullname
	AliasName2 string // from benchmark
	AliasName3 string // from ct
}

func threeCheckModelMap(assets map[string]*WetestAsset, bench map[string]Data, ct map[string]*ModelDetail) map[string]*ModelDetail2 {
	out := make(map[string]*ModelDetail2)
	for _, item := range assets {
		model := item.Model

		detail, ok := wetestModelMap[model]
		if !ok {
			log.Fatalf("error: not find model:%v, name:%v", model, item.AssetTag)
		}

		if _, ok := out[model]; !ok {
			out[model] = &ModelDetail2{
				ModelDetail: *detail,
				AliasName1:  item.FullName,
			}
		}

		if out[model].AliasName2 == "" {
			if item, ok := bench[model]; ok {
				out[model].AliasName2 = item.Name
			}
		}

		if out[model].AliasName3 == "" {
			if item, ok := ct[model]; ok {
				out[model].AliasName3 = item.AliasName
			}
		}
	}

	return out
}