package main

import (
	"io/ioutil"
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
