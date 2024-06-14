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

// PKCS7 padding
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// PKCS7 unpadding
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

// Aes256EncryptBase64
//
//	@Description: AES256加密，CBC模式，PKCS7Padding，iv就是密钥
//	@param key 加密key
//	@param plaintext 原始数据
//	@return string 返回base64之后数据
//	@return error
func Aes256EncryptBase64(key []byte, plaintext string) (string, error) {
	if len(key) != 32 {
		return "", errors.New("key length must be 32 bytes")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	paddedData := pkcs7Padding([]byte(plaintext), aes.BlockSize)
	ciphertext := make([]byte, len(paddedData))

	mode := cipher.NewCBCEncrypter(block, key[:aes.BlockSize])
	mode.CryptBlocks(ciphertext, paddedData)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Aes256DecryptBase64
//
//	@Description: AES256解密,CBC模式，PKCS7Padding，iv就是密钥
//	@param key 加密key
//	@param ciphertext base64的密文
//	@return []byte 解密后数据
//	@return error
func Aes256DecryptBase64(key []byte, ciphertext string) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("key length must be 32 bytes")
	}
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertextBytes) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	plaintext := make([]byte, len(ciphertextBytes))

	mode := cipher.NewCBCDecrypter(block, key[:aes.BlockSize])
	mode.CryptBlocks(plaintext, ciphertextBytes)

	plaintext, err = pkcs7Unpadding(plaintext)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// Aes256Encrypt
//
//	@Description: AES256加密，CBC模式，PKCS7Padding，iv就是密钥
//	@param key 加密key
//	@param plaintext 原始数据
//	@return byte 返回加密后数据
//	@return error
func Aes256Encrypt(key []byte, plaintext []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("key length must be 32 bytes")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	paddedData := pkcs7Padding(plaintext, aes.BlockSize)
	ciphertext := make([]byte, len(paddedData))

	mode := cipher.NewCBCEncrypter(block, key[:aes.BlockSize])
	mode.CryptBlocks(ciphertext, paddedData)

	return ciphertext, nil
}

// Aes256Decrypt
//
//	@Description: AES256解密,CBC模式，PKCS7Padding，iv就是密钥
//	@param key 加密key
//	@param ciphertext 密文
//	@return []byte 解密后数据
//	@return error
func Aes256Decrypt(key []byte, cipherBytes []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("key length must be 32 bytes")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(cipherBytes) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	plaintext := make([]byte, len(cipherBytes))

	mode := cipher.NewCBCDecrypter(block, key[:aes.BlockSize])
	mode.CryptBlocks(plaintext, cipherBytes)

	plaintext, err = pkcs7Unpadding(plaintext)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// GenerateCRC32
//
//	@Description: 生成CRC32校验码
//	@param data
//	@return string
func GenerateCRC32(data []byte) string {
	return fmt.Sprintf("%08x", crc32.ChecksumIEEE(data))
}
