package cmethod

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//gzip压缩
func GzipBytes(byt []byte) []byte {
	var buf bytes.Buffer
	//zw := gzip.NewWriter(&buf)
	zw, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)

	zw.Write(byt)
	if err := zw.Close(); err != nil {
	}
	return buf.Bytes()
}

//gzip解压缩
func UnGzipBytes(byt []byte) ([]byte, error) {
	if len(byt) == 0 {
		return nil, nil
	}
	var buf bytes.Buffer
	buf.Write(byt)
	zr, err := gzip.NewReader(&buf)
	defer func() {
		if zr != nil {
			zr.Close()
		}
	}()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(zr)
}

func EncDec(byt []byte, divisor int) []byte {
	for i, v := range byt {
		byt[i] = v ^ byte(i*5+divisor)
	}
	return byt
}

func EncDec2(byt []byte, divisor int) []byte {
	for i, v := range byt {
		byt[i] = (byte(i+divisor) & (^v)) | (v & (^byte(i + divisor)))
	}
	return byt
}

func Enc3(byt, pKeyTemp []byte, divisor int) []byte {
	for i, _ := range byt {
		byt[i] = (byte(i+divisor) & (^byt[i])) | (byt[i] & (^byte(i + divisor)))
		byt[i] = byt[i] ^ pKeyTemp[i%len(pKeyTemp)]
	}
	return byt
}
func Dec3(byt, pKeyTemp []byte, divisor int) []byte {
	for i, _ := range byt {
		byt[i] = byt[i] ^ pKeyTemp[i%len(pKeyTemp)]
		byt[i] = (byte(i+divisor) & (^byt[i])) | (byt[i] & (^byte(i + divisor)))
	}
	return byt
}

func WriteFile(byt []byte, path, name string) error {
	err := CreateDir(path)
	if err != nil {
		return err
	}

	tmpFileName := path + name + ".tmp"
	f, err := os.Create(tmpFileName)
	if err != nil {
		return err
	}
	_, err = f.Write(byt)
	f.Close()
	if err != nil {
		return err
	}
	return os.Rename(tmpFileName, path+name)
}

//调用os.MkdirAll递归创建文件夹,只会创建文件夹
func CreateDir(filePath string) error {
	if !IsExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		return err
	}
	return nil
}

// 判断所给路径文件/文件夹是否存在(返回true是存在)
func IsExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

//获取md5值
func Byt2md5(byt []byte) string {
	h := md5.New()
	h.Write(byt)
	return hex.EncodeToString(h.Sum(nil))
}

//http请求
func HTTPrequest(method, url string, headMap map[string]string, bodybytes []byte, timeoutSecond int) ([]byte, error) {
	body := bytes.NewReader(bodybytes)
	request, err := http.NewRequest(method, url, body) //创建请求体
	if err != nil {
		return nil, err
	}
	//添加头信息
	for k, v := range headMap {
		request.Header[k] = []string{v}
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				conn, err := net.Dial("tcp", addr)
				if err != nil {
					return nil, err
				}
				_ = conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeoutSecond))) //设置 发送+接受  数据超时时间
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

//http请求
func HTTPSrequest(method, url string, headMap map[string]string, bodybytes []byte, timeoutSecond int) ([]byte, error) {
	body := bytes.NewReader(bodybytes)
	request, err := http.NewRequest(method, url, body) //创建请求体
	if err != nil {
		return nil, err
	}
	//添加头信息
	for k, v := range headMap {
		request.Header[k] = []string{v}
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //设置tls.Config的InsecureSkipVerify为true，client将不再对服务端的证书进行校验
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				conn, err := net.Dial("tcp", addr)
				if err != nil {
					return nil, err
				}
				_ = conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeoutSecond))) //设置 发送+接受  数据超时时间
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

//func CreatRsaKey() error {
//	// 生成私钥文件
//	privateKey, err := rsa.GenerateKey(rand.Reader, 1024) //默认1024长度
//	if err != nil {
//		return err
//	}
//	derStream, _ := x509.MarshalPKCS8PrivateKey(privateKey)
//	//derStream := x509.MarshalPKCS1PrivateKey(privateKey)
//	priBlock := &pem.Block{
//		Type:  "RSA PRIVATE KEY",
//		Bytes: derStream,
//	}
//
//	fmt.Printf("=======私钥文件内容=========\n%v", string(pem.EncodeToMemory(priBlock)))
//	// 生成公钥文件
//	publicKey := &privateKey.PublicKey
//	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
//	if err != nil {
//		return err
//	}
//	publicBlock := &pem.Block{
//		Type:  "PUBLIC KEY",
//		Bytes: derPkix,
//	}
//
//	fmt.Printf("=======公钥文件内容=========\n%v", string(pem.EncodeToMemory(publicBlock)))
//
//	if err != nil {
//		return err
//	}
//	return nil
//}

//rsa分段加密   密钥长度，默认为1024位  超过不适用
func RsaSubEnc(dataBytes []byte, pubKey []byte) ([]byte, error) {
	strBytLenth := len(dataBytes)
	encByt := []byte{}
	encLenth := 117
	for i := 0; i < strBytLenth; i = i + encLenth { //一次只能加密117个字节
		var data []byte
		var err error
		if i+encLenth <= strBytLenth {
			data, err = RsaEnc(dataBytes[i:i+encLenth], pubKey)
		} else {
			data, err = RsaEnc(dataBytes[i:], pubKey)
		}
		if err != nil {
			return nil, err
		}
		encByt = append(encByt, data...)
	}
	return encByt, nil
}

//rsa分段解密   密钥长度，默认为1024位  超过不适用
func RsaSubDec(dataBytes []byte, privKey []byte) ([]byte, error) {
	strBytLenth := len(dataBytes)
	decByt := []byte{}
	decLenth := 128
	for i := 0; i < strBytLenth; i = i + decLenth { //一次只能解密128个字节
		var data []byte
		var err error
		if i+decLenth <= strBytLenth {
			data, err = RsaDec(dataBytes[i:i+decLenth], privKey)
		} else {
			data, err = RsaDec(dataBytes[i:], privKey)
		}
		if err != nil {
			return nil, err
		}
		decByt = append(decByt, data...)
	}
	return decByt, nil
}

//公钥加密
func RsaEnc(origData []byte, pubKey []byte) ([]byte, error) {
	block, _ := pem.Decode(pubKey)
	if block == nil {
		return nil, errors.New("public key error!")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

//私钥解密
func RsaDec(origData []byte, privKey []byte) ([]byte, error) {
	block, _ := pem.Decode(privKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	privInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	priv := privInterface.(*rsa.PrivateKey)
	return rsa.DecryptPKCS1v15(rand.Reader, priv, origData)

}

//获取http的请求ip
func GetRemoteIp(r *http.Request) string {
	remoteAddr := r.RemoteAddr
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else if ip = r.Header.Get("X-Forwarded-For"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return remoteAddr
}

func Timestr2Unix(timestr, timeFormat string) (int64, error) { //timeFormat格式2006-01-02 15:04:05
	//时间 to 时间戳
	loc, _ := time.LoadLocation("Asia/Shanghai") //设置时区
	tt, err := time.ParseInLocation(timeFormat, timestr, loc)
	return tt.Unix() * 1000, err
}

//判断是ipv4 的ip且是公网ip
func checkIp(ip string) bool {
	//判断是否是ipv4
	regStr := `^(([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.)(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){2}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`
	if match, _ := regexp.MatchString(regStr, ip); !match {
		return false
	}

	//判断是不是公网(不在内网的ip段内)
	intIP := inet_aton(ip)
	if intIP == inet_aton("127.0.0.1") {

	} else if inet_aton("10.0.0.0") <= intIP && intIP <= inet_aton("10.255.255.255") {

	} else if inet_aton("172.16.0.0") <= intIP && intIP <= inet_aton("172.31.255.255") {

	} else if inet_aton("192.168.0.0") <= intIP && intIP <= inet_aton("192.168.255.255") {

	} else {
		return true
	}
	return false
}

//ip换算成数字
func inet_aton(ip string) int64 {
	bits := strings.Split(ip, ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64

	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}

//根据ifconfig获取本机公网地址,多地址用/拼接
func GetLocalIP() string {
	//Get List of interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	ip := ""
	//Take one interface from interfaces list
	for _, iface := range interfaces {
		//Get interface name
		ifaceName, err := net.InterfaceByName(iface.Name)
		if err != nil {
			return ""
		}
		//Get IPs assigned to the interface name
		cidrs, _ := ifaceName.Addrs()
		//Loop through each IP/netmask

		for _, cidr := range cidrs {
			ipNet, _, _ := net.ParseCIDR(cidr.String())
			tempIP := ipNet.String()
			if checkIp(tempIP) {
				ip = ip + tempIP + "/"
			}
		}

	}
	ip = strings.Trim(ip, "/")
	return ip
}
