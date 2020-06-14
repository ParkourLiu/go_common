package main

import (
	"fmt"
	"strings"
)

func main() {
	path := strings.TrimRight("/aaa/bbb/ccc/ddd/", "/") //去除后/符号
	lastIndex := strings.LastIndex(path, "/")           //判断倒数第一个文件夹分隔符位置
	path = path[:lastIndex+1]
	fmt.Println(path)
}
func f(i int) {
	fmt.Println(i)
}
