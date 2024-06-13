// @Author Eric
// @Date 2024/6/9 23:30:00
// @Desc 加解密相关的处理
package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"hash/crc32"
)

var iv = []byte{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2} // IV should be 16 bytes for AES-256

func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

func pkcs7Unpadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("invalid padding size")
	}
	padding := int(data[length-1])
	if padding > length {
		return nil, errors.New("padding > length")
	}
	return data[:(length - padding)], nil
}

// Aes256Encrypt
//
//	@Description: AES256加密，CBC模式，PKCS7Padding，iv固定
//	@param key 加密key
//	@param plaintext 原始数据
//	@return string 返回base64之后数据
//	@return error
func Aes256Encrypt(key, plaintext []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plaintext = pkcs7Padding(plaintext, block.BlockSize())
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	//iv := ciphertext[:aes.BlockSize]

	//if _, err := io.ReadFull(rand.Reader, iv); err != nil {
	//	return "", err
	//}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Aes256Decrypt
//
//	@Description: AES256解密,CBC模式，PKCS7Padding，iv固定
//	@param key 加密key
//	@param ciphertext base64的密文
//	@return []byte 解密后数据
//	@return error
func Aes256Decrypt(key []byte, ciphertext string) ([]byte, error) {
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	fmt.Print(ciphertextBytes)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertextBytes) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	//iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	fmt.Print(ciphertextBytes)

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertextBytes, ciphertextBytes)

	fmt.Print(ciphertextBytes)

	return pkcs7Unpadding(ciphertextBytes)
}

// GenerateCRC32
//
//	@Description: 生成CRC32校验码
//	@param data
//	@return string
func GenerateCRC32(data []byte) string {
	return fmt.Sprintf("%08x", crc32.ChecksumIEEE(data))
}
