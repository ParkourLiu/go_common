package cmethod_test

import (
	"fmt"
	"go_common/cmethod"
	"testing"
)

func TestHTTPrequest(t *testing.T) {
	for {
		b, err := cmethod.HTTPrequest("GET", "http://www.baidu.com", nil, nil, 3)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(len(b))
	}

}

func TestGenRsaKey(t *testing.T) {
	cmethod.CreatRsaKey()
}

var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDRIiAaFZrBFEzNlaGLqFPiqCP4
7CrR7wYphraHOevWWUSFzxCUQm8+n/pNUGBkTRtMbOixgVvmAqIzv3zXqmK3/mGB
XhCbw+D7NILPBMwu0cy7O2eJfcCKTzonEg6nnPEytCZV/L0ht9Z7XtOXulQ5VrwC
p7rX1UOIWMDW6A4kRQIDAQAB
-----END PUBLIC KEY-----
`)

var privKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBANEiIBoVmsEUTM2V
oYuoU+KoI/jsKtHvBimGtoc569ZZRIXPEJRCbz6f+k1QYGRNG0xs6LGBW+YCojO/
fNeqYrf+YYFeEJvD4Ps0gs8EzC7RzLs7Z4l9wIpPOicSDqec8TK0JlX8vSG31nte
05e6VDlWvAKnutfVQ4hYwNboDiRFAgMBAAECgYEAwDu/CFsNiicvxdWhza7nlLN7
hWcIoTo2Dtu+UiSirMAXZWwFUFKU0RraSFD2mZvq2OBPMEK5B38qO6jrh44d+Fqp
A+10azawc5mrMIVbKiSxMaVDJcLHwn85uZ/za/S2AnMHhzC6crQt3qsKC/S4+aKH
O2DljAWU/XR2aFfrRYECQQDaRxJqUckD24T9R3fOjOrgb300ysbASVqT0f5klS10
ZSVb8wYSqGpBfqIAGOOxed4Qegx8SnAv/pTDnzIxPsphAkEA9UZ+t1TnR4ivEhL/
mZHvWMFoejTe5fMDqNf7wsGdT3x9SqCYCIBGwHU8V/K634e4ur1cn2gbElBfzzLR
HePMZQJAe1GgA9VE/hrtnbLc6yMOJ9KVKFhPxZ8rv0vqr6TgU1w5qSM6ERx5O5tx
pyBos4IohaKOn0Hm9BaesY6latEQwQJAY1p56+NtiBF68TRW6ystK+O0YYRXIght
XBCZP8vT4CXKTtd8njzv6/fRSMLfJbrBfotEIKI4DRQXq0OnZ5cl7QJBAIWt7x4l
je1FqyDWaUP1vCkJ9q9pzoYEPr7aQ/Wg/SHHIkT8Sg3gO/WyrwJHxmGIm+3E69ve
ECRWgO8jiXUewW0=
-----END RSA PRIVATE KEY-----
`)

func Test_RSA(t *testing.T){
		str := "草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马草泥马"
		strByt := []byte(str)
		encStr, err := cmethod.RsaSubEnc(strByt, publicKey)
		if err != nil {
			fmt.Println(111, err)
			return
		}
		fmt.Println(encStr)
		b, err := cmethod.RsaSubDec(encStr, privKey)
		if err != nil {
			fmt.Println(222, err)
			return
		}
		fmt.Println(string(b))
		enc, err := cmethod.RsaEnc(strByt, publicKey)
		if err != nil {
			fmt.Println(111, err)
			return
		}
		fmt.Println(cmethod.RsaDec(enc, privKey))

}