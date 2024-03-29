package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Exists 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	// os.Stat获取文件信息
	if _, err := os.Stat(path); err != nil {
		return os.IsExist(err)
	}
	return true
}

// IsHiddenDirectory 路径是否是隐藏路径
func IsHiddenDirectory(path string) bool {
	return len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".")
}

// SubDir 输出所有子目录，目录名
func SubDir(folder string) ([]string, error) {
	subs, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	var ret []string
	for _, sub := range subs {
		if sub.IsDir() {
			ret = append(ret, sub.Name())
		}
	}
	return ret, nil
}
