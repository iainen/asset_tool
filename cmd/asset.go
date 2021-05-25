package main

import (
	"git.code.oa.com/zhongkaizhu/assets_manager/excel"
	"github.com/gocarina/gocsv"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
	"strings"
)

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

// 由http://eam.oa.com/上导出的个人设备的csv表
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
				// log.Printf("%#v", line)
				matchList = append(matchList, line)
			}
		}
	}

	return all, matchList
}

func exportCsv(outCsvPath string, out interface{}) {
	outCsv, err := os.OpenFile(outCsvPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
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

func exportCheckCsv(snipeItCsvPath string, inCheckXlsxPath string, outCheckCsvPath string, out2 string) {
	all := make([]*CtLine, 0)
	allMap := make(map[string]*Line, 0)

	inCsv, err := os.OpenFile(snipeItCsvPath, os.O_RDONLY, os.ModePerm)
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

func diffCsv(epoCsv string, snipeItCsv string) {
	importCsv := func(csvPath string, out interface{}) {
		inCsv, err := os.OpenFile(csvPath, os.O_RDONLY, os.ModePerm)
		if err != nil {
			panic(err)
		}
		defer inCsv.Close()

		_, err = fixInCsvUtf8(inCsv)
		if err != nil {
			panic(err)
		}

		if err := gocsv.UnmarshalFile(inCsv, out); err != nil {
			panic(err)
		}
	}

	snipeAll := make([]*CtLine, 0)
	importCsv(snipeItCsv, &snipeAll)

	innerMap := make(map[string]*Line, 0)
	for _, line := range snipeAll {
		innerMap[line.AssetTag] = &Line{
			Company:      line.Company,
			AssetTag:     line.AssetTag,
			Model:        line.Model,
			Brand:        line.Brand,
			Manufacturer: line.Manufacturer,
			Category:     line.Category,
			Status:       line.Status,
			Location:     line.Location,
		}
		//log.Printf("-->: %#v", line)
	}

	epoAll := make([]*EamLine, 0)
	importCsv(epoCsv, &epoAll)

	for _, line := range epoAll {
		if _, ok := innerMap[line.AssetTag]; !ok {
			log.Printf("not found: %#v", line)
		}
	}
}

func main() {
	//flag.Parse()
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	app := &cli.App{
		Commands: []*cli.Command {
			{
				Name:    "eam",
				Usage:   "处理eam.om.com上导出的文件",
				Action:  func(c *cli.Context) error {
					log.Printf("prefix:%v", c.String("prefix"))
					_, filterList := loadEamCsv(c.String("input"), c.String("prefix"))
					exportCsv(c.String("output"), filterList)

					return nil
				},

				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "prefix",
						Aliases: []string{"p"},
						Value:   "TKMB",
						Usage:   "filter with prefix, like `TKMB/TKMR` etc.",
					},

					&cli.StringFlag{
						Name:    "input",
						Aliases: []string{"i"},
						Required: true,
						Usage:   "csv file exported by eam.oa.com",
					},

					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Required: true,
						Usage:   "output file to save",
					},
				},
			},

			{
				Name:    "check",
				Usage:   "用于设备盘点，产生Snipe-IT资产管理系统导入所需的csv文件",
				Action: func(c *cli.Context) error {
					inFile := c.String("input")
					assetFile := c.String("all")

					outDir, outName := filepath.Split(inFile)
					i := strings.LastIndex(outName, filepath.Ext(inFile))
					out := filepath.Join(outDir, "checked_" + outName[0:i] + ".csv")
					out2:= filepath.Join(outDir, "unknown_" + outName[0:i] + ".csv")
					exportCheckCsv(assetFile, inFile, out, out2)
					return nil
				},

				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "input",
						Aliases:  []string{"i"},
						Required: true,
						Usage:    "csv文件，由机房管理员执行盘点时产生",
					},

					&cli.StringFlag{
						Name:     "all",
						Aliases:  []string{"A"},
						Required: true,
						Usage:    "csv文件, 从Snipe-IT资产管理系统中导出的总资产表",
					},
				},
			},

			{
				Name:  "merge",
				Usage: "合并多个文件，搜索子目录下所有指定文件(asset.csv)",
				Action: func(c *cli.Context) error {
					inDir := c.String("dir")
					outFile := c.String("out")
					all := mergeAssets(inDir, "asset.csv", c.String("prefix"))
					exportCsv(outFile, &all)
					return nil
				},

				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "dir",
						Aliases:  []string{"d"},
						Required: true,
						Usage:    "输入目录",
					},

					&cli.StringFlag{
						Name:     "out",
						Aliases:  []string{"o"},
						Required: true,
						Usage:    "输出文件",
					},

					&cli.StringFlag{
						Name:     "prefix",
						Aliases:  []string{"p"},
						Value: "TKMB",
						Usage:    "过滤指定的资产类型，如TKMB、TKNB等",
					},
				},
			},
			{
				Name:    "diff",
				Usage:   "用于设备盘点，产生Snipe-IT资产管理系统导入所需的csv文件",
				Action: func(c *cli.Context) error {
					eamCsv := c.String("eam")
					snipeItCsv := c.String("snipe-it")
					diffCsv(eamCsv, snipeItCsv)
					return nil
				},

				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "eam",
						Aliases:  []string{"e"},
						Required: true,
						Usage:    "csv文件，由eam过滤导出的个人资产设备列表",
					},

					&cli.StringFlag{
						Name:     "snipe-it",
						Aliases:  []string{"s"},
						Required: true,
						Usage:    "csv文件, 从Snipe-IT资产管理系统中导出的总资产表",
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
