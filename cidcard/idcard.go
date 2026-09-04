package cidcard

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var p idcard

func init() {
	fbs, err := ioutil.ReadFile("./idcardData.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
		return
	}
	data := strings.TrimSpace(string(fbs))
	data = strings.ReplaceAll(data, "\r\n", "\n")
	data = strings.Trim(data, "\n")
	dataSplit := strings.Split(data, "\n")

	idcards := []int32{}
	addrs := []addr{}
	for _, v := range dataSplit {
		vSplit := strings.Split(v, "\t")
		if len(vSplit) != 4 {
			continue
		}
		codeStr := vSplit[0]
		province := vSplit[1]
		city := vSplit[2]
		district := vSplit[3]
		codeInt, _ := strconv.Atoi(codeStr)
		idcards = append(idcards, int32(codeInt))
		addrs = append(addrs, addr{
			Province: province,
			City:     city,
			District: district,
		})
	}
	p = idcard{
		Idcards: idcards,
		Addrs:   addrs,
	}
}

type idcard struct {
	Idcards []int32
	Addrs   []addr
}

type addr struct {
	Province string
	City     string
	District string
}

func Idcard2Addr(idcard6 string) (addr addr, ok bool) {
	if len(idcard6) > 6 {
		idcard6 = idcard6[:6]
	}
	nInt, err := strconv.Atoi(idcard6)
	if err != nil {
		return
	}
	idcardInt32 := int32(nInt)
	start := 0
	end := len(p.Idcards) - 1
	for start <= end {
		mid := (start + end) / 2
		if idcardInt32 < p.Idcards[mid] {
			end = mid - 1
		} else if idcardInt32 > p.Idcards[mid] {
			start = mid + 1
		} else {
			return p.Addrs[mid], true
		}
	}
	return
}

var (
	xiShu  = []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	jieGuo = map[int]string{
		0:  "1",
		1:  "0",
		2:  "X",
		3:  "9",
		4:  "8",
		5:  "7",
		6:  "6",
		7:  "5",
		8:  "4",
		9:  "3",
		10: "2",
	}
)

func IsIdcard(idcard string) (isIdcard bool) {
	split := strings.Split(idcard, "")
	if len(split) != 18 {
		return
	}
	add := 0
	for i, v := range split[:len(split)-1] {
		nInt, err := strconv.Atoi(v)
		if err != nil {
			return
		}
		add = add + (nInt * xiShu[i])
	}
	if split[len(split)-1] == "x" {
		split[len(split)-1] = "X"
	}
	if jieGuo[add%11] != split[len(split)-1] {
		return
	}
	return true
}
func CalcLast(idcard string) (last string) {
	split := strings.Split(idcard, "")
	if len(split) != 17 {
		return
	}
	add := 0
	for i, v := range split {
		nInt, err := strconv.Atoi(v)
		if err != nil {
			return
		}
		add = add + (nInt * xiShu[i])
	}
	return jieGuo[add%11]
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
func randInt(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Intn(max-min+1) + min
}
func randDate() string {
	return time.Now().AddDate(-randInt(5, 80), -randInt(1, 12), -randInt(1, 30)).Format("20060102")
}
func Check() {
	fbs, err := ioutil.ReadFile("./data.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	data := strings.TrimSpace(string(fbs))
	data = strings.ReplaceAll(data, "\r\n", "\n")
	data = strings.Trim(data, "\n")
	dataSplit := strings.Split(data, "\n")
	file, err := os.OpenFile("D:/idcard.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("D:/idcard.txt", err)
		return
	}
	for _, v := range dataSplit {
		vSplit := strings.Split(v, "\t")
		if len(vSplit) != 4 {
			continue
		}
		codeStr := vSplit[0]
		idcards := codeStr + fmt.Sprintf("%s%d", randDate(), randInt(100, 999))
		idcards = idcards + CalcLast(idcards)
		fmt.Println(idcards)
		file.WriteString(idcards + "\n")
		//	headMap := map[string]string{
		//		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.35 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36 Edg/106.0.1418.52",
		//		"Host":       "www.haoshudi.com",
		//	}
		//agen:
		//	jbs, err := cmethod.HTTPrequest("GET", fmt.Sprintf("https://www.haoshudi.com/api/id/query/?userid=%s", idcards), headMap, nil, 30)
		//	if err != nil {
		//		fmt.Println("请求错误", codeStr, err)
		//		time.Sleep(time.Second * 5)
		//		goto agen
		//	}
		//	var c checkStruct
		//	err = json.Unmarshal(jbs, &c)
		//	if err != nil {
		//		fmt.Println("数据转换错误", codeStr, string(jbs), err)
		//		time.Sleep(time.Second * 5)
		//		goto agen
		//	}
		//	if !c.Status {
		//		fmt.Println("结果错误", codeStr, string(jbs))
		//		time.Sleep(time.Second * 5)
		//		goto agen
		//	}
		//	if len(c.Data.NewAddress) != 0 {
		//		c.Data.Address = strings.Join(c.Data.NewAddress, " ")
		//	}
		//	c.Data.Address = strings.ReplaceAll(c.Data.Address, " ", "	")
		//	file.WriteString(fmt.Sprintf("%s	%s\n", codeStr, c.Data.Address))
		//	time.Sleep(time.Second*5 + (time.Millisecond * time.Duration(rand.Intn(20))))
	}
}

type checkStruct struct {
	Status bool   `json:"status"`
	Msg    string `json:"msg"`
	Data   info   `json:"data"`
}
type info struct {
	Address    string   `json:"address"`
	NewAddress []string `json:"new_address"`
}
