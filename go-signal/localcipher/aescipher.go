package localcipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

// 生成指定长度的随机 IV（nonce）
// GCM 推荐使用 12 字节（96位）的 nonce
func generateNonce(length int) ([]byte, error) {
	nonce := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, nonce)
	return nonce, err
}

// AesDecrypt 解密
func AesGCMEncrypt(plaintext []byte, key []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//使用gcm
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	// 3. 生成或验证 nonce（IV）
	nonceSize := aesgcm.NonceSize() // GCM 推荐的 nonce 长度（12字节）
	// 自动生成随机 nonce
	nonce, err := generateNonce(nonceSize)
	if err != nil {
		return nil, fmt.Errorf("generate nonce faild: %v", err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	tag := RandomString(3)
	fullData := append([]byte(tag), nonce...)
	fullData = append(fullData, ciphertext...)
	return fullData, nil
}

// AesGCMDecrypt 解密
func AesGCMDecrypt(data []byte, key []byte) ([]byte, error) {
	iv := data[3:15]

	ciphertext := data[15:]
	// aad:=data[79:]
	//创建实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//使用gcm
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return aesgcm.Open(nil, iv, ciphertext, nil)
}

func AesCBCEncrypt(plaintext []byte, key []byte) (iv []byte, ciphertext []byte, err error) {
	// 1. 创建 AES 加密块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	// 2. 生成随机 IV（长度 = 块大小 = 16 字节）
	iv = make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, nil, err
	}

	// 3. 对明文进行 PKCS#7 填充（确保长度为块大小的整数倍）
	plaintext = pkcs7Pad(plaintext, aes.BlockSize)

	// 4. 初始化 CBC 加密器
	mode := cipher.NewCBCEncrypter(block, iv)

	// 5. 执行加密
	ciphertext = make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)

	// 6. 拼接 IV 和密文，再进行 base64 编码（方便传输）

	return iv, ciphertext, nil
}

// AES-CBC 解密
// ciphertextBase64: 加密后的 base64 字符串（格式：IV + 密文）
// key: 密钥（与加密时相同）
// 返回值：解密后的明文
func AesCBCDecrypt(iv []byte, ciphertext []byte, key []byte) ([]byte, error) {

	// 2. 验证数据长度（至少包含 IV）
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("密文长度不足")
	}

	// 4. 创建 AES 解密块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 5. 初始化 CBC 解密器
	mode := cipher.NewCBCDecrypter(block, iv)

	// 6. 执行解密
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// 7. 去除 PKCS#7 填充
	return pkcs7Unpad(plaintext)
}

// PKCS#7 填充（确保数据长度为 blockSize 的整数倍）
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(data, padtext...)
}

// PKCS#7 去除填充
func pkcs7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("数据长度为 0")
	}
	padding := int(data[len(data)-1])
	if padding > len(data) {
		return nil, errors.New("无效的填充长度")
	}
	return data[:len(data)-padding], nil
}
