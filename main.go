package main

import (
	"flag"
	"log"
)

var benchmark = flag.String("benchmark", "", "benchmark json file")
var levy4500 = flag.String("levy4500", "", "levy 4500 xlsx file")
var out = flag.String("output", "", "export to a excel file")
var epo = flag.String("epo", "", "加载解析epo导出的excel表")
var wetest = flag.String("wetest", "", "加载解析levy导出的数据库表")

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile | log.LstdFlags)



	if *epo != "" {
		_, _ = loadEpoExcel2AssetMap(*epo)
	}

	if *wetest != "" {
		_, _ = loadWetestGoodExcel2Map(*wetest)
	}

	// step1:
	//	tag->model->fullname
	//	model->fullname
	//	fullname->model

	// step2:
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
