package cdpapi

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	test("我的回归3都给我iu全国丢高度为孤独和我给丢额外给都给u我干爹我大概iu额外给i的外观丢为其单独违规idyu的我国其余都给我i一点归额外的归额外各位给我额 奋斗史非让我分为氛围 俄方微风额外分为氛围24如WS")
	test("fdhiuy2ty87f")
	test("dh7932h9h98")
	test("yd98g872gf87342g8f7fer")
}

func test(txt string) {
	a := []byte(txt)
	bs, err := Enc(a)
	fmt.Println(1, string(bs), err)
	bs, err = Dec(bs)
	fmt.Println(2, string(bs), err)
}
