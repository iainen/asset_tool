package main

import (
	"github.com/gocarina/gocsv"
	"io"
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

func fixInCsvUtf8(csv *os.File) (bool, error) {
	// support utf-8 bom! `ef bb bf`
	buf := make([]byte, 3)
	l, err := csv.Read(buf)
	if err != nil || l != len(buf) {
		return false, err
	}

	if buf[0] == 0xEF && buf[1] == 0xBB && buf[2] == 0xBF {
		return true, nil
	} else {
		csv.Seek(0, io.SeekStart)
		return false, nil
	}
}

func fixOutCsvUtf8(csv *os.File) error {
	// support utf-8 bom! `ef bb bf`
	buf := []byte{0xEF, 0xBB, 0xBF}
	_, err := csv.Seek(0, io.SeekStart)
	_, err = csv.Write(buf)
	return err
}

func loadEamTransferCsv(csvPath string, filter string) ([]*EamTransferLine, []*EamTransferLine) {
	all := make([]*EamTransferLine, 0)
	inCsv, err := os.OpenFile(csvPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}

	_, err = fixInCsvUtf8(inCsv)
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