package cbank

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var p Bankcard

type Bankcard struct {
	Bankcards []int //银行卡
	BankAddr  []BankAddr
}

type BankAddr struct {
	BankId   string //63020000
	BankName string //上海银行
	CardName string //阳光商旅信用卡
	CardType string //借记卡
	Province string
	City     string
}

func init() {
	fbs, err := ioutil.ReadFile("./bankData.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
		return
	}
	data := strings.TrimSpace(string(fbs))
	data = strings.ReplaceAll(data, "\r\n", "\n")
	data = strings.Trim(data, "\n")
	dataSplit := strings.Split(data, "\n")

	idcards := []int{}
	addrs := []BankAddr{}
	for _, v := range dataSplit {
		vSplit := strings.Split(v, "\t")
		if len(vSplit) != 7 {
			continue
		}
		codeStr := vSplit[0]
		bankId := vSplit[1]
		BankName := vSplit[2]
		CardName := vSplit[3]
		CardType := vSplit[4]
		province := vSplit[5]
		city := vSplit[6]
		codeInt, _ := strconv.Atoi(codeStr)
		idcards = append(idcards, codeInt)
		addrs = append(addrs, BankAddr{
			BankId:   bankId,
			BankName: BankName,
			CardName: CardName,
			CardType: CardType,
			Province: province,
			City:     city,
		})
	}
	p = Bankcard{
		Bankcards: idcards,
		BankAddr:  addrs,
	}
}

func Bankcard2Addr(Bankcard string) (addr BankAddr, ok bool) {
	if len(Bankcard) > 10 {
		Bankcard = Bankcard[:10]
	}
	nInt, err := strconv.Atoi(Bankcard)
	if err != nil {
		return
	}
	start := 0
	end := len(p.Bankcards) - 1
	for start <= end {
		mid := (start + end) / 2
		if nInt < p.Bankcards[mid] {
			end = mid - 1
		} else if nInt > p.Bankcards[mid] {
			start = mid + 1
		} else {
			return p.BankAddr[mid], true
		}
	}
	return
}

// 判断银联卡算法
func CheckBankcard(card string) (isUnion bool) {
	endNum, err := strconv.Atoi(card[len(card)-1:]) //取最后一位校验码
	if err != nil {
		return
	}
	card = card[:len(card)-1] //实际卡号

	double := true
	source := strings.Split(card, "")
	checksum := 0
	for i := len(source) - 1; i > -1; i-- {
		t, err := strconv.ParseInt(source[i], 10, 8)
		if err != nil {
			return
		}
		n := int(t)

		if double {
			n = n * 2
		}

		double = !double

		if n >= 10 {
			n = n - 9
		}

		checksum += n
	}
	controlDigit := checksum % 10
	if controlDigit != 0 {
		controlDigit = 10 - controlDigit
	}
	if endNum == controlDigit {
		isUnion = true
	}
	return
}
