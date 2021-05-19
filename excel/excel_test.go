package excel

import (
	"github.com/gocarina/gocsv"
	"log"
	"os"
	"testing"
)

type CtLine struct {
	Company string `csv:"公司"`
	AssetTag string `csv:"资产标签"`
	Model string `csv:"型号"`
	Brand string `csv:"品牌"`
	Manufacturer string `csv:"生产厂家"`
	Category string `csv:"类别："`
	Status string `csv:"状态"`
	Location string `csv:"位置"`
}

type Line struct {
	Company string
	AssetTag string `excel:"Asset Tag" csv:"Asset Tag"`
	Serial string `excel:"Serial Number" csv:"Serial Number"`
	Model string `excel:"Model name" csv:"Model name"`
	Brand string `excel:"Model Number" csv:"Model Number"`
	Manufacturer string
	Category string
	Status string
	Location string
}

func TestName(t *testing.T) {
	lines := make([]*Line, 0)
	Load("../ct/custom-assets-report-2021-05-18-051041.xlsx", &lines)
	for _, line := range lines {
		log.Printf("%#v", line)
	}
}

type CheckLine struct {
	AssetTag string `excel:"资产标签"`
	Status string `excel:"状态"`
	Location string `excel:"默认位置"`
}

func TestName2(t *testing.T) {
	lines := make([]*CheckLine, 0)
	Load("../ct/G10.xlsx", &lines)
	for _, line := range lines {
		log.Printf("%#v", line)
	}
}

func TestAll(t *testing.T) {
	allMap := make(map[string]*Line, 0)
	all := make([]*Line, 0)

	inCsv, err := os.OpenFile("../ct/custom-assets-report-2021-05-19-112946.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer inCsv.Close()
	if err := gocsv.UnmarshalFile(inCsv, &all); err != nil {
		panic(err)
	}
	for _, line := range all {
		allMap[line.AssetTag] = line
		log.Printf("-->: %#v", line)
	}

	export := make([]*Line, 0)

	toChecks := make([]*CheckLine, 0)
	Load("../ct/G10.xlsx", &toChecks)
	for _, line := range toChecks {
		if f, ok := allMap[line.AssetTag]; ok {
			f.Location = line.Location
			f.Status = line.Status

			// log.Printf("found: %#v", f)
			export = append(export, f)
		} else {
			log.Printf("not found: %#v", line)
		}
	}

	outCsv, err := os.OpenFile("G10.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer outCsv.Close()

	//csvContent, err := gocsv.MarshalString(&export)
	//fmt.Println(csvContent) // Display all clients as CSV string
	err = gocsv.MarshalFile(&export, outCsv)
	if err != nil {
		panic(err)
	}
}
