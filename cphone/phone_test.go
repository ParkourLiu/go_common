package cphone_test

import (
	"fmt"
	"testing"
)
import "go_common/cphone"

func TestFind(t *testing.T) {
	cphone.NewPhone(`D:\Go\GoWorkkkkkk\src\go_common\cphone\phone.dat`)
	r, err := cphone.Find("17671774535")
	fmt.Println(r, err)
}
