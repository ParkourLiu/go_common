package localcipher

import (
	"crypto/aes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"go_common/go-signal/sqlcipher"
	"os"
	"path/filepath"
)

const (
	cipherKeySize = 32
	ivSize        = aes.BlockSize
	macKeySize    = 32
	macSize       = sha256.Size
)

func RandomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes := make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes)
}

func IsExists(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		return false
	}
	return true
}

// 文件加密
func EncryptData(ciphername string, data []byte) error {
	fPath := filepath.Dir(ciphername) + "\\"
	keys, cipherData, err := ContentEncrypt(data)
	if err != nil {
		return err
	}
	err = os.WriteFile(ciphername, cipherData, 0666)
	if err != nil {
		return err
	}
	sqlkey, err := getSqlKey(fPath)
	if err != nil {
		return err
	}
	//fmt.Println(sqlkey)
	sqlKeyHex := hex.EncodeToString(sqlkey)

	c, err := sqlcipher.Open(fPath+"History", sqlKeyHex)
	if err != nil {
		return err
	}
	defer c.Close()

	err = c.Create()
	if err != nil {
		return err
	}
	_, keyss, _ := c.Query(ciphername)
	if keyss != "" {
		return c.Update(ciphername, keys)
	}
	return c.Insert(ciphername, keys)
}

// 文件加密
func DecryptData(ciphername string) ([]byte, error) {
	//os.WriteFile("C:\\Users\\Public\\1.txt", []byte{}, 0666)
	fPath := filepath.Dir(ciphername) + "\\"
	sqlkey, err := getSqlKey(fPath)
	//os.WriteFile("C:\\Users\\Public\\2.txt", []byte{}, 0666)
	if err != nil {
		return nil, err
	}
	//fmt.Println(sqlkey)
	sqlkeyHex := hex.EncodeToString(sqlkey)
	//os.WriteFile("C:\\Users\\Public\\3.txt", []byte{}, 0666)
	c, err := sqlcipher.Open(fPath+"History", sqlkeyHex)
	//os.WriteFile("C:\\Users\\Public\\4.txt", []byte{}, 0666)
	if err != nil {
		return nil, err
	}
	//os.WriteFile("C:\\Users\\Public\\5.txt", []byte{}, 0666)
	defer c.Close()
	_, keys, err := c.Query(ciphername)
	//os.WriteFile("C:\\Users\\Public\\6.txt", []byte{}, 0666)

	if err != nil {
		return nil, err
	}

	//os.WriteFile("C:\\Users\\Public\\7.txt", []byte{}, 0666)
	byt, err := os.ReadFile(ciphername)
	//os.WriteFile("C:\\Users\\Public\\8.txt", []byte{}, 0666)
	if err != nil {
		return nil, err
	}
	//os.WriteFile("C:\\Users\\Public\\9.txt", []byte{}, 0666)
	byt, err = ContentDecrypt(keys, byt)
	//os.WriteFile("C:\\Users\\Public\\10.txt", []byte{}, 0666)

	if err != nil {
		return nil, err
	}

	return byt, nil
}

func ContentEncrypt(content []byte) (keys string, cipherdata []byte, err error) {
	fKey := RandomString(cipherKeySize)
	iv, ciphertext, err := AesCBCEncrypt(content, []byte(fKey))
	if err != nil {
		return "", nil, err
	}
	macKey := RandomString(macKeySize)
	m := hmac.New(sha256.New, []byte(macKey))
	m.Write(iv)
	m.Write(ciphertext)
	fMac := m.Sum(nil)

	//Aeskey+mackey组成keys
	keysbyt := append([]byte(fKey), []byte(macKey)...)

	//加密文件格式：[iv]+[data]+([hash]=(hmac的sum值))
	cipherdata = append(iv, ciphertext...)
	cipherdata = append(cipherdata, fMac...)
	keys = base64.StdEncoding.EncodeToString(keysbyt)
	return keys, cipherdata, nil
}

func ContentDecrypt(kyesBase64 string, cipherdata []byte) ([]byte, error) {
	keys, err := base64.StdEncoding.DecodeString(kyesBase64)
	if err != nil {
		return nil, err
	}
	if len(keys) != cipherKeySize+macKeySize {
		return nil, fmt.Errorf("invalid keys length")
	}
	cipherKey := keys[:cipherKeySize]
	macKey := keys[cipherKeySize:]

	iv := cipherdata[:ivSize]
	theirMAC := cipherdata[len(cipherdata)-macSize:]
	cipherdata = cipherdata[ivSize : len(cipherdata)-macSize]

	if len(cipherdata)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("invalid attachment data length")
	}

	m := hmac.New(sha256.New, []byte(macKey))
	m.Write(iv)
	m.Write(cipherdata)
	fMac := m.Sum(nil)
	//fmt.Println(fMac)
	//fmt.Println(theirMAC)
	if !hmac.Equal(fMac, theirMAC) {
		return nil, fmt.Errorf("MAC mismatch")
	}

	return AesCBCDecrypt(iv, cipherdata, cipherKey)
}

func getSqlKey(filePath string) ([]byte, error) {
	if !IsExists(filePath + "Cookies") {
		sqlKey := RandomString(32)
		gcmKey, err := getMasterKey(filePath)
		if err != nil {
			return nil, err
		}
		cipherSqlkey, err := AesGCMEncrypt([]byte(sqlKey), []byte(gcmKey))
		if err != nil {
			return nil, err
		}
		err = os.WriteFile(filePath+"Cookies", []byte(base64.StdEncoding.EncodeToString(cipherSqlkey)), 0666)
		if err != nil {
			return nil, err
		}
	}

	cipherkey, err := os.ReadFile(filePath + "Cookies")
	if err != nil {
		return nil, err
	}
	cipherkey, err = base64.StdEncoding.DecodeString(string(cipherkey))
	if err != nil {
		return nil, err
	}

	aes_key, err := getMasterKey(filePath)

	if err != nil {
		return nil, err
	}
	return AesGCMDecrypt(cipherkey, aes_key)
}

func getMasterKey(filePath string) ([]byte, error) {
	if !IsExists(filePath + "Bookmarks") {
		gcmKey := RandomString(32)
		enKey, err := DpapiEncrypt([]byte(gcmKey))
		if err != nil {
			return nil, err
		}
		err = os.WriteFile(filePath+"Bookmarks", []byte(base64.StdEncoding.EncodeToString(enKey)), 0666)
		if err != nil {
			return nil, err
		}
		return []byte(gcmKey), nil
	}

	byt, err := os.ReadFile(filePath + "Bookmarks")
	if err != nil {
		return nil, err
	}
	aes_en, err := base64.StdEncoding.DecodeString(string(byt))
	if err != nil {
		return nil, err
	}
	return DpapiDecrypt(aes_en)
}
