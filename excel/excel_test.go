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
	Load("/Users/zhongkai/workplace/CloudTesting/AssetTool/ct/custom-assets-report-2021-05-18-051041.xlsx", lines)
	for _, line := range lines {
		log.Printf("line:%v", line)
	}
}
