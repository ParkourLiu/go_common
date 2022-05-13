package cphone

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

var (
	phoneDatPath = "" //文件的路径
	phoneMap     = map[string]PhoneRecord{}
)

type PhoneRecord struct {
	//PhoneNum string
	Province string
	City     string
	//CardType string
}

func (pr *PhoneRecord) String() string {
	pj, _ := json.Marshal(pr)
	return string(pj)
}

func NewPhone(datpath string) (err error) {
	content, err := ioutil.ReadFile(datpath)
	if err != nil {
		return
	}
	content = bytes.Replace(content, []byte("市"), []byte(""), -1)
	content = bytes.Replace(content, []byte("彝族自治州"), []byte(""), -1)
	content = bytes.Replace(content, []byte("藏族自治州"), []byte(""), -1)
	content = bytes.Replace(content, []byte("朝鲜族自治州"), []byte(""), -1)
	content = bytes.Replace(content, []byte("回族自治州"), []byte(""), -1)
	content = bytes.Replace(content, []byte("白族自治州"), []byte(""), -1)
	content = bytes.Replace(content, []byte("苗族侗族自治州"), []byte(""), -1)
	content = bytes.Replace(content, []byte("布依族苗族自治州"), []byte(""), -1)
	content = bytes.Replace(content, []byte("蒙古自治州"), []byte(""), -1)
	content = bytes.Replace(content, []byte("哈萨克自治州"), []byte(""), -1)
	content = bytes.Replace(content, []byte("土家族苗族自治州"), []byte(""), -1)
	content = bytes.Replace(content, []byte("壮族苗族自治州"), []byte(""), -1)
	content = bytes.Replace(content, []byte("傣族自治州"), []byte(""), -1)
	content = bytes.Replace(content, []byte("柯尔克孜自治州"), []byte(""), -1)
	content = bytes.Replace(content, []byte("傈僳族自治州"), []byte(""), -1)
	content = bytes.Replace(content, []byte("傣族景颇族自治州"), []byte(""), -1)
	content = bytes.Replace(content, []byte("藏族羌族自治州"), []byte(""), -1)
	content = bytes.Replace(content, []byte("哈尼族"), []byte(""), -1)
	content = bytes.Replace(content, []byte("蒙古族"), []byte(""), -1)
	content = bytes.Replace(content, []byte("地区"), []byte(""), -1)
	contentStr := string(content)
	phoneInfoLineList := strings.Split(contentStr, "\n") //按行切割
	for _, v := range phoneInfoLineList {
		phoneInfo := strings.Split(v, "\t") //按tab切割出具体列  1300000	甘肃	临夏	中国联通网络
		if len(phoneInfo) != 4 {
			continue
		}
		phoneMap[phoneInfo[0]] = PhoneRecord{
			//PhoneNum: phoneInfo[0], //号码
			Province: phoneInfo[1], //省
			City:     phoneInfo[2], //市(或者直辖县等)
			//CardType: phoneInfo[3], //中国联通网络   中国移动网络  中国电信网络
		}
	}
	return
}

func Find(phone_num string) (pr *PhoneRecord, err error) {
	if len(phone_num) < 7 || len(phone_num) > 11 {
		return nil, errors.New("illegal phone length")
	}

	p := phoneMap[phone_num[0:7]]
	return &p, nil
}

func Addr() {
	sortAddr := []string{}
	s := map[string]string{}
	for _, v := range phoneMap {
		if _, ok := s[v.City]; !ok {
			s[v.City] = v.Province
			sortAddr = append(sortAddr, v.Province+"\t"+v.City+"\n")
		}
	}
	sort.Strings(sortAddr)
	sb := bytes.Buffer{}
	for _, v := range sortAddr {
		sb.WriteString(v)
	}
	ioutil.WriteFile(`D:\Go\GoWorkkkkkk\src\go_common\cphone\Addr.dat`, sb.Bytes(), 0644)
}

func Province2City(Province string) {
	s := map[string]string{}
	for _, v := range phoneMap {
		if Province == v.Province {
			s[v.City] = v.Province
		}
	}

	for k, v := range s {
		fmt.Println(v, k)
	}
}
