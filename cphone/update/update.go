package main

import (
	"bytes"
	"fmt"
	"go_common/clogs"
	"io/ioutil"
	"sort"
	"strings"
)

var (
	log             = clogs.NewLog("7", false)
	oldPhoneDatPath = `D:\Go\GoWorkkkkkk\src\go_common\cphone\phone.dat`                             //文件的路径
	newPhoneDatPath = `D:\Go\GoWorkkkkkk\src\go_common\cphone\update\phone-qqzeng-202209-492961.txt` //文件的路径
	outputPath      = `D:\Go\GoWorkkkkkk\src\go_common\cphone\update\phone.dat`                      //文件的路径
	phoneMap        = map[string]PhoneRecord{}
	newPhoneMap     = map[string]PhoneRecord{}
)

type PhoneRecord struct {
	//PhoneNum string
	Province string
	City     string
	CardType string
}

func main() {
	content, err := ioutil.ReadFile(oldPhoneDatPath)
	if err != nil {
		log.Error(err)
		return
	}
	contentStr := string(content)
	phoneInfoLineList := strings.Split(contentStr, "\n") //按行切割
	for _, v := range phoneInfoLineList {
		phoneInfo := strings.Split(v, "\t") //按tab切割出具体列  1300000	甘肃	临夏	中国联通网络
		if len(phoneInfo) != 4 {
			log.Warn(len(phoneInfo), v)
			continue
		}
		phoneMap[phoneInfo[0]] = PhoneRecord{
			//PhoneNum: phoneInfo[0], //号码
			Province: phoneInfo[1], //省
			City:     phoneInfo[2], //市(或者直辖县等)
			//CardType: phoneInfo[3], //中国联通网络   中国移动网络  中国电信网络
		}
	}
	//---------------------------------------------------------------------------------------------
	newContent, err := ioutil.ReadFile(newPhoneDatPath)
	if err != nil {
		log.Error(err)
		return
	}
	newContentStr := string(newContent)
	newContentStr = strings.ReplaceAll(newContentStr, "\r\n", "\n")
	newPhoneInfoLineList := strings.Split(newContentStr, "\n") //按行切割
	phoneList := []string{}
	for i, v := range newPhoneInfoLineList {
		if i == 0 { //第一行是字段头
			continue
		}
		phoneInfo := strings.Split(v, "\t") //按tab切割出具体列  170	1703075	陕西	西安	移动/虚拟	710000	029	610100
		if len(phoneInfo) != 8 {
			log.Warn(len(phoneInfo), v)
			continue
		}
		phoneList = append(phoneList, phoneInfo[1])
		newPhoneMap[phoneInfo[1]] = PhoneRecord{
			//PhoneNum: phoneInfo[0], //号码
			Province: phoneInfo[2], //省
			City:     phoneInfo[3], //市(或者直辖县等)
			CardType: phoneInfo[4], //中国联通网络   中国移动网络  中国电信网络
		}
	}

	//-----------------------------------------------------------------------------------------------
	//多区域查询
	count := 0
	temp := map[string]int{}
	for k, v := range newPhoneMap {
		if strings.Contains(v.City, "/") {
			log.Info(k, v)
			count++
			temp[v.City] = temp[v.City] + 1
		}
	}
	log.Info("多区域总数量：", count)
	for k, v := range temp {
		log.Info(v, k)
	}

	//老的有新的没有
	for k, v := range phoneMap {
		if _, ok := newPhoneMap[k]; !ok {
			log.Info("新版丢失：", k, v)
		}
	}
	//-------------------------------------------------------------------------------------------------------
	sort.Strings(phoneList) //排好序
	buf := bytes.Buffer{}
	for _, phone := range phoneList {
		v := newPhoneMap[phone]
		if v.City == "仙桃/潜江/天门" {
			v.City = "江汉"
		}

		if strings.Contains(v.City, "/") {
			v.City = v.City[:strings.Index(v.City, "/")]
		}
		buf.WriteString(fmt.Sprintf("%s\t%s\t%s\t%s\n", phone, v.Province, v.City, v.CardType))
	}
	ioutil.WriteFile(outputPath, buf.Bytes(), 0644)
}
