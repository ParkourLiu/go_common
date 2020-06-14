package csftp_test

import (
	"fmt"
	"go_common/csftp"
	"testing"
	"time"
)

var (
	sshSftpClient *csftp.SftpClient
)

func init() {
	sshSftpClient = csftp.NewSftpClient("172.16.5.137:22", "root", "1")
}

func Test_Create(t *testing.T) {
	defer sshSftpClient.Close()
	for i := 0; i < 1000; i++ {
		s, err := sshSftpClient.GetSftpClient().Getwd() //每次使用时，用内置的GetSftpClient方法而不要自己定义一个变量再使用，因为使用GetSftpClient方法可实现断开自动重新重连
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(s)
		}
		time.Sleep(3 * time.Second)
	}

}
func Test_Close(t *testing.T) {
	sshSftpClient.Close()
	time.Sleep(5 * time.Minute)
}
func Test_Upload(t *testing.T) {
	defer sshSftpClient.Close()
	err := sshSftpClient.Upload([]byte("123"), "/", "123")
	if err != nil {
		t.Error(err)
		return
	}
}
func Test_Download(t *testing.T) {
	defer sshSftpClient.Close()
	b, err := sshSftpClient.Download("/123")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(b))
}

func TestSftpClient_Remove(t *testing.T) {
	defer sshSftpClient.Close()
	err := sshSftpClient.Remove("/root/aaa/")
	if err != nil {
		t.Error(err)
		return
	}
}

