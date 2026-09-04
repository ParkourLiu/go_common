package cid

import (
	"crypto/md5"
)

// 根据id获取对应的倍数id
var (
	//zhishu = int64(9)
	zhishu = int64(786499599943985)
)

func CreatZhiShuIDs(min, count int64) (ids []string) {
	for i := int64(0); i < count; i++ {
		ids = append(ids, jinzhi36ToString(i*zhishu+min))
	}
	return
}

//
//func trim(min int64) (newMin int64) {
//	for i := int64(0); ; i++ { //修剪最小值
//		tmp := zhishu * i
//		if tmp >= min {
//			newMin = tmp
//			break
//		}
//	}
//	return
//}

func Id2ZhiShuId(id string, min, count int64) string {
	c := Id2Int(id)
	c = c % count
	c = c*zhishu + min
	return jinzhi36ToString(c)
}

func Id2Int(id string) (z int64) {
	bs := md5.Sum([]byte(id))
	for i, d := range bs {
		weishu := int64(1)
		for j := 0; j < i; j++ {
			weishu = weishu * 10
		}
		z = z + weishu*int64(d)
	}
	return
}
