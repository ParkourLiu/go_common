package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"go_common/go-signal/localcipher"
	"go_common/go-signal/sqlcipher"
	"log"
	"os"
)

func main() {
	encryptTest()
}
func encryptTest() {
	ke := `aaaa`
	err := localcipher.EncryptData(`E:\gowork\src\go_common\go-signal\data\aaa`, []byte(ke))
	if err != nil {
		log.Fatalln(1, err)
	}
	jbs, err := localcipher.DecryptData(`E:\gowork\src\go_common\go-signal\data\aaa`)
	if err != nil {
		log.Fatalln(2, err)
	}
	fmt.Println(string(jbs))
}

func decryptTest() {
	byt, err := localcipher.DecryptData("keke")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("keke")
	fmt.Println(string(byt))

	byt, err = localcipher.DecryptData("k1k1")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("k1k1")
	fmt.Println(string(byt))
}

// 文件加密
func EncryptData(ciphername string, data []byte) error {

	keys, cipherData, err := localcipher.ContentEncrypt(data)
	if err != nil {
		return err
	}
	err = os.WriteFile(ciphername, cipherData, 0666)
	if err != nil {
		return err
	}
	sqlKey := localcipher.RandomString(32)
	fmt.Println("sqlkey ", sqlKey)

	sqlKeyHex := hex.EncodeToString([]byte(sqlKey))
	c, err := sqlcipher.Open("./localconfig", sqlKeyHex)
	if err != nil {
		return err
	}
	defer c.Close()

	err = c.Create()
	if err != nil {
		return err
	}
	err = c.Insert(ciphername, keys)
	if err != nil {
		return err
	}
	gcmKey := localcipher.RandomString(32)
	fmt.Println("gcmKey ", gcmKey)

	cipherkey, err := localcipher.AesGCMEncrypt([]byte(sqlKey), []byte(gcmKey))
	if err != nil {
		return err
	}
	enKey, err := localcipher.DpapiEncrypt([]byte(gcmKey))
	if err != nil {
		return err
	}
	os.WriteFile("config", []byte(base64.StdEncoding.EncodeToString(enKey)), 0666)

	os.WriteFile("local.ini", []byte(base64.StdEncoding.EncodeToString(cipherkey)), 0666)

	return nil
}

// 文件解密
func DecryptData(ciphername string) ([]byte, error) {
	byt, err := os.ReadFile("config")
	if err != nil {
		return nil, err
	}
	aes_en, err := base64.StdEncoding.DecodeString(string(byt))
	if err != nil {
		return nil, err
	}
	aes_key, err := localcipher.DpapiDecrypt(aes_en)
	if err != nil {
		return nil, err
	}
	fmt.Println("gcmKey ", string(aes_key))

	cipherkey, err := os.ReadFile("local.ini")
	if err != nil {
		return nil, err
	}
	cipherkey, err = base64.StdEncoding.DecodeString(string(cipherkey))

	if err != nil {
		return nil, err
	}
	sqlkey, err := localcipher.AesGCMDecrypt(cipherkey, aes_key)
	if err != nil {
		return nil, err
	}
	fmt.Println("sqlkey ", string(sqlkey))

	sqlkeyHex := hex.EncodeToString(sqlkey)
	c, err := sqlcipher.Open("localconfig", sqlkeyHex)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	_, keys, err := c.Query(ciphername)

	if err != nil {
		return nil, err
	}
	byt, err = os.ReadFile(ciphername)
	if err != nil {
		return nil, err
	}
	byt, err = localcipher.ContentDecrypt(keys, byt)

	if err != nil {
		return nil, err
	}

	return byt, nil
}
