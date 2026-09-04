package cid

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type UnixId struct {
	divisor      int64
	lastUnixTime int64
	machineId    string
	l            *sync.Mutex
}

func NewUnixId(machineId ...string) (uid *UnixId) {
	uid = &UnixId{
		divisor:      0,
		lastUnixTime: time.Now().Unix(),
		l:            &sync.Mutex{},
	}
	if len(machineId) > 0 {
		uid.machineId = machineId[0]
	}
	return
}

func (u *UnixId) getId() (id, div int64) {
	u.l.Lock()
	defer u.l.Unlock()
	nowUnixTime := time.Now().Unix()
	if nowUnixTime > u.lastUnixTime {
		u.lastUnixTime = nowUnixTime
		u.divisor = 0
	}
	u.divisor++
	return u.lastUnixTime << 2, u.divisor
}

func (u *UnixId) NextId() (id string) {
	lastUnixTime, divisor := u.getId()
	return fmt.Sprintf("%d%d%s", lastUnixTime, divisor, u.machineId)
}
func (u *UnixId) NextId36() (id string) {
	lastUnixTime, divisor := u.getId()
	lastUnixTime36 := jinzhi36ToString(lastUnixTime)
	return fmt.Sprintf("%s%d%s", lastUnixTime36, divisor, u.machineId)
}

func (u *UnixId) NextId62() (id string) {
	lastUnixTime, divisor := u.getId()
	lastUnixTime62 := jinzhi62ToString(lastUnixTime)
	return fmt.Sprintf("%s%d%s", lastUnixTime62, divisor, u.machineId)
}

// randomLen 生成id的长度
func (u *UnixId) NextId36Random(randomLen int64) (id string) {
	lastUnixTime, divisor := u.getId()
	lastUnixTime36 := jinzhi36ToString(lastUnixTime)
	id = fmt.Sprintf("%s%d", lastUnixTime36, divisor)
	randomLen = randomLen - int64(len(id)) //获取长度差值
	randomStr := strings.Builder{}
	for i := int64(1); i <= randomLen; i++ {
		randomStr.WriteString(datas36[(lastUnixTime-i*divisor*randomLen)%36])
	}
	return fmt.Sprintf("%s%s%s", id, randomStr.String(), u.machineId)
}

// randomLen 生成id的长度
func (u *UnixId) NextId62Random(randomLen int64) (id string) {
	lastUnixTime, divisor := u.getId()
	lastUnixTime62 := jinzhi62ToString(lastUnixTime)
	id = fmt.Sprintf("%s%d", lastUnixTime62, divisor)
	randomLen = randomLen - int64(len(id)) //获取长度差值
	randomStr := strings.Builder{}
	for i := int64(1); i <= randomLen; i++ {
		randomStr.WriteString(datas62[(lastUnixTime-i*divisor*randomLen)%62])
	}
	return fmt.Sprintf("%s%s%s", id, randomStr.String(), u.machineId)
}

var datas62 = [62]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
var datas36 = [36]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

func jinzhi62ToString(num int64) (s string) {
	jinzhiNum := jinZhiZhuanHuan(num, 62)
	for _, i := range jinzhiNum {
		s += datas62[i]
	}
	return
}

func jinzhi36ToString(num int64) (s string) {
	jinzhiNum := jinZhiZhuanHuan(num, 36)
	for _, i := range jinzhiNum {
		s += datas36[i]
	}
	return
}

// 进制转换
func jinZhiZhuanHuan(num, jinzhi int64) (jinzhiNum []int64) {
	for num > 0 {
		yushu := num % jinzhi
		num = num / jinzhi
		jinzhiNum = append(jinzhiNum, yushu)
	}
	jinzhis := []int64{}
	for i := len(jinzhiNum) - 1; i >= 0; i-- {
		jinzhis = append(jinzhis, jinzhiNum[i])
	}
	return jinzhis
}
