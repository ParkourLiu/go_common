package cbank

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	fmt.Println(Bankcard2Addr("353003094793939"))
	fmt.Println(CheckBankcard("353003094793939"))
}
