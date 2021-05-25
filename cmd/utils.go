package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

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

func AddUtf8Bom(inFile string, outFile string) error {
	source, err := os.OpenFile(inFile, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	_, err = fixInCsvUtf8(source)
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