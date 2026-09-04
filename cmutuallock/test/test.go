package main

import (
	"fmt"
	"hash/crc32"
	"os"
	"strings"
)

func main() {
	//a := os.Args[0]
	fmt.Println(os.Args[0])
	fmt.Println(os.Args[0][strings.LastIndex(os.Args[0], "/")+1:])
	fmt.Println(int(crc32.ChecksumIEEE([]byte(os.Args[0][strings.LastIndex(os.Args[0], "/")+1:]))))
}

//字节数(大端)组转成int(无符号的)
func String(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}
