package csftp

import (
	"github.com/pkg/sftp"
	"go_common/clogs"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net"
	"strings"
	"time"
)

type SftpClient struct {
	sSH_ADDR     string
	sSH_USER     string
	sSH_PASSWORD string
	sshClient    *ssh.Client
	sftpClient   *sftp.Client
	proxyInfo    *ProxyInfo //代理信息
	conn         net.Conn
	checkBreak   bool //是否停止检查连通性的开关
	log          *clogs.Log
}

//代理信息的结构体
type ProxyInfo struct {
	Addres       string //"127.0.0.1:2008"
	AuthUser     string //代理用户
	AuthPassword string //代理密码
}

func NewSftpClient(sshAddr, sshUser, sshPassword string, proxyInfo *ProxyInfo, checkConnect bool, log *clogs.Log) *SftpClient {
	var sshClient *ssh.Client
	var err error
	var proxyConn net.Conn
agenSSH:
	if proxyInfo == nil { //不挂代理创建ssh连接
		sshClient, err = ssh.Dial("tcp", sshAddr, &ssh.ClientConfig{
			User:            sshUser,
			Auth:            []ssh.AuthMethod{ssh.Password(sshPassword)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		})
		if err != nil {
			log.Error(err)
			time.Sleep(5 * time.Second)
			goto agenSSH
		}
	} else { //需要挂代理创建ssh连接
		//创建代理拨号器
		proxyDialer, err := proxy.SOCKS5("tcp", proxyInfo.Addres, &proxy.Auth{User: proxyInfo.AuthUser, Password: proxyInfo.AuthPassword}, proxy.Direct)
		if err != nil {
			log.Error(err)
			time.Sleep(5 * time.Second)
			goto agenSSH
		}
		//用代理拨号器连接目标服务器:端口
		proxyConn, err = proxyDialer.Dial("tcp", sshAddr)
		if err != nil {
			log.Error(err)
			time.Sleep(5 * time.Second)
			goto agenSSH
		}
		//使用代理连接器 创建ssh连接
		c, chans, reqs, err := ssh.NewClientConn(proxyConn, sshAddr, &ssh.ClientConfig{
			User:            sshUser,
			Auth:            []ssh.AuthMethod{ssh.Password(sshPassword)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		})
		if err != nil {
			log.Error(err)
			if proxyConn != nil {
				_ = proxyConn.Close() //关闭上一步创建的连接器
			}
			time.Sleep(5 * time.Second)
			goto agenSSH
		}
		//创建出sshCilent
		sshClient = ssh.NewClient(c, chans, reqs)
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		log.Error(err)
		time.Sleep(5 * time.Second)
		goto agenSSH
	}
	sc := &SftpClient{
		sSH_ADDR:     sshAddr,
		sSH_USER:     sshUser,
		sSH_PASSWORD: sshPassword,
		sshClient:    sshClient,
		sftpClient:   sftpClient,
		proxyInfo:    proxyInfo,
		conn:         proxyConn,
		log:          log,
	}
	if checkConnect {
		go sc.checkConn() //监控连通性
	}

	return sc
}

//用pwd检测连通性，保证断开后会自联
func (c *SftpClient) checkConn() {
	for {
		if c.checkBreak { //是否停止检查连接,并关闭通道
			if c.sftpClient != nil {
				_ = c.sftpClient.Close()
			}
			if c.conn != nil {
				_ = c.conn.Close()
			}
			if c.sshClient != nil {
				_ = c.sshClient.Close()
			}
			c.checkBreak = false //复位开关
			c.log.Info("sshSftpCilent conn is close")
			return
		}

		//用pwd检测连通性
		_, err := c.sftpClient.Getwd()
		if err != nil {
			//连接出错,开始重新初始化
			c.log.Error(err)

			if c.sftpClient != nil {
				_ = c.sftpClient.Close()
			}
			if c.conn != nil {
				_ = c.conn.Close()
			}
			if c.sshClient != nil {
				_ = c.sshClient.Close()
			}
			sc := NewSftpClient(c.sSH_ADDR, c.sSH_USER, c.sSH_PASSWORD, c.proxyInfo, false, c.log) //创建新的链接
			c.sshClient = sc.sshClient
			c.sftpClient = sc.sftpClient
			c.conn = sc.conn
		}
		time.Sleep(5 * time.Second)
		continue
	}
}

//每次使用时，用内置的GetSftpClient方法而不要自己定义一个变量再使用，因为使用GetSftpClient方法可实现断开自动重新重连
func (c *SftpClient) GetSftpClient() *sftp.Client {
	return c.sftpClient
}

//每次使用时，用内置的GetSshClient方法而不要自己定义一个变量再使用，因为使用GetSshClient方法可实现断开自动重新重连
func (c *SftpClient) GetSshClient() *ssh.Client {
	return c.sshClient
}

//关闭ssh和sftp连接，并停止连接性检查
func (c *SftpClient) Close() {
	c.checkBreak = true //停止检查连接
	if c.sftpClient != nil {
		_ = c.sftpClient.Close()
	}
	if c.conn != nil {
		_ = c.conn.Close()
	}
	if c.sshClient != nil {
		_ = c.sshClient.Close()
	}
}

//上传文件
func (c *SftpClient) Upload(fileBytes []byte, dstPath, dstFileName string) error {
	if !strings.HasSuffix(dstPath, "/") {
		dstPath = dstPath + "/"
	}
	fileName := dstPath + dstFileName
	tmpFileName := fileName + ".tmp"
	err := c.sftpClient.MkdirAll(dstPath) //创建目录
	if err != nil {
		return err
	}
	dstFile, err := c.sftpClient.Create(tmpFileName) //远程
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = dstFile.Write(fileBytes)
	if err != nil {
		return err
	}
	return c.sftpClient.Rename(tmpFileName, fileName) //重命名
}

//下载文件
func (c *SftpClient) Download(srcPath string) ([]byte, error) {
	srcFile, err := c.sftpClient.Open(srcPath) //远程
	if err != nil {
		return nil, err
	}
	defer srcFile.Close()
	return ioutil.ReadAll(srcFile)
}

//删除文件或者文件夹
func (c *SftpClient) Remove(path string) error {
	return c.sftpClient.Remove(path)
}

//移动文件夹
func (c *SftpClient) Rename(oldPath, newPath string) error {
	return c.sftpClient.Rename(oldPath, newPath)
}
