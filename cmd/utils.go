package main

import (
	"github.com/gocarina/gocsv"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func FixInCsvUtf8(csv *os.File) (bool, error) {
	// support utf-8 bom! `ef bb bf`
	buf := make([]byte, 3)
	l, err := csv.Read(buf)
	if err != nil || l != len(buf) {
		return false, err
	}

	if buf[0] == 0xEF && buf[1] == 0xBB && buf[2] == 0xBF {
		return true, nil
	} else {
		csv.Seek(0, io.SeekStart)
		return false, nil
	}
}

func FixOutCsvUtf8(csv *os.File) error {
	// support utf-8 bom! `ef bb bf`
	buf := []byte{0xEF, 0xBB, 0xBF}
	_, err := csv.Seek(0, io.SeekStart)
	_, err = csv.Write(buf)
	return err
}

func AddUtf8Bom(inFile string, outFile string) error {
	source, err := os.OpenFile(inFile, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	_, err = FixInCsvUtf8(source)
	if err != nil {
		return err
	}

	dest, err := os.OpenFile(outFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	utf8Bom := []byte{0xEF, 0xBB, 0xBF}
	if _, err := dest.Write(utf8Bom); err != nil {
		return err
	}

	buf := make([]byte, 64*1024)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := dest.Write(buf[:n]); err != nil {
			return err
		}
	}

	return nil
}

//获取指定目录下的所有文件,包含子目录下的文件
func GetAllFiles(dirPth string, filter string) (files []string, err error) {
	files = make([]string, 0)

	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	for _, fi := range dir {
		subPath := filepath.Join(dirPth, fi.Name())
		if fi.IsDir() { // 目录, 递归遍历
			subFiles, _ := GetAllFiles(subPath, filter)
			for _, f := range subFiles {
				files = append(files, f)
			}
		} else {
			ok := strings.HasPrefix(fi.Name(), filter)
			if ok {
				files = append(files, subPath)
			}
		}
	}

	return files, nil
}

func loadCsv(csvPath string, out interface{}) {
	inCsv, err := os.OpenFile(csvPath, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer inCsv.Close()

	_, err = FixInCsvUtf8(inCsv)
	if err != nil {
		panic(err)
	}

	if err := gocsv.UnmarshalFile(inCsv, out); err != nil {
		panic(err)
	}
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

	if err := FixOutCsvUtf8(outCsv); err != nil {
		panic(err)
	}

	err = gocsv.MarshalFile(out, outCsv)
	if err != nil {
		panic(err)
	}
}