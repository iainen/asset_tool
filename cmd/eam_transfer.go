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
	defer inCsv.Close()
	if err := gocsv.UnmarshalFile(inCsv, &all); err != nil {
		panic(err)
	}

	matchList := make([]*EamTransferLine, 0)
	if filter != "" {
		for _, line := range all {
			if strings.HasPrefix(line.AssetTag, filter) {
				log.Printf("%#v", line)
				matchList = append(matchList, line)
			}
		}
	}

	return all, matchList
}