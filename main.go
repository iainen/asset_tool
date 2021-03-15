package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "github.com/360EntSecGroup-Skylar/excelize/v2"
    "log"
    "os"
)

type Data struct {
    Name                   string `json:"name"`
    Model                  string `json:"model"`
    Price                  string `json:"price"`
    Date                   string `json:"date"`
    Manu                   string `json:"manu"`
    Cpuname                string `json:"cpuname"`
    Corename               string `json:"corename"`
    Rate                   string `json:"rate"`
    Num                    string `json:"num"`
    NumCores               string `json:"num_cores"`
    CPUScore               string `json:"cpu_score"`
    GpuScore               string `json:"gpu_score"`
    BenchmarkScore         string `json:"benchmark_score"`
    RAM                    string `json:"ram"`
    GpuName                string `json:"gpu_name"`
    GpuMarks               string `json:"gpu_marks"`
    UserIndex              string `json:"user_index"`
    BenchmarkIndexIntop500 string `json:"benchmark_index_inTop500"`
    MainTop                string `json:"main_top"`
    CPUFrequency           string `json:"cpu_frequency"`
    Zram                   string `json:"zram"`
    Sp                     string `json:"sp"`
    ScreenRatio            string `json:"screen_ratio"`
    Vulkan                 string `json:"vulkan"`
    OpenglEs               string `json:"opengl_es"`
    Os                     string `json:"os"`
    Notch                  string `json:"notch"`
}

type MobileDeviceBenchmark struct {
    Status int `json:"status"`
    Data   []Data `json:"data"`
    /*
    Data   []struct {
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
    } `json:"data"`
     */
}

var benchmark = flag.String("benchmark", "", "benchmark json file")

var modelMap = make(map[string][]Data)

func exportExcel(outfile string) {
    f := excelize.NewFile()
    // Create a new sheet.
    index := f.NewSheet("Sheet")
    // Set active sheet of the workbook.
    f.SetActiveSheet(index)


    titles := []string {
        "型号",
        "手机名",
        "厂商",
        "CPU名",
        "异性屏",
        "分辨率",
        "上市日期",
    }

    for i,v := range titles {
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
       line ++
    }

    // Save spreadsheet by the given path.
    if err := f.SaveAs(outfile); err != nil {
      fmt.Println(err)
    }
}

func main() {
    flag.Parse()
    log.SetFlags(log.Lshortfile | log.LstdFlags)



    jsonBytes, err := os.ReadFile(*benchmark)
    if err != nil {
        log.Printf("fail to open logfile: %v", err)
        return
    }

    //var mobile MobileDescription
    //err = json.Unmarshal(jsonBytes, &mobile)
    //if err != nil {
    //    fmt.Println(err.Error())
    //}

    var mobileBenchmark MobileDeviceBenchmark
    err = json.Unmarshal(jsonBytes, &mobileBenchmark)
    if err == nil {
       fmt.Printf("fail to parse json: %v", err)
       return
    }


    log.Printf("length of array: %v", len(mobileBenchmark.Data))
    for _, item := range mobileBenchmark.Data {
        //log.Printf("%v: model:%v，name:%v\n", index, item.Model, item.Name)
        if names, ok := modelMap[item.Model]; !ok {
            names = make([]Data, 0)
            modelMap[item.Model] = append(names, item)
        } else {
            modelMap[item.Model] = append(names, item)
        }
    }

    //index := 0
    //for model, names := range modelMap {
    //    index++
    //    //log.Printf("%v: model:%v, data:%v", index, model, names)
    //
    //    name := names[0]
    //    log.Printf("%v: model:%v，name:%v, manu:%v, cpu:%v, Notch:%v, screenRation:%v, date:%v\n", index,
    //       model,
    //       name.Name,
    //       name.Manu,
    //       name.Cpuname,
    //       name.Notch,
    //       name.ScreenRatio,
    //       name.Date)
    //}

    exportExcel("out.xlsx")
}
