package main

import (
	"encoding/json"
	"log"
	"os"
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

func loadBenchmark2ModelMap(inputfile string) (map[string]Data, error) {
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

	benchMap := make(map[string]Data)
	//log.Printf("model number: %v", len(benchMap))
	index := 0
	for _, item := range mobileBenchmark.Data {
		index ++
		benchMap[item.Model] = item
		log.Printf("%v: [%v] --> [%v]", index, item.Model, item.Name)
	}

	return benchMap, nil
}

func loadBenchmark2NameMap(inputfile string) (map[string][]Data, error) {
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