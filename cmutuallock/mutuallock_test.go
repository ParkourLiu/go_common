package cmutuallock_test

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	a := "aaa"
	ab := []byte(a)
	for k, v := range ab {
		fmt.Println(k, v)
	}
}
