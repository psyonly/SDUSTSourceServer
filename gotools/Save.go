package gotools

// 本包用于实现一些底层对文件存取的操作

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// 文件目录种类
var species = map[string]string{
	"headImg": "/HeadImage",
	"root":    "",
}

type FilePath struct {
	FileType byte   // 表示文件种类默认0为目录，其他为文件
	FileAddr string // 目录或文件名
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

// 遍历root目录下所有文件，将结果放在一个切片里返回
func WalkThroughDir(root string, pfL int) []FilePath {
	fmt.Println(root)
	fp := make([]FilePath, 0)
	walk := func(path string, info os.FileInfo, err error) error {
		//fmt.Println(path, path[pfL:], pfL)
		if info.IsDir() {
			fp = append(fp, FilePath{
				FileType: 0,
				FileAddr: path[pfL:],
			})
		} else {
			fp = append(fp, FilePath{
				FileType: 1,
				FileAddr: path[pfL:],
			})
		}
		return nil
	}

	err := filepath.Walk("./"+root, walk)
	if err != nil {
		log.Fatal(err)
	}
	return fp
}
