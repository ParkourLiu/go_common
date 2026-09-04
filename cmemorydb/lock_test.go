package cmemorydb

import (
	"fmt"
	"testing"
)

var ll = newLock()

func TestName(t *testing.T) {
	for i := 0; i < 9999; i++ {
		go a(i)
		go c(i)
		go b(i)
	}

}
func a(i int) {
	ll.Lock()
	defer ll.UnLock()
	fmt.Println(i, 11)
}
func b(i int) {
	ll.Lock()
	defer ll.UnLock()
	fmt.Println(i, 222)
}
func c(i int) {
	ll.Lock()
	defer ll.UnLock()
	fmt.Println(i, 333)
}
