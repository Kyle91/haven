// @Author Eric
// @Date 2024/9/4 16:58:00
// @Desc
package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/Kyle91/haven/crypto"
	"io"
	"strconv"
	"strings"
	"time"
)

type AuthToken struct {
	SecretKey []byte //base64的
	Salt      string
}

// 初始化 AuthToken 类
// secretKey是base64的
func NewAuthToken(secretKey, salt string) *AuthToken {
	key, _ := base64.StdEncoding.DecodeString(secretKey)
	return &AuthToken{
		SecretKey: key,
		Salt:      salt,
	}
}

// 生成随机字符串
func generateRandomString() (string, error) {
	randBytes := make([]byte, 16) // 生成16字节的随机数
	if _, err := io.ReadFull(rand.Reader, randBytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(randBytes), nil
}

// 生成 Token
func (a *AuthToken) GenerateToken(userID int64, expirationTime int64) (string, error) {
	// 1. 生成过期时间戳
	expiration := time.Now().Add(time.Duration(expirationTime) * time.Second).Unix()

	// 2. 生成一个随机字符串
	randomString, err := generateRandomString()
	if err != nil {
		return "", err
	}

	// 3. 拼接原始数据: userID:expiration:salt:randomString
	rawData := fmt.Sprintf("%d:%d:%s:%s", userID, expiration, a.Salt, randomString)

	// 4. 对原始数据进行加密，生成 token
	token, err := crypto.Aes256Encrypt(a.SecretKey, []byte(rawData))
	if err != nil {
		return "", err
	}

	// 5. 将加密后的token转换为十六进制字符串
	return hex.EncodeToString(token), nil
}

// 解析 Token
func (a *AuthToken) ParseToken(token string) (int64, bool, error) {
	// 1. 将十六进制字符串解码为字节数组
	tokenBytes, err := hex.DecodeString(token)
	if err != nil {
		return 0, false, errors.New("invalid token format")
	}

	// 2. 对 token 进行解密，获取原始数据
	decryptedData, err := crypto.Aes256Decrypt(a.SecretKey, tokenBytes)
	if err != nil {
		return 0, false, err
	}

	// 3. 拆分原始数据
	parts := strings.Split(string(decryptedData), ":")
	if len(parts) != 4 {
		return 0, false, errors.New("invalid token format")
	}

	userID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, false, errors.New("invalid user id")
	}
	expiration, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, false, errors.New("invalid expiration time")
	}
	tokenSalt := parts[2]

	// 4. 校验盐值是否匹配
	if tokenSalt != a.Salt {
		return 0, false, errors.New("salt mismatch")
	}

	// 5. 校验是否过期
	currentTime := time.Now().Unix()
	if currentTime > expiration {
		return userID, false, errors.New("token expired")
	}

	return userID, true, nil
}
