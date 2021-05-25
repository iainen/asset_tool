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

func TestCheckFile(t *testing.T) {
	toChecks := make([]*CheckLine, 0)
	excel.Load("../ct/bug_tocheck.xlsx", &toChecks)
	for _, line := range toChecks {
		log.Printf("%#v", line)
	}
}

func TestCheck2File(t *testing.T) {
	toChecks := make([]*CheckLine, 0)
	excel.Load("../ct/bug_tocheck2.xlsx", &toChecks)
	for _, line := range toChecks {
		log.Printf("%#v", line)
	}
}

func TestTrim(t *testing.T) {
	//fmt.Printf("[%q]\n", strings.Trim(" !!! Achtung! Achtung! !!! ", "!"))
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