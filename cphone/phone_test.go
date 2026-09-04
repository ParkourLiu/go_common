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

var ps = []string{
	"1340722",
	"1340723",
	"1340729",
	"1342999",
	"1345112",
	"1345115",
	"1346970",
	"1347741",
	"1347748",
	"1347751",
	"1354594",
	"1354598",
	"1359391",
	"1359396",
	"1359399",
	"1359740",
	"1359741",
	"1359745",
	"1361728",
	"1361729",
	"1367729",
	"1369737",
	"1369739",
	"1378991",
	"1387297",
	"1387298",
	"1388698",
	"1390722",
	"1397217",
	"1397294",
	"1398692",
	"1399799",
	"1470716",
	"1470722",
	"1500722",
	"1502727",
	"1502728",
	"1502731",
	"1502733",
	"1517151",
	"1517157",
	"1520722",
	"1527114",
	"1527115",
	"1527119",
	"1527122",
	"1571722",
	"1577101",
	"1579729",
	"1579926",
	"1580722",
	"1582688",
	"1582693",
	"1582793",
	"1582795",
	"1582796",
	"1582799",
	"1587183",
	"1587185",
	"1587187",
	"1587188",
	"1590861",
	"1592607",
	"1592608",
	"1592609",
	"1597198",
	"1597199",
	"1597228",
	"1780716",
	"1830728",
	"1837280",
	"1837283",
	"1837286",
	"1837287",
	"1837288",
	"1872731",
	"1872733",
	"1872734",
	"1872735",
	"1872736",
	"1872737",
	"1877115",
	"1877117",
	"1880722",
	"1887139",
	"1887232",
	"1887261",
	"1887263",
	"1952541",
	"1980728",
	"1987117",
	"1987186",
	"1988837",
}

func TestName(t *testing.T) {
	for _, v := range ps {
		r, _ := cphone.Find(v)
		fmt.Println(v, r.Province, r.City)
	}
}
