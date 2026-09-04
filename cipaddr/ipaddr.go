package cipaddr

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sort"
	"strings"
)

var ipInfos []iPInfoStruct

func init() {
	err := loadIPDatabase("./ipv4.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
		return
	}
}

// IPInfoStruct 存储IP段的信息
type iPInfoStruct struct {
	StartIP  uint32 // 起始IP的整数形式
	EndIP    uint32 // 结束IP的整数形式
	Country  string
	Province string
	City     string
}

// ipToUint32 将点分十进制的IP转换为uint32整数
func IpToUint32(ipStr string) (uint32, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return 0, fmt.Errorf("invalid IP address: %s", ipStr)
	}
	ip = ip.To4() // 确保是IPv4
	if ip == nil {
		return 0, fmt.Errorf("not an IPv4 address: %s", ipStr)
	}
	var result uint32
	for i := 0; i < 4; i++ {
		result = result<<8 + uint32(ip[i])
	}
	return result, nil
}
func Uint32ToIP(ipUint uint32) string {
	// 拆分uint32为4个8位字节（对应IPv4的4个段）
	// 注意字节顺序：高位在前（大端序），与IP地址的书写顺序一致
	byte1 := byte(ipUint >> 24) // 最高8位（第一个段）
	byte2 := byte(ipUint >> 16) // 次高8位（第二个段）
	byte3 := byte(ipUint >> 8)  // 次低8位（第三个段）
	byte4 := byte(ipUint)       // 最低8位（第四个段）

	// 拼接成IPv4地址字节数组
	ipBytes := []byte{byte1, byte2, byte3, byte4}
	// 转换为net.IP类型并返回字符串
	return net.IPv4(ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3]).String()
}

// LoadIPDatabase 从文件加载IP数据库并构建索引
func loadIPDatabase(filePath string) (err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) < 5 {
			continue // 跳过格式不正确的行
		}

		startIP, err := IpToUint32(parts[0])
		if err != nil {
			continue
		}
		endIP, err := IpToUint32(parts[1])
		if err != nil {
			continue
		}

		ipInfo := iPInfoStruct{
			StartIP:  startIP,
			EndIP:    endIP,
			Country:  parts[2],
			Province: parts[3],
			City:     parts[4],
			//ISP:      parts[5],
			//Region:   parts[6],
		}
		ipInfos = append(ipInfos, ipInfo)
	}

	// 按起始IP排序，以便后续二分查找
	sort.Slice(ipInfos, func(i, j int) bool {
		return ipInfos[i].StartIP < ipInfos[j].StartIP
	})

	return
}

// QueryIP 查询IP所属的信息
func Ip2Addr(ipStr string) (ipInfo iPInfoStruct) {
	ip, err := IpToUint32(ipStr)
	if err != nil {
		return
	}

	// 二分查找
	low, high := 0, len(ipInfos)-1
	for low <= high {
		mid := (low + high) / 2
		midInfo := ipInfos[mid]

		if ip < midInfo.StartIP {
			high = mid - 1
		} else if ip > midInfo.EndIP {
			low = mid + 1
		} else {
			// 找到匹配的IP段
			return midInfo
		}
	}

	return
}
