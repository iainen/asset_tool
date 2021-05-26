package main

import (
	"fmt"
	"git.code.oa.com/zhongkaizhu/assets_manager/excel"
	"log"
	"strings"
	"testing"
)

func TestEamCsv(t *testing.T) {
	// iconv -f GB18030 -t UTF8 2021-05-19_个人实物查询_zhongkaizhu.csv > 0519-zhongkai.csv
	_, list2 := loadEamCsv("../ct/eam-0519-zhongkai.csv", "TKMB")
	for _, line := range list2 {
		log.Printf("%#v", line)
	}

	log.Printf("len: %v", len(list2))

	exportCsv("../ct/TKMB-eam-0519-zhongkai.csv", list2)
}

func TestEamTransferCsv(t *testing.T) {
	_, list2 := loadEamTransferCsv("../ct/eb/2021-03-08_181台/asset.csv", "TKMB")
	for _, line := range list2 {
		log.Printf("%#v", line)
	}
	log.Printf("len: %v", len(list2))
}

func TestTrim(t *testing.T) {
	fmt.Printf("[%q]\n", strings.TrimSpace("    "))
}

func TestGetAllFiles(t *testing.T) {
	all := mergeAssets("../ct/eb/", "asset.csv", "TKMB")
	// export all
	exportCsv("../ct/eb/all-asset.csv", &all)
}

func TestFilterSelfA1Assets(t *testing.T) {
	ebMap := make(map[string]*EamTransferLine)
	eb := mergeAssets("../ct/eb/", "asset.csv", "TKMB")
	for _, line := range eb {
		ebMap[line.AssetTag] = line
	}

	selfAssets := make([]*EamLine, 0)
	_, selfAll := loadEamCsv("../ct/eam-0519-zhongkai.csv", "TKMB")
	for _, line := range selfAll {
		if _, ok := ebMap[line.AssetTag]; !ok {
			selfAssets = append(selfAssets, line)
		}
	}

	exportCsv("../ct/eb/all-a1-zhongkai.csv", &selfAssets)
}

func TestA6Asset(t *testing.T) {
	toChecks := make([]*A6Line, 0)
	loadCsv("../ct/A6-0525-assets.csv", &toChecks)

	log.Printf("len: %v\n", len(toChecks))
	for _, line := range toChecks {
		if  line.AssetTag == "" {
			log.Printf("%#v", line)
		}
	}

	all := make([]*SnipeItLine, 0)
	loadCsv("../ct/custom-assets-report-2021-05-25-035754.csv", &all)

	log.Printf("len: %v\n", len(all))
	for _, line := range all {
		if  line.AssetTag == "" {
			log.Printf("%#v", line)
		}
	}

	allMap := make(map[string]*Line, 0)
	for _, line := range all {
		allMap[line.AssetTag] = &Line{
			Company:      line.Company,
			AssetTag:     line.AssetTag,
			Model:        line.Model,
			Brand:        line.Brand,
			Manufacturer: line.Manufacturer,
			Category:     line.Category,
			Status:       line.Status,
			Location:     line.Location,
		}
	}

	founded := make([]*Line, 0)
	notFounded := make([]*A6Line, 0)

	for _, line := range toChecks {
		if strings.TrimSpace(line.AssetTag) == "" {
			continue
		}

		if f, ok := allMap[line.AssetTag]; ok {
			f.Location = "A6_Checked"
			f.Status = "上线"

			//log.Printf("found: %#v", f)
			founded = append(founded, f)
		} else {
			notFounded = append(notFounded, line)
			log.Printf("not found: %#v", line)
		}
	}

	exportCsv("a6_founded.csv", &founded)
	if len(notFounded) > 0 {
		exportCsv("a6_not_founded.csv", &notFounded)
	}
}

func TestAssignSnipeItAsset(t *testing.T) {
	eamList := make([]*EamLine, 0)
	loadCsv("2021-05-25_个人实物查询_colerwang.csv", &eamList)
	allMap := make(map[string]*EamLine, 0)

	for _, line := range eamList {
		allMap[line.AssetTag] = line
	}

	notFound := make([]*SnipeItLine, 0)
	loadCsv("../ct/0525_all_not_found.csv", &notFound)

	notFound2 := make([]*EamLine, 0)

	log.Printf("len: %v\n", len(notFound))
	for _, line := range notFound {
		if f, ok := allMap[line.AssetTag]; ok {
			f.Type = line.Model
			notFound2 = append(notFound2, f)
		} else {
			log.Printf("--> %#v", line)
		}
	}

	if len(notFound2) > 0 {
		exportCsv("0525_all_not_found_withname.csv", &notFound2)
	}

}
