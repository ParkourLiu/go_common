package a_package_file

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// 组装二进制包
// ocrJsonStr: OCR识别结果json字符串 ["啊啊啊啊","大苏打实打实打算"]
// imgBytes: 图片原始byte数组
// return: 组装完成的二进制流
func PackOCRData(jsonBytes []byte, imgBytes []byte) []byte {
	jsonLen := uint32(len(jsonBytes))

	// 预分配buffer，减少内存扩容
	totalSize := 4 + len(jsonBytes) + len(imgBytes)
	buf := bytes.NewBuffer(make([]byte, 0, totalSize))

	// 1.写入4字节大端长度
	_ = binary.Write(buf, binary.BigEndian, jsonLen)
	// 2.写入json内容
	_, _ = buf.Write(jsonBytes)
	// 3.写入图片字节
	_, _ = buf.Write(imgBytes)

	return buf.Bytes()
}

// UnpackOCRData 解析二进制包
// data: PackOCRData输出的完整二进制
// return: ocrJson字符串, imageBytes, error
func UnpackOCRData(data []byte) (ocrJson string, imgBytes []byte, err error) {
	if len(data) < 4 {
		return "", nil, errors.New("data length less than 4 bytes header")
	}
	// 读取前4字节 json长度
	jsonLen := binary.BigEndian.Uint32(data[:4])
	endPos := 4 + int(jsonLen)
	if endPos > len(data) {
		return "", nil, errors.New("json length exceed total data size")
	}

	jsonBytes := data[4:endPos]
	imgBytes = data[endPos:]
	return string(jsonBytes), imgBytes, nil
}
