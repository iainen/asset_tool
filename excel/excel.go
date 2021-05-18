package excel

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"log"
	"reflect"
	"strconv"
)

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

func populate(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)

	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)

	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(b)

	default:
		return fmt.Errorf("unsupported kind %s", v.Type())
	}
	return nil
}

func Load(input string, ptr interface{}) error {
	names := make([]string, 0)
	fields := make(map[string]int)

	rv := reflect.ValueOf(ptr)
	log.Printf("kind:%v, type:%v", rv.Kind(), rv.Type())

	if rv.Kind() != reflect.Slice {
		log.Fatalf("ptr must be Slice")
	}

	v := rv.Type().Elem() // the struct variable
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Field(i) // a reflect.StructField
		tag := fieldInfo.Tag           // a reflect.StructTag
		name := tag.Get("excel")
		if name == "" {
			name = fieldInfo.Name
		}

		names = append(names, name)
		fields[name] = i
		//log.Printf("name:%v", name)
	}

	f, err := excelize.OpenFile(input)
	if err != nil {
		log.Fatalf("fail to open file: %v", err)
	}

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)

	titleRow := rows[0]
	titleMap := detainTitles(names, titleRow)
	for _, row := range rows[1:] {
		line := reflect.New(rv.Type().Elem()).Elem()
		for _, name := range names {
			if _, ok := titleMap[name]; !ok {
				continue
			}
			_ = populate(line.Field(fields[name]), row[titleMap[name]])
		}
		log.Printf("line:%v", line)
		//rv.Set(reflect.Append(rv, line))
	}
	return nil
}