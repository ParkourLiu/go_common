package cid_test

import (
	"fmt"
	"go_common/cid"
	"testing"
)

var IdUtil = cid.NewUnixId()

func TestName(t *testing.T) {
	fmt.Println(IdUtil.NextId())
	fmt.Println(IdUtil.NextId62Random(17))
	fmt.Println(IdUtil.NextId62Random(12))
}

// 并发
func TestName2(t *testing.T) {
	fmt.Println(111111111111)
	c := make(chan bool, 300)
	ss := make([]string, 9999999)
	for i := 0; i < 9999999; i++ {
		c <- true
		go func(ii int) {
			id := IdUtil.NextId62()
			ss[ii] = id
			<-c
		}(i)
	}
	m := map[string]bool{}
	for _, s := range ss {
		if m[s] {
			fmt.Println("重复", s)
			return
		}
		m[s] = true
	}
	fmt.Println("完成")
}

func TestCreatZhiShuIds(t *testing.T) {
	fmt.Println(cid.CreatZhiShuIDs(30, 66666))
	dllId := cid.Id2ZhiShuId("dsds", 8765876548975, 10)
	fmt.Println(dllId)
}

func TestId2ZhiShuId(t *testing.T) { //5lm1vouaq 8ef2tj9g3 b783rdolg e014p83q
	cids := cid.CreatZhiShuIDs(8765876548975, 10)
	fmt.Println(cids)
	fmt.Println(len(cids))
	fmt.Println(cid.Id2ZhiShuId("a", 8765876548975, 1000))
}
