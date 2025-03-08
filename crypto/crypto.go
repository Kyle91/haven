// @Author Eric
// @Date 2024/6/9 23:30:00
// @Desc 加解密相关的处理
package crypto

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"math/rand"
	"time"
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

// MD5
//
//	@Description: 生成MD5校验码
//	@param data
//	@return string
func MD5(data []byte) string {
	hash := md5.New()
	hash.Write(data)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// SHA256
//
//	@Description: 生成SHA256校验码
//	@param data
//	@return string
func SHA256(data []byte) string {
	hash := sha256.New()
	hash.Write(data)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// HMACSHA256
//
//	@Description: 生成基于密钥的 HMAC-SHA256 签名
//	@param key 签名密钥
//	@param message 签名内容
//	@return string 返回签名的十六进制字符串
//	@return error
func HMACSHA256(key []byte, message string) (string, error) {
	if len(key) == 0 {
		return "", errors.New("key must not be empty")
	}

	mac := hmac.New(sha256.New, key)
	_, err := mac.Write([]byte(message))
	if err != nil {
		return "", err
	}

	signature := mac.Sum(nil)
	return fmt.Sprintf("%x", signature), nil
}

func long2str(v []uint64, w bool) string {
	lenV := len(v)
	n := (lenV - 1) << 3 // n 是以字节为单位，表示长度

	if w {
		m := v[lenV-1]
		if m < uint64(n-7) || m > uint64(n) {
			return ""
		}
		n = int(m)
	}

	var buf bytes.Buffer
	for i := 0; i < lenV; i++ {
		// 将 uint64 数字转换为字节流
		var b [8]byte
		binary.LittleEndian.PutUint64(b[:], v[i])
		buf.Write(b[:])
	}

	if w {
		return buf.String()[:n] // 如果 w 为 true，返回前 n 个字节
	}

	return buf.String()
}

func str2long(s string, w bool) []uint64 {
	var v []uint64
	// Pad string to be a multiple of 8 bytes (64 bits)
	padded := s + string(make([]byte, (8-len(s)%8)%8))
	for i := 0; i < len(padded); i += 8 {
		var n uint64
		// 将每 8 个字节转为 uint64
		binary.Read(bytes.NewReader([]byte(padded[i:i+8])), binary.LittleEndian, &n)
		v = append(v, n)
	}
	if w {
		v = append(v, uint64(len(s)))
	}
	return v
}

// XXTeaEncrypt
//
//	@Description: 带过期时间的加密
//	@param str
//	@param key
//	@param expiry
//	@return string
func XXTeaEncrypt(str, key string, expiry int) string {
	if str == "" {
		return ""
	}
	ckeyLength := 8
	str += random(ckeyLength)

	// Expiry time
	if expiry != 0 {
		str = fmt.Sprintf("%010d", expiry+int(time.Now().Unix())) + str
	} else {
		str = fmt.Sprintf("%010d", 0) + str
	}

	v := str2long(str, true)
	k := str2long(key, false)

	// Ensure key length is at least 4
	for len(k) < 4 {
		k = append(k, 0)
	}

	n := len(v) - 1
	z := v[n]
	y := v[0]
	delta := uint64(0x9E3779B9)
	q := uint64(6 + 52/(n+1))
	sum := uint64(0)

	// Start the encryption loop
	for q > 0 {
		q--
		// Update sum, keep it as uint64
		sum = sum + delta
		e := sum >> 2 & 3
		for p := 0; p < n; p++ {
			y = v[p+1]
			// Ensure p is uint64 for bitwise operation with e
			mx := ((z >> 5 & 0x07ffffff) ^ y<<2) + ((y >> 3 & 0x1fffffff) ^ z<<4) ^ (sum ^ y) + (k[p&3^int(e)] ^ z)
			v[p] = v[p] + mx // First, update v[p]
			z = v[p]         // Then update z
		}
		y = v[0]
		mx := ((z >> 5 & 0x07ffffff) ^ y<<2) + ((y >> 3 & 0x1fffffff) ^ z<<4) ^ (sum ^ y) + (k[n&3^int(e)] ^ z)
		v[n] = v[n] + mx // Update v[n]
		z = v[n]         // Then update z
	}

	return long2str(v, false)
}

// XXTeaDecrypt
//
//	@Description: 带过期时间的解密
//	@param str
//	@param key
//	@return string
func XXTeaDecrypt(str, key string) string {
	if str == "" {
		return ""
	}
	ckeyLength := 8

	v := str2long(str, false)
	k := str2long(key, false)

	// Ensure key length is at least 4
	for len(k) < 4 {
		k = append(k, 0)
	}

	n := len(v) - 1
	y := v[0]
	delta := uint64(0x9E3779B9)
	q := uint64(6 + 52/(n+1))
	sum := uint64(q * delta)

	// Decryption loop
	for sum != 0 {
		e := sum >> 2 & 3
		for p := n; p > 0; p-- {
			z := v[p-1]
			// Ensure p is uint64 for bitwise operation with e
			mx := ((z >> 5 & 0x07ffffff) ^ y<<2) + ((y >> 3 & 0x1fffffff) ^ z<<4) ^ (sum ^ y) + (k[p&3^int(e)] ^ z)
			// Separate assignment to avoid chain assignment
			v[p] = v[p] - mx
			y = v[p] // Then update y
		}
		z := v[n]
		mx := ((z >> 5 & 0x07ffffff) ^ y<<2) + ((y >> 3 & 0x1fffffff) ^ z<<4) ^ (sum ^ y) + (k[n&3^int(e)] ^ z)
		v[0] = v[0] - mx  // Update v[0]
		y = v[0]          // Then update y
		sum = sum - delta // Update sum
	}

	// Convert result to string
	ret := long2str(v, true)
	dateLen := len(ret)
	ret = ret[:dateLen-ckeyLength]

	// Check if the expiry time is valid
	if ret[:10] == "0000000000" || (time.Now().Unix()-int64(toInt(ret[:10]))) > 0 {
		ret = ret[10:]
	} else {
		ret = ""
	}

	return ret
}

func random(length int) string {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz"
	var randStr []byte
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		randStr = append(randStr, chars[rand.Intn(len(chars))])
	}
	return string(randStr)
}

func toInt(str string) int64 {
	var result int64
	for i := 0; i < len(str); i++ {
		result = result*10 + int64(str[i]-'0')
	}
	return result
}

// RSAEncrypt
//
// @Description: RSA公钥加密
// @param pubKey 公钥
// @param plaintext 原始数据
// @return string 返回加密后的数据
// @return error
func RSAEncrypt(pubKey *rsa.PublicKey, plaintext []byte) (string, error) {
	ciphertext, err := rsa.EncryptPKCS1v15(cryptoRand.Reader, pubKey, plaintext)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// RSADecrypt
//
// @Description: RSA私钥解密
// @param privKey 私钥
// @param ciphertext base64加密的密文
// @return []byte 返回解密后的数据
// @return error
func RSADecrypt(privKey *rsa.PrivateKey, ciphertext string) ([]byte, error) {
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}
	plaintext, err := rsa.DecryptPKCS1v15(cryptoRand.Reader, privKey, ciphertextBytes)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// AesGCMEncrypt
//
//	@Description: AES-GCM 加密
//	@param key 加密key
//	@param plaintext 原始数据
//	@return []byte 返回加密后的数据（包括nonce和ciphertext）
//	@return error
func AesGCMEncrypt(key []byte, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := cryptoRand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	return append(nonce, ciphertext...), nil
}

// AesGCMDecrypt
//
//	@Description: AES-GCM 解密
//	@param key 加密key
//	@param ciphertext 密文（包括nonce和ciphertext）
//	@return []byte 返回解密后的数据
//	@return error
func AesGCMDecrypt(key []byte, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// GenerateRandomKey
//
//	@Description: 生成指定长度的随机密钥，用于加密或签名等场景
//	@param length 密钥的长度（字节数）
//	@return []byte 返回生成的随机密钥
//	@return error 如果生成失败，返回错误
func GenerateRandomKey(length int) ([]byte, error) {
	key := make([]byte, length)
	_, err := cryptoRand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func ZlibCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := zlib.NewWriter(&buf)
	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}
	writer.Close()
	return buf.Bytes(), nil
}

func ZlibDecompress(data []byte) ([]byte, error) {
	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(reader)
}
