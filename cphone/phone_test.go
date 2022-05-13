package cphone_test

import (
	"bytes"
	"fmt"
	"go_common/cphone"
	"io/ioutil"
	"sort"
	"strings"
	"testing"
	"time"
)

func init() {
	cphone.NewPhone(`D:\Go\GoWorkkkkkk\src\go_common\cphone\phone.dat`)
}
func TestFind(t *testing.T) {
	r, err := cphone.Find("1997558")
	fmt.Println(r, err)
	for {
		time.Sleep(time.Hour)
	}
}
func TestAddr(t *testing.T) {
	cphone.Addr()

}

func TestProvince2City(t *testing.T) {
	cphone.Province2City("湖北")
}

type PhoneRecord struct {
	PhoneNum string
	Province string
	City     string
	CardType string
}

func TestSort(t *testing.T) {
	phoneMap := map[string]PhoneRecord{}
	phoneList := []PhoneRecord{}
	content, err := ioutil.ReadFile(`D:\Go\GoWorkkkkkk\src\go_common\cphone\phoneSortByCity.txt`)
	if err != nil {
		return
	}
	contentStr := string(content)
	phoneInfoLineList := strings.Split(contentStr, "\n") //按行切割
	for _, v := range phoneInfoLineList {
		phoneInfo := strings.Split(v, "\t") //按tab切割出具体列  1300000	甘肃	临夏	中国联通网络
		if len(phoneInfo) != 4 {
			continue
		}
		p := PhoneRecord{
			PhoneNum: phoneInfo[0], //号码
			Province: phoneInfo[1], //省
			City:     phoneInfo[2], //市(或者直辖县等)
			CardType: phoneInfo[3], //中国联通网络   中国移动网络  中国电信网络
		}
		phoneMap[phoneInfo[0]] = p
		phoneList = append(phoneList, p)
	}

	sort.SliceStable(phoneList, func(i int, j int) bool {
		return phoneList[i].City < phoneList[j].City
	})

	sb := bytes.Buffer{}
	for _, v := range phoneList {
		sb.WriteString(v.PhoneNum + "\t" + v.Province + "\t" + v.City + "\t" + v.CardType + "\n")
	}
	ioutil.WriteFile(`D:\Go\GoWorkkkkkk\src\go_common\cphone\phoneSortByCity.txt`, sb.Bytes(), 0644)
}
