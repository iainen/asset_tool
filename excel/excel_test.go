package excel

import (
	"log"
	"testing"
)

type CheckLine struct {
	AssetTag string `excel:"资产标签"`
	Status string `excel:"状态"`
	Location string `excel:"默认位置"`
}

func TestCheckFile(t *testing.T) {
	toChecks := make([]*CheckLine, 0)
	Load("../resources/bug_tocheck.xlsx", &toChecks)
	for _, line := range toChecks {
		log.Printf("%#v", line)
	}
}

func TestCheck2File(t *testing.T) {
	toChecks := make([]*CheckLine, 0)
	Load("../resources/bug_tocheck2.xlsx", &toChecks)
	for _, line := range toChecks {
		log.Printf("%#v", line)
	}
}
