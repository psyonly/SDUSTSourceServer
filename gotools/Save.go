package gotools

// 本包用于实现一些底层对文件存取的操作

import (
	"os"
	"path/filepath"
)

// 文件目录种类
var species = map[string]string{
	"headImg": "/HeadImage",
	"root":    "",
}

// 获得指定地址的目录大小
func GetDirSize(path string) (size int64, err error) {
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return
}

// 创建'spec' 类型的文件夹在 'dst' 目录下
func CreateDir(dst string, spec string) {
	if err := os.Mkdir(dst+species[spec], os.ModePerm); err != nil {
		panic(err)
	}
}
