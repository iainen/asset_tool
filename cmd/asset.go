package main

import (
	"flag"
	"git.code.oa.com/zhongkaizhu/assets_manager/excel"
	"github.com/gocarina/gocsv"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var inCtFile = flag.String("in-ct-all", "", "加载解析ct导出的csv总表")
var inCheckFile = flag.String("in-check", "", "加载解析盘点生成的xlsx表")
var eamFile = flag.String("eam", "", "加载epo导出的csv表")
var ebDir = flag.String("eb", "", "加载eb资产列表目录下所有asset.csv")

// 资产管理系统上导出的总资产，csv格式
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

// 线下盘点生成的资产，xlsx格式
type CheckLine struct {
	AssetTag string `excel:"资产标签" csv:"资产标签"`
	Status string `excel:"状态" csv:"状态"`
	Location string `excel:"默认位置" csv:"默认位置"`
}

// 用于导入资产管理系统的表，csv格式
type Line struct {
	Company string
	AssetTag string `excel:"Asset Tag" csv:"Asset Tag"`
	//	Serial string `excel:"Serial Number" csv:"Serial Number"`
	Model string `excel:"Model name" csv:"Model name"`
	Brand string `excel:"Model Number" csv:"Model Number"`
	Manufacturer string
	Category string
	Status string
	Location string
}

// 用于导入http://eam.oa.com/上导出的个人设备
type EamLine struct {
	AssetTag string `csv:"资产编码"`
	Name string `csv:"规格型号"`
	Type string `csv:"资产名称"`
	Brand string `csv:"品牌名称"`
}

func loadEamCsv(csvPath string, filter string) ([]*EamLine, []*EamLine) {
	all := make([]*EamLine, 0)
	inCsv, err := os.OpenFile(csvPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer inCsv.Close()

	_, err = fixInCsvUtf8(inCsv)
	if err != nil {
		panic(err)
	}

	if err := gocsv.UnmarshalFile(inCsv, &all); err != nil {
		panic(err)
	}

	matchList := make([]*EamLine, 0)
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

func exportCsv(outCsvPath string, out interface{}) {
	outCsv, err := os.OpenFile(outCsvPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer outCsv.Close()
	if _, err := outCsv.Seek(0, 0); err != nil { // Go to the start of the file
		panic(err)
	}

	if err := fixOutCsvUtf8(outCsv); err != nil {
		panic(err)
	}

	err = gocsv.MarshalFile(out, outCsv)
	if err != nil {
		panic(err)
	}
}

func exportCheckCsv(inCtAllCsvPath string, inCheckXlsxPath string, outCheckCsvPath string, out2 string) {
	all := make([]*CtLine, 0)
	allMap := make(map[string]*Line, 0)

	inCsv, err := os.OpenFile(inCtAllCsvPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer inCsv.Close()

	_, err = fixInCsvUtf8(inCsv)
	if err != nil {
		panic(err)
	}

	if err := gocsv.UnmarshalFile(inCsv, &all); err != nil {
		panic(err)
	}
	for _, line := range all {
		allMap[line.AssetTag] = &Line{
			Company: line.Company,
			AssetTag: line.AssetTag,
			Model: line.Model,
			Brand: line.Brand,
			Manufacturer: line.Manufacturer,
			Category: line.Category,
			Status: line.Status,
			Location: line.Location,
		}
		//log.Printf("-->: %#v", line)
	}

	founded := make([]*Line, 0)
	notFounded := make([]*CheckLine, 0)

	toChecks := make([]*CheckLine, 0)
	excel.Load(inCheckXlsxPath, &toChecks)
	for _, line := range toChecks {
		if f, ok := allMap[line.AssetTag]; ok {
			f.Location = line.Location
			f.Status = line.Status

			//log.Printf("found: %#v", f)
			founded = append(founded, f)
		} else {
			notFounded = append(notFounded, line)
			log.Printf("not found: %#v", line)
		}
	}

	exportCsv(outCheckCsvPath, &founded)
	if len(notFounded) > 0 {
		exportCsv(out2, &notFounded)
	}
}

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	if *inCtFile != "" && *inCheckFile != "" {
		outDir, outName := filepath.Split(*inCheckFile)
		i := strings.LastIndex(outName, filepath.Ext(*inCheckFile))
		out := filepath.Join(outDir, "checked_" + outName[0:i] + ".csv")
		out2:= filepath.Join(outDir, "unknown_" + outName[0:i] + ".csv")
		exportCheckCsv(*inCtFile, *inCheckFile, out, out2)
	}

	if *eamFile != "" {
		loadEamCsv(*eamFile, "TKMB")
	}

	if *ebDir != ""  {
		mergeEbAssets(*ebDir, filepath.Join(*ebDir, "all-asset.csv"))
	}
}
