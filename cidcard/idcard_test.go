package cidcard_test

import (
	"fmt"
	"go_common/cidcard"
	"math/rand"
	"runtime/debug"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	fmt.Println(cidcard.Idcard2Addr("330700198811082304"))

}
func TestAddr(t *testing.T) {
	for {
		fmt.Println(cidcard.Idcard2Addr("653125199601173212"))
		debug.FreeOSMemory()
		time.Sleep(time.Second)
	}
}
func TestIsIdcard(t *testing.T) {
	fmt.Println(cidcard.IsIdcard("421087199602293223"))
	fmt.Println(310 % 11)
}
func TestCalcLast(t *testing.T) {
	fmt.Println(cidcard.CalcLast("310123200001011297"))
}

func TestIdcard2Addr(t *testing.T) {
	fmt.Println(randInt(1, 2))
}

func randInt(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Intn(max-min+1) + min
}
