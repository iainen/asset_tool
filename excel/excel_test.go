package excel

import (
	"log"
	"testing"
)

type Line struct {
	Company string
	AssetTag string `excel:"Asset Tag"`
	Serial string `excel:"Serial Number"`
	Model string `excel:"Model name"`
	Brand string `excel:"Model Number"`
	Manufacturer string
	Category string
	Status string
	Location string
}

func TestName(t *testing.T) {
	lines := make([]Line, 0)
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
	lines := make([]CheckLine, 0)
	Load("../ct/G10.xlsx", &lines)
	for _, line := range lines {
		log.Printf("%#v", line)
	}
}

func TestAll(t *testing.T) {
	allMap := make(map[string]Line, 0)
	all := make([]Line, 0)
	Load("../ct/custom-assets-report-2021-05-18-051041.xlsx", &all)
	for _, line := range all {
		allMap[line.AssetTag] = line
	}

	toChecks := make([]CheckLine, 0)
	Load("../ct/G10.xlsx", &toChecks)
	for _, line := range toChecks {
		if f, ok := allMap[line.AssetTag]; ok {
			f.Location = line.Location
			f.Status = line.Status

			log.Printf("found: %#v", f)
		} else {
			log.Printf("not found: %#v", line)
		}
	}
}
