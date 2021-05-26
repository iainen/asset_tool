package main

import (
	"github.com/gocarina/gocsv"
	"log"
	"os"
	"strings"
)

// 预约发货后，epo提供的资产表
type EamTransferLine struct {
	AssetTag string `csv:"固资码"`
	Name string `csv:"名称"`
	Type string `csv:"物料码"`
}

func loadEamTransferCsv(csvPath string, filter string) ([]*EamTransferLine, []*EamTransferLine) {
	all := make([]*EamTransferLine, 0)
	inCsv, err := os.OpenFile(csvPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}

	_, err = FixInCsvUtf8(inCsv)
	if err != nil {
		panic(err)
	}

	defer inCsv.Close()
	if err := gocsv.UnmarshalFile(inCsv, &all); err != nil {
		panic(err)
	}

	matchList := make([]*EamTransferLine, 0)
	if filter != "" {
		for _, line := range all {
			if strings.HasPrefix(line.AssetTag, filter) {
				matchList = append(matchList, line)
			}
		}
	}

	return all, matchList
}

func mergeAssets(dir string, fileName string, assetPrefix string) []*EamTransferLine {
	all := make([]*EamTransferLine, 0)

	csvList, _ := GetAllFiles(dir, fileName)
	for _, csv := range csvList {
		log.Printf("--> %#v", csv)
		_, list2 := loadEamTransferCsv(csv, assetPrefix)
		for _, line := range list2 {
			all = append(all, line)
		}
	}

	return all
}