package csftp

import (
	"github.com/pkg/sftp"
	"go_common/clogs"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"time"
)

var (
	log        *clogs.Log
	checkBreak = false //是否停止检查连通性的开关
)

func init() {
	log = clogs.NewLog("7")
}

type SftpClient struct {
	aaa          string
	sSH_ADDR     string
	sSH_USER     string
	sSH_PASSWORD string
	sshClient    *ssh.Client
	sftpClient   *sftp.Client
}

func NewSftpClient(SSH_ADDR, SSH_USER, SSH_PASSWORD string) *SftpClient {
agenSSH:
	sshClient, err := ssh.Dial("tcp", SSH_ADDR, &ssh.ClientConfig{
		User:            SSH_USER,
		Auth:            []ssh.AuthMethod{ssh.Password(SSH_PASSWORD)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		log.Error(err)
		time.Sleep(5 * time.Second)
		goto agenSSH
	}
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		log.Error(err)
		time.Sleep(5 * time.Second)
		goto agenSSH
	}
	sc := &SftpClient{
		sSH_ADDR:     SSH_ADDR,
		sSH_USER:     SSH_USER,
		sSH_PASSWORD: SSH_PASSWORD,
		sshClient:    sshClient,
		sftpClient:   sftpClient,
	}
	go sc.checkConn() //监控联通性
	return sc
}

//用pwd检测连通性，保证断开后会自联
func (s *SftpClient) checkConn() {
	for {
		if checkBreak { //是否停止检查连接,并关闭通道
			if s.sshClient != nil {
				_ = s.sshClient.Close()
			}
			if s.sftpClient != nil {
				_ = s.sftpClient.Close()
			}
			checkBreak = false //复位开关
			log.Info("sshSftpCilent conn is close")
			return
		}

		//用pwd检测连通性
		_, err := s.sftpClient.Getwd()
		if err != nil {
			//连接出错,开始重新初始化
			log.Error(err)
			sc := NewSftpClient(s.sSH_ADDR, s.sSH_USER, s.sSH_PASSWORD)
			s.sshClient = sc.sshClient
			s.sftpClient = sc.sftpClient
		}
		time.Sleep(5 * time.Second)
		continue
	}
}

//每次使用时，用内置的GetSftpClient方法而不要自己定义一个变量再使用，因为使用GetSftpClient方法可实现断开自动重新重连
func (s *SftpClient) GetSftpClient() *sftp.Client {
	return s.sftpClient
}

//每次使用时，用内置的GetSshClient方法而不要自己定义一个变量再使用，因为使用GetSshClient方法可实现断开自动重新重连
func (s *SftpClient) GetSshClient() *ssh.Client {
	return s.sshClient
}

func (s *SftpClient) Close() {
	checkBreak = true
}

func (s *SftpClient) Upload(fileBytes []byte, dstPath, dstFileName string) error {
	fileName := dstPath + dstFileName
	tmpFileName := fileName + ".tmp"
	dstFile, err := s.sftpClient.Create(tmpFileName) //远程
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = dstFile.Write(fileBytes)
	if err != nil {
		return err
	}
	return s.sftpClient.Rename(tmpFileName, fileName) //重命名
}

func (s *SftpClient) Download(srcPath string) ([]byte, error) {
	srcFile, err := s.sftpClient.Open(srcPath) //远程
	if err != nil {
		return nil, err
	}
	defer srcFile.Close()
	return ioutil.ReadAll(srcFile)
}
func (s *SftpClient) Remove(path string) error {
	return s.sftpClient.Remove(path)
}
