package cmethod

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"os"
)

//gzip压缩
func GzipBytes(byt []byte) []byte {
	var buf bytes.Buffer
	//zw := gzip.NewWriter(&buf)
	zw, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)

	zw.Write(byt)
	if err := zw.Close(); err != nil {
	}
	return buf.Bytes()
}

//gzip解压缩
func UnGzipBytes(byt []byte) []byte {
	var buf bytes.Buffer
	buf.Write(byt)
	zr, _ := gzip.NewReader(&buf)
	defer func() {
		if zr != nil {
			zr.Close()
		}
	}()
	a, _ := ioutil.ReadAll(zr)
	return a
}

func EncDec(byt []byte) []byte {
	for i, v := range byt {
		byt[i] = v ^ byte(i*5+74)
	}
	return byt
}

func EncDec2(byt []byte) []byte {
	for i, v := range byt {
		byt[i] = (byte(i+75) & (^v)) | (v & (^byte(i + 75)))
	}
	return byt
}
func WriteFile(byt []byte, path, name string) error {
	err := CreateFile(path)
	if err != nil {
		return err
	}

	tmpFileName := path + name + ".tmp"
	f, err := os.Create(tmpFileName)
	if err != nil {
		return err
	}
	_, err = f.Write(byt)
	f.Close()
	if err != nil {
		return err
	}
	return os.Rename(tmpFileName, path+name)
}

//调用os.MkdirAll递归创建文件夹
func CreateFile(filePath string) error {
	if !IsExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		return err
	}
	return nil
}

// 判断所给路径文件/文件夹是否存在(返回true是存在)
func IsExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

//获取md5值
func byt2md5(byt []byte) string {
	h := md5.New()
	h.Write(byt)
	return hex.EncodeToString(h.Sum(nil))
}
