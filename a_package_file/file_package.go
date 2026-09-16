package a_package_file

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

// ArchiveJob 生产侧任务
type ArchiveJob struct {
	Path string
	Data []byte
}

type BufferArchiver struct {
	destPaths   []string //最终文件写入到哪里
	packageSize int      //包阈值大小
	jobChan     chan ArchiveJob
	buf         *bytes.Buffer
}

func NewBufferArchiver(destPaths []string, packageSize, poolSize int) *BufferArchiver {
	ba := &BufferArchiver{
		destPaths:   destPaths,
		packageSize: packageSize,
		jobChan:     make(chan ArchiveJob, poolSize),
		buf:         bytes.NewBuffer(nil),
	}
	go ba.consumerLoop()
	return ba
}

// AddFile 生产协程并发投递任务
func (ba *BufferArchiver) AddFile(path string, data []byte) {
	ba.jobChan <- ArchiveJob{
		Path: path,
		Data: data,
	}
}

func (ba *BufferArchiver) consumerLoop() {
	for job := range ba.jobChan {
		pathBytes := []byte(job.Path)
		pathLen := uint32(len(pathBytes))
		dataLen := uint64(len(job.Data))

		// 直接完整写入本条记录，不做预先预判
		_ = binary.Write(ba.buf, binary.LittleEndian, pathLen)
		_, _ = ba.buf.Write(pathBytes)
		_ = binary.Write(ba.buf, binary.LittleEndian, dataLen)
		_, _ = ba.buf.Write(job.Data)

		// 写完本条，判断是否达到阈值；消费协程内直接调用上层写盘回调
		if ba.buf.Len() < ba.packageSize {
			continue
		}
		nowDay := time.Now().Format("20060102/")
		nowUnixMilli := fmt.Sprint(time.Now().UnixMilli())
		for _, destPath := range ba.destPaths {
			destPath = destPath + nowDay
			filePath := destPath + nowUnixMilli
			doFlush(destPath, filePath, ba.buf.Bytes())
		}
		ba.buf.Reset()
	}
}

// doFlush 内部：拷贝buffer数据，调用上层回调，清空buffer
func doFlush(destPath, filePath string, fbs []byte) {
	//fbs = ZstdCompress(fbs) //zstd压缩
	err := os.MkdirAll(destPath, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = os.WriteFile(filePath+".tmp", fbs, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = os.Rename(filePath+".tmp", filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	//log.Error("写入", filePath)
}

//解包

// UnmarshalArchivePartialSilent 解析归档字节，遇到损坏直接停止，不返回error
// 返回已经成功解析的条目；无论截断、越界、损坏都只返回已读出数据，err永远nil
func UnmarshalArchivePartialSilent(raw []byte) []*ArchiveJob {
	var entries []*ArchiveJob
	offset := 0
	total := len(raw)

	for offset < total {
		// pathLen uint32
		if offset+4 > total {
			break
		}
		pathLen := binary.LittleEndian.Uint32(raw[offset : offset+4])
		offset += 4

		pathEnd := offset + int(pathLen)
		if pathEnd > total {
			break
		}
		path := string(raw[offset:pathEnd])
		offset = pathEnd

		// dataLen uint64
		if offset+8 > total {
			break
		}
		dataLen := binary.LittleEndian.Uint64(raw[offset : offset+8])
		offset += 8
		// ----------------核心修复----------------
		// 剩余字节转 uint64，无符号对比，避免溢出
		rem := uint64(total - offset)
		// 如果声明dataLen大于缓冲区剩余，判定损坏直接退出
		if dataLen > rem {
			break
		}
		dataEnd := offset + int(dataLen)
		if dataEnd > total {
			break
		}
		data := raw[offset:dataEnd]
		offset = dataEnd

		entries = append(entries, &ArchiveJob{
			Path: path,
			Data: data,
		})
	}
	return entries
}
