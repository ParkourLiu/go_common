package main

import (
	"fmt"
	"net"
)

//"188.131.240.128:2008", &proxy.Auth{User: "user@#3344", Password: "user@#3344telangpu"

func main() {
	// create a socks5 dialer
	//proxyDialer, err := proxy.SOCKS5("tcp", "188.131.240.128:2008", &proxy.Auth{User: "user@#3344", Password: "user@#3344telangpu"}, proxy.Direct)
	//if err != nil {
	//	fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
	//	os.Exit(1)
	//}
	//conn, err := proxyDialer.Dial("tcp", "119.29.173.65:20220")
	//if err != nil {
	//	fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
	//	os.Exit(1)
	//}
	//c, chans, reqs, err := ssh.NewClientConn(conn, "119.29.173.65:20220", &ssh.ClientConfig{
	//	User:            "root",
	//	Auth:            []ssh.AuthMethod{ssh.Password("tengxunyun@#3344")},
	//	HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	//})
	//sshClient := ssh.NewClient(c, chans, reqs)
	//sftpClient, err := sftp.NewClient(sshClient)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fis, err := sftpClient.ReadDir("/root/sftpTest/")
	//for k, v := range fis {
	//	fmt.Println(k, v.Name())
	//}
	var proxyConn net.Conn
	fmt.Println(proxyConn)
}
