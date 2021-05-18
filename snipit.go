package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"log"
	"strings"
)

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
	ROM               = "存储"
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

func exportSnipITExcel(outfile string, assets map[string]*WetestAsset) {
	f := excelize.NewFile()
	sheet := f.GetSheetName(0)

	titles := []string{
		COMPANY,
		ASSET_TAG,
		SERIAL,
		MODEL_NAME,
		MODEL_NUMBER,
		CATEGORY,
		STATUS,
		Manufacturer,
		EPO_NAME,
		IMEI,
	}

	for i, v := range titles {
		f.SetCellValue(sheet, fmt.Sprintf("%c%d", int('A')+i, 1), v)
	}

	line := 2
	c := []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J'}
	for tag, item := range assets {
		category := "Android"
		if strings.Contains(item.Model, "iPhone") || strings.Contains(item.Model, "iPad") {
			category = "iOS"
		}

		i := 0
		f.SetCellValue(sheet, fmt.Sprintf("%c%d", c[i], line), "腾讯"); i++
		f.SetCellValue(sheet, fmt.Sprintf("%c%d", c[i], line), tag); i++
		f.SetCellValue(sheet, fmt.Sprintf("%c%d", c[i], line), item.Serial); i++
		f.SetCellValue(sheet, fmt.Sprintf("%c%d", c[i], line), item.Model); i++
		f.SetCellValue(sheet, fmt.Sprintf("%c%d", c[i], line), item.Brand); i++ //model_number
		f.SetCellValue(sheet, fmt.Sprintf("%c%d", c[i], line), category); i++
		f.SetCellValue(sheet, fmt.Sprintf("%c%d", c[i], line), "新增"); i++
		f.SetCellValue(sheet, fmt.Sprintf("%c%d", c[i], line), item.Manu); i++
		f.SetCellValue(sheet, fmt.Sprintf("%c%d", c[i], line), item.FullName); i++
		f.SetCellValue(sheet, fmt.Sprintf("%c%d", c[i], line), item.IMEI); i++
		line++
	}

	// Save spreadsheet by the given path.
	if err := f.SaveAs(outfile); err != nil {
		fmt.Println(err)
	}
}