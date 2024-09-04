// @Author Eric
// @Date 2024/9/4 16:58:00
// @Desc
package auth

import (
	"errors"
	"fmt"
	"github.com/Kyle91/haven/crypto"
	"strconv"
	"strings"
	"time"
)

type AuthToken struct {
	SecretKey []byte
	Salt      string
}

// 初始化 AuthToken 类
func NewAuthToken(secretKey, salt string) *AuthToken {
	return &AuthToken{
		SecretKey: []byte(secretKey),
		Salt:      salt,
	}
}

// 生成 Token
func (a *AuthToken) GenerateToken(userID int64, expirationTime int64) (string, error) {
	// 1. 生成过期时间戳
	expiration := time.Now().Add(time.Duration(expirationTime) * time.Second).Unix()

	// 2. 拼接原始数据: userID:expiration:salt
	rawData := fmt.Sprintf("%d:%d:%s", userID, expiration, a.Salt)

	// 3. 对原始数据进行加密，生成 token
	token, err := crypto.Aes256Encrypt(a.SecretKey, []byte(rawData))
	if err != nil {
		return "", err
	}

	return string(token), nil
}

// 解析 Token
func (a *AuthToken) ParseToken(token string) (int64, bool, error) {
	// 1. 对 token 进行解密，获取原始数据
	decryptedData, err := crypto.Aes256Decrypt(a.SecretKey, []byte(token))
	if err != nil {
		return 0, false, err
	}

	// 2. 拆分原始数据
	parts := strings.Split(string(decryptedData), ":")
	if len(parts) != 3 {
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

	// 3. 校验盐值是否匹配
	if tokenSalt != a.Salt {
		return 0, false, errors.New("salt mismatch")
	}

	// 4. 校验是否过期
	currentTime := time.Now().Unix()
	if currentTime > expiration {
		return userID, false, errors.New("token expired")
	}

	return userID, true, nil
}
