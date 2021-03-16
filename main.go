package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type Data struct {
	Name                   string      `json:"name"`
	Model                  string      `json:"model"`
	Rom                    string      `json:"rom"`
	Price                  string      `json:"price"`
	Date                   string      `json:"date"`
	Manu                   string      `json:"manu"`
	Cpuname                string      `json:"cpuname"`
	Corename               string      `json:"corename"`
	Rate                   int         `json:"rate"`
	Num                    string      `json:"num"`
	NumCores               string      `json:"num_cores"`
	CPUScore               string      `json:"cpu_score"`
	GpuScore               string      `json:"gpu_score"`
	BenchmarkScore         string      `json:"benchmark_score"`
	RAM                    string      `json:"ram"`
	CPUFrequency           string      `json:"cpu_frequency"`
	GpuNAme                string      `json:"gpu_n  ame,omitempty"`
	GpuMarks               string      `json:"gpu_marks"`
	UserIndex              int         `json:"user_index"`
	BenchmarkIndexIntop500 int         `json:"benchmark_index_inTop500"`
	MainTop                string      `json:"main_top"`
	Top3                   string      `json:"top3"`
	Top4                   string      `json:"top4"`
	Top5                   string      `json:"top5"`
	Zram                   string      `json:"zram"`
	ScreenRatio            string      `json:"screen_ratio"`
	Sp                     string      `json:"sp"`
	OpenglEs               string      `json:"opengl_es"`
	Vulkan                 string      `json:"vulkan"`
	Os                     interface{} `json:"os"`
	Notch                  string      `json:"notch"`
	GpuName                string      `json:"gpu_name,omitempty"`
}

type MobileDeviceBenchmark struct {
	Status int    `json:"status"`
	Data   []Data `json:"data"`
}

func exportExcel(outfile string, modelMap map[string][]Data) {
	f := excelize.NewFile()
	// Create a new sheet.
	index := f.NewSheet("Sheet")
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)

	titles := []string{
		"型号",
		"手机名",
		"厂商",
		"CPU名",
		"异性屏",
		"分辨率",
		"上市日期",
	}

	for i, v := range titles {
		//log.Printf("%c%d:%v", int('A')+i, 1, v)
		f.SetCellValue("Sheet", fmt.Sprintf("%c%d", int('A')+i, 1), v)
	}

	line := 2
	for model, names := range modelMap {
		name := names[0]
		f.SetCellValue("Sheet", fmt.Sprintf("A%d", line), model)
		f.SetCellValue("Sheet", fmt.Sprintf("B%d", line), name.Name)
		f.SetCellValue("Sheet", fmt.Sprintf("C%d", line), name.Manu)
		f.SetCellValue("Sheet", fmt.Sprintf("D%d", line), name.Cpuname)
		f.SetCellValue("Sheet", fmt.Sprintf("E%d", line), name.Notch)
		f.SetCellValue("Sheet", fmt.Sprintf("F%d", line), name.ScreenRatio)
		f.SetCellValue("Sheet", fmt.Sprintf("G%d", line), name.Date)
		line++
	}

	// Save spreadsheet by the given path.
	if err := f.SaveAs(outfile); err != nil {
		fmt.Println(err)
	}
}

func loadFromBenchmark(inputfile string) (map[string][]Data, error) {
	jsonBytes, err := os.ReadFile(inputfile)
	if err != nil {
		log.Printf("fail to open logfile: %v", err)
		return nil, err
	}

	var mobileBenchmark MobileDeviceBenchmark
	err = json.Unmarshal(jsonBytes, &mobileBenchmark)
	//if err != nil {
	//	log.Printf("fail to parse json: %v", err)
	//	return nil, err
	//}

	benchMap := make(map[string][]Data)
	for _, item := range mobileBenchmark.Data {
		if datas, ok := benchMap[item.Name]; !ok {
			datas = make([]Data, 0)
			benchMap[item.Name] = append(datas, item)
		} else {
			benchMap[item.Name] = append(datas, item)
		}
	}

	return benchMap, nil
}

const (
	COMPANY           = "Company"
	ASSET_TAG         = "Asset Tag"
	SERIAL            = "Serial Number"
	MODEL_NAME        = "Model name"
	CATEGORY          = "Category"
	STATUS            = "Status"
	LOCATION          = "Location"
	Manufacturer      = "Manufacturer"
	CPU               = "CPU"
	NOTCH             = "异形屏"
	EPO_NAME          = "epo全名"
	IMEI              = "IMEI"
	USB_TYPE          = "USB类型"
	USB_LEVEL         = "USB级别"
	SCREEN_RESOLUTION = "分辨率"
	RAM               = "内存"
	ROM               = "存储大小"
	WIFI              = "WiFi"
	IS_5G             = "5G"
	MODEL_NUMBER      = "Model Number"
)

var AssetTitles = []string{
	COMPANY, ASSET_TAG, SERIAL, MODEL_NAME, CATEGORY, STATUS, LOCATION, Manufacturer, CPU, NOTCH,
	EPO_NAME, IMEI, USB_TYPE, USB_LEVEL, SCREEN_RESOLUTION, RAM, ROM, WIFI, IS_5G, MODEL_NUMBER,
}

func detainTitles(titles []string, titleRow []string) map[string]int {
	m := make(map[string]int)
	for _, t1 := range titles {
		for i2, t2 := range titleRow {
			if t1 == t2 {
				m[t1] = i2
				break
			}
		}
	}
	return m
}

type AssetMini struct {
	AssetTag string
	FullName string
}

type ModelMini struct {
	ModelName   string
	ModelNumber string
	Manu        string
	Catergory   string
}

type Model struct {
	Models []ModelMini
	Assets []AssetMini
}

func loadFrom4500(inputFile string) (map[string]*Model, error) {
	f, err := excelize.OpenFile(inputFile)
	if err != nil {
		log.Fatalf("fail to open file: %v", err)
	}

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)

	titleRow := rows[0]
	titleMap := detainTitles(AssetTitles, titleRow)
	log.Printf("len:%v, %v", len(titleRow), titleRow)
	log.Printf("%v", titleMap)

	var modelMap = make(map[string]*Model)
	for _, row := range rows[1:] {
		modelName := row[titleMap[MODEL_NAME]]
		manu := row[titleMap[Manufacturer]]

		var modelNumber string
		if len(row) > titleMap[MODEL_NUMBER] {
			modelNumber = row[titleMap[MODEL_NUMBER]]
		}

		catergory := row[titleMap[CATEGORY]]

		//log.Printf("name:%v, manu:%v, %v, %v", row[titleMap[MODEL_NAME]], manu, modelNumber, catergory)
		//continue

		items, ok := modelMap[modelName]
		if !ok {
			items = &Model{
				Models: make([]ModelMini, 0),
				Assets: make([]AssetMini, 0),
			}

			modelMap[modelName] = items
		}

		found := false
		for _, j := range items.Models {
			if j.ModelName == modelName && j.ModelNumber == modelNumber && j.Manu == manu && j.Catergory == catergory {
				found = true
				break
			}
		}
		if !found {
			items.Models = append(items.Models, ModelMini{
				ModelName:   modelName,
				ModelNumber: modelNumber,
				Manu:        manu,
				Catergory:   catergory,
			})
		}

		items.Assets = append(items.Assets, AssetMini{
			AssetTag: row[titleMap[ASSET_TAG]],
			FullName: row[titleMap[EPO_NAME]],
		})
	}

	return modelMap, nil
	// show
	/*
		index := 0
		for name, item := range modelMap {
			if len(item.Assets) > 1 {
				log.Printf("%v: model:%v", index, name)
				for _, asset := range item.Assets {
					log.Printf("  %v", asset.FullName)
				}
			}

			index++
		}
	*/
}

var benchmark = flag.String("benchmark", "", "benchmark json file")
var levy4500 = flag.String("levy4500", "", "levy 4500 xlsx file")
var out = flag.String("output", "", "export to a excel file")

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	var benchModelMap map[string][]Data
	var err error
	if *benchmark != "" {
		benchModelMap, err = loadFromBenchmark(*benchmark)
		if err != nil {
			return
		}
	}

	if *levy4500 != "" {
		maps, err := loadFrom4500(*levy4500)
		if err != nil {
			log.Printf("fail to load levy excel: %v", err)
			return
		}

		if benchModelMap != nil {
			nameList := make([]string, len(benchModelMap), len(benchModelMap))
			for name, _ := range benchModelMap {
				nameList = append(nameList, name)
			}

			sort.Strings(nameList)
			for _, item := range maps {
				for _, asset := range item.Assets {
					//log.Printf("  %v", asset.FullName)
					found := false
					for i := len(nameList)-1; i >= 0; i-- {
						if strings.Contains(strings.ToLower(asset.FullName), strings.ToLower(nameList[i])) {
							found = true
							log.Printf("%v --> %v", nameList[i], asset.FullName)
							break
						}
					}

					if !found {
						modelName := item.Models[0].ModelName
						log.Printf("X %v: %v: %v", asset.FullName, modelName, benchModelMap[modelName])
					}
				}
			}
		}
	}

	if *out != "" {
		exportExcel(*out, benchModelMap)
	}
}
