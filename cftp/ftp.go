package cftp

import (
	"bytes"
	"github.com/jlaffaye/ftp"
	"time"
)

type FtpClient struct {
	addr string
	user string
	pwd  string
	conn *ftp.ServerConn
}

func NewFtpClient(addr, user, pwd string, checkConnect bool) (fc *FtpClient, err error) {
	fc = &FtpClient{
		addr: addr,
		user: user,
		pwd:  pwd,
	}
	c, err := ftp.Dial(addr, ftp.DialWithTimeout(20*time.Second))
	if err != nil {
		return
	}
	err = c.Login(user, pwd)
	if err != nil {
		return
	}
	fc.conn = c
	if checkConnect {
		go fc.online()
	}
	return
}

func (f *FtpClient) Upload(fileName string, fbs []byte) (err error) {
	err = f.conn.Stor(fileName, bytes.NewReader(fbs))
	return
}

func (f *FtpClient) online() {
	for {
		time.Sleep(time.Second * 5)
		_, err := f.conn.CurrentDir()
		if err == nil {
			continue
		}
		f.conn.Quit()
		tempF, err := NewFtpClient(f.addr, f.user, f.pwd, false)
		if err != nil {
			continue
		}
		f.conn = tempF.conn
	}

}
