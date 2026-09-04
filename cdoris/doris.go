package cdoris

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type DorisClient struct {
	host       string //10.0.0.214:8030
	auth       string
	httpClient *http.Client
}

func NewDorisClient(host, user, pass string, timeOut int) (dc *DorisClient) {
	dc = &DorisClient{
		host: host,
		auth: "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+pass)),
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
					conn, err := net.Dial("tcp", addr)
					if err != nil {
						return nil, err
					}
					_ = conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeOut))) //设置 发送+接受  数据超时时间
					return conn, nil
				},
			},
			Timeout: time.Second * time.Duration(timeOut),
		},
	}
	return
}

// fields====ct,domain,location,title,content
func (dc *DorisClient) PutCsv(dbName, tableName, fields, csvData string) (err error) {
	hMap := map[string]string{
		"group_commit":     "async_mode",
		"Authorization":    dc.auth,
		"Expect":           "100-continue",
		"column_separator": ",",  //默认逗号分隔符
		"enclose":          "'",  //包围符
		"escape":           "\\", //转义符
		"columns":          fields,
	}
	r, err := HTTPrequest(dc.httpClient, "PUT", fmt.Sprintf("http://%s/api/%s/%s/_stream_load", dc.host, dbName, tableName), hMap, []byte(csvData))
	if err != nil {
		return
	}
	if !bytes.Contains(r, []byte(`"Status": "Success"`)) {
		err = fmt.Errorf(string(r))
		return
	}
	return
}

// fields====ct,domain,location,title,content
func (dc *DorisClient) PutCsvByNative(dbName, tableName, fields string, csvData []byte) (err error) {
	hMap := map[string]string{
		"Authorization": dc.auth,
		"Expect":        "100-continue",
		"columns":       fields,
	}
	r, err := HTTPrequest(dc.httpClient, "PUT", fmt.Sprintf("http://%s/api/%s/%s/_stream_load", dc.host, dbName, tableName), hMap, csvData)
	if err != nil {
		return
	}
	if !bytes.Contains(r, []byte(`"Status": "Success"`)) {
		err = fmt.Errorf(string(r))
		return
	}
	return
}

func (dc *DorisClient) PutJsonNoUpdate(dbName, tableName string, jsonData []byte) (err error) {
	if len(jsonData) == 0 {
		return
	}
	hMap := map[string]string{
		"Authorization": dc.auth,
		"Expect":        "100-continue",
		"format":        "json",
		"group_commit":  "async_mode",
	}
	if jsonData[0] == '[' {
		hMap["strip_outer_array"] = "true"
	}
	r, err := HTTPrequest(dc.httpClient, "PUT", fmt.Sprintf("http://%s/api/%s/%s/_stream_load", dc.host, dbName, tableName), hMap, jsonData)
	if err != nil {
		return
	}
	if !bytes.Contains(r, []byte(`"Status": "Success"`)) {
		err = fmt.Errorf(string(r))
		return
	}
	return
}

func (dc *DorisClient) PutJson(dbName, tableName string, jsonData []byte) (err error) {
	if len(jsonData) == 0 {
		return
	}
	hMap := map[string]string{
		"Authorization":          dc.auth,
		"Expect":                 "100-continue",
		"format":                 "json",
		"group_commit":           "async_mode",
		"unique_key_update_mode": "UPDATE_FLEXIBLE_COLUMNS", //灵活列更新
	}
	if jsonData[0] == '[' {
		hMap["strip_outer_array"] = "true"
	}
	r, err := HTTPrequest(dc.httpClient, "PUT", fmt.Sprintf("http://%s/api/%s/%s/_stream_load", dc.host, dbName, tableName), hMap, jsonData)
	if err != nil {
		return
	}
	if !bytes.Contains(r, []byte(`"Status": "Success"`)) {
		err = fmt.Errorf(string(r))
		return
	}
	return
}

// http请求
func HTTPrequest(httpClient *http.Client, method, url string, headMap map[string]string, bodybytes []byte) ([]byte, error) {
	body := bytes.NewReader(bodybytes)
	request, err := http.NewRequest(method, url, body) //创建请求体
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
	//添加头信息
	for k, v := range headMap {
		request.Header[k] = []string{v}
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
