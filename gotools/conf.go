package gotools

import (
	"fmt"
	"io/ioutil"
)

// 本包包含若干与配置文件/解释文件相关的功能实现

func GetVersionScript(version string) (ver string) {
	dir := ""
	//fmt.Println(version)
	dir = fmt.Sprintf("./Files/versions/%s.html", version)
	//fmt.Println(dir)
	content, err := ioutil.ReadFile(dir)
	if err != nil {
		fmt.Println("Read file err", err)
		return
	}
	ver = string(content[:])
	//fmt.Println(ver)
	return
}
