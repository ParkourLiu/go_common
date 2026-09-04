package cregcenter

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	code200 = 200
	code301 = 301
)

type returnData struct {
	Code    int          `json:"code,omitempty"`
	Msg     string       `json:"msg,omitempty"`
	Token   string       `json:"token,omitempty"`
	Servers []*paramData `json:"servers,omitempty"`
}

func (r *returnData) Status(code int, msg string) {
	r.Code = code
	r.Msg = msg
}

type paramData struct {
	LoginName     string `json:"name,omitempty"`
	LoginPassWord string `json:"pass,omitempty"`
	Workspace     string `json:"workspace,omitempty"`
	ServerName    string `json:"serverName,omitempty"`
	Port          string `json:"port,omitempty"`

	Host    string `json:"ipPort,omitempty"`
	Healthy int    `json:"healthy,omitempty"`
}

type RegInfo struct {
	LoginName     string
	LoginPassword string
	RegCenterAddr string
	Workspace     string
	ServerName    string
	Port          string
	callAddr      map[string][]string
}

type CallServerClient struct {
	l          *sync.RWMutex
	serverName string
	hosts      []string
}

func (c *CallServerClient) Call(method, url string, headMap map[string]string, bodybytes []byte, timeoutSecond int) ([]byte, error) {
	l := len(c.hosts)
	if l == 0 {
		return nil, fmt.Errorf("%s service does not exist", c.serverName)
	}
	i := rand.Intn(l)
	var send = func(method, url string, headMap map[string]string, bodybytes []byte, timeoutSecond int) ([]byte, error) {
		body := bytes.NewReader(bodybytes)
		request, err := http.NewRequest(method, fmt.Sprintf("http://%s%s", c.hosts[i], url), body)
		if err != nil {
			return nil, err
		}
		if request != nil {
			if request.Body != nil {
				defer func() {
					_ = request.Body.Close()
					request.Close = true
				}()
			}
		}
		for k, v := range headMap {
			request.Header[k] = []string{v}
		}
		httpClient := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
					conn, err := net.Dial("tcp", addr)
					if err != nil {
						return nil, err
					}
					_ = conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeoutSecond)))
					return conn, nil
				},
			},
			Timeout: time.Second * time.Duration(timeoutSecond),
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
			return nil, err
		}

		return ioutil.ReadAll(resp.Body)
	}
	return send(method, url, headMap, bodybytes, timeoutSecond)
}
func (c *CallServerClient) callFindServer() {
	for {
		c.findServer()
		time.Sleep(time.Minute)
	}
}
func (c *CallServerClient) findServer() (err error) {
	servers := getServer(c.serverName)
	if servers == nil || len(servers) == 0 {
		return fmt.Errorf("%s service does not exist", c.serverName)
	}
	hosts := make([]string, len(servers))
	for i, s := range servers {
		hosts[i] = s.Host
	}
	c.l.Lock()
	defer c.l.Unlock()
	c.hosts = hosts
	return
}
