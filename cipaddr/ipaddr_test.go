package cipaddr_test

import (
	"fmt"
	"go_common/cipaddr"
	"testing"
)

func TestUint322Ip(t *testing.T) {
	a, err := cipaddr.IpToUint32("180.167.221.138")
	fmt.Println(a, err)
	ip := cipaddr.Uint32ToIP(a)
	fmt.Println(ip)
}

func TestIp2Addr(t *testing.T) {
	fmt.Println(cipaddr.Ip2Addr("180.167.221.138"))

}
