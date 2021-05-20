package main

import (
	"log"
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

func TestGetAllFiles(t *testing.T) {
	mergeEbAssets("../ct/eb/", "../ct/eb/all-asset.csv")
}