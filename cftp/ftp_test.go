package cftp_test

import (
	"fmt"
	"go_common/cftp"
	"os"
	"testing"
	"time"
)

var ftpClient *cftp.FtpClient

func init() {
	var err error
	ftpClient, err = cftp.NewFtpClient("172.16.5.129:2021", "ftpUser", "Ansight@123", true)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
func TestName(t *testing.T) {
	for i := 0; ; i++ {
		time.Sleep(time.Second * 2)
		err := ftpClient.Upload(fmt.Sprint(i)+"aa.txt", []byte("abcd"))
		if err != nil {
			fmt.Println(111, err)
			continue
		}
		fmt.Println(i)
	}

	fmt.Println("完成")
}
