package cregcenter

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	regInfo *RegInfo
	token   string
)

func InitRegInfo(ri *RegInfo) (err error) {
	if ri.LoginName == "" || ri.LoginPassword == "" || ri.RegCenterAddr == "" || ri.Workspace == "" || ri.ServerName == "" || ri.Port == "" {
		return fmt.Errorf("Registration information cannot have empty fields")
	}
	regInfo = ri

	err = login()
	if err != nil {
		return
	}
	go regServer()
	return
}

func NewCallClient(serverName string) (ct *CallServerClient, err error) {
	if regInfo == nil {
		return ct, fmt.Errorf("Please initialize the registration information before performing this operation")
	}
	ct = &CallServerClient{l: &sync.RWMutex{}, serverName: serverName, hosts: []string{}}
	err = ct.findServer()
	if err != nil {
		return
	}
	go ct.callFindServer()
	return ct, err
}
func login() (err error) {
	p := paramData{
		LoginName:     regInfo.LoginName,
		LoginPassWord: regInfo.LoginPassword,
	}
	r, err := httpRequest(fmt.Sprintf("http://%s/login", regInfo.RegCenterAddr), p)
	if err != nil {
		fmt.Printf("%s : reg server auth error %s \n", time.Now().Format("2006-01-02 15:04:05"), err.Error())
		return
	}
	if r.Code != code200 {
		fmt.Printf("%s : reg server auth error %s \n", time.Now().Format("2006-01-02 15:04:05"), r.Msg)
		return fmt.Errorf("%d:%s", r.Code, r.Msg)
	}
	token = r.Token
	return
}
func regServer() {
	for {
		reg()
		time.Sleep(5 * time.Second)
	}
}

func reg() {
	p := paramData{
		Workspace:  regInfo.Workspace,
		ServerName: regInfo.ServerName,
		Port:       regInfo.Port,
	}
	r, err := httpRequest(fmt.Sprintf("http://%s/regServer", regInfo.RegCenterAddr), p)
	if err != nil {
		fmt.Printf("%s : reg server error %s \n", time.Now().Format("2006-01-02 15:04:05"), err.Error())
		return
	}
	if r.Code == code301 {
		_ = login()
	} else if r.Code != code200 {
		fmt.Printf("%s : reg server error %s \n", time.Now().Format("2006-01-02 15:04:05"), r.Msg)
		return
	}
}

func GetServers(healthy int) (servers []*paramData) {
	p := paramData{
		Workspace: regInfo.Workspace,
		Healthy:   healthy,
	}
	r, err := httpRequest(fmt.Sprintf("http://%s/getServers", regInfo.RegCenterAddr), p)
	if err != nil {
		fmt.Printf("%s : get servers error %s \n", time.Now().Format("2006-01-02 15:04:05"), err.Error())
		return
	}
	if r.Code == code301 {
		_ = login()
	} else if r.Code != code200 {
		fmt.Printf("%s : get servers error %s \n", time.Now().Format("2006-01-02 15:04:05"), r.Msg)
		return
	}
	return r.Servers
}

func getServer(serverName string) (servers []*paramData) {
	p := paramData{
		Workspace:  regInfo.Workspace,
		ServerName: serverName,
		Healthy:    1,
	}
	r, err := httpRequest(fmt.Sprintf("http://%s/getServer", regInfo.RegCenterAddr), p)
	if err != nil {
		fmt.Printf("%s : get server error %s \n", time.Now().Format("2006-01-02 15:04:05"), err.Error())
		return
	}
	if r.Code == code200 {
		return r.Servers
	}
	return
}
func httpRequest(url string, p paramData) (r returnData, err error) {
	var EncDec2 = func(byt []byte, divisor int) []byte {
		for i, v := range byt {
			byt[i] = (byte(i+divisor) & (^v)) | (v & (^byte(i + divisor)))
		}
		return byt
	}
	var post = func(url string, bodybytes []byte) (r returnData, err error) {
		body := bytes.NewReader(bodybytes)
		request, err := http.NewRequest("POST", url, body) //创建请求体
		if err != nil {
			return
		}
		if request != nil {
			if request.Body != nil {
				defer func() {
					_ = request.Body.Close()
					request.Close = true
				}()
			}
		}
		request.Header.Set("Token", token)
		httpClient := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
					conn, err := net.Dial("tcp", addr)
					if err != nil {
						return nil, err
					}
					_ = conn.SetDeadline(time.Now().Add(time.Second * time.Duration(10))) //设置 发送+接受  数据超时时间
					return conn, nil
				},
			},
			Timeout: time.Second * time.Duration(10),
		}

		resp, err := httpClient.Do(request)
		if resp != nil {
			if resp.Body != nil {
				defer func() {
					_ = resp.Body.Close()
					resp.Close = true
				}()
			}
		}
		if err != nil {
			return
		}
		bodybytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		bodybytes = EncDec2(bodybytes, 88)
		err = json.Unmarshal(bodybytes, &r)
		return
	}
	jsonBytes, _ := json.Marshal(p)
	jsonBytes = EncDec2(jsonBytes, 88)
	r, err = post(url, jsonBytes)
	return
}
