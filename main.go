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

	// step3:
	//  model->aliasname

	/*
	var benchNameMap map[string][]Data
	var err error

	if *benchmark != "" {
		benchNameMap, err = loadBenchmark2NameMap(*benchmark)
		if err != nil {
			return
		}
	}

	if *levy4500 != "" {
		maps, err := loadFrom4500(*levy4500)
		if err != nil {
			log.Printf("fail to load levy excel: %v", err)
			return
		}

		if benchNameMap != nil {
			nameList := make([]string, 0, len(benchNameMap))
			for name, _ := range benchNameMap {
				if len(name) >= 2 { // some bad typename
					nameList = append(nameList, name)
				}
			}

			//sort.Strings(nameList)
			//for i, v := range nameList {
			//	log.Printf("%v: %v", i, v)
			//}
			//return

			for _, item := range maps {
				for _, asset := range item.Assets {
					//log.Printf("  %v", asset.FullName)
					found := false
					for i := len(nameList)-1; i >= 0; i-- {
						if strings.Contains(strings.ToLower(asset.FullName), strings.ToLower(nameList[i])) {
							found = true
							//log.Printf("[Y] %v --> %v", nameList[i], asset.FullName)
							break
						}
					}

					if !found {
						log.Printf("[X] %v", asset.FullName)
					}
				}
			}
		}
	}

	if *out != "" {
		benchModelMap, _ := loadBenchmark2ModelMap(*benchmark)
		exportExcel(*out, benchModelMap)
	}
	 */
}
