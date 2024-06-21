// @Author Eric
// @Date 2024/6/21 17:05:00
// @Desc
package utils

import (
	"encoding/base64"
	"fmt"
	"github.com/Kyle91/haven/crypto"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

// GenerateSerialNumber 生成序列号
func GenerateSerialNumber() string {
	// 获取当前时间的纳秒时间戳
	nanoTime := time.Now().UnixNano()

	// 生成UUID并取前8位
	uuidStr := uuid.New().String()
	shortUUID := uuidStr[:8]

	// 生成12位随机数字
	rnd := rand.New(rand.NewSource(nanoTime))
	randomDigits := fmt.Sprintf("%012d", rnd.Int63n(1e12))

	// 拼接短UUID和12位随机数字，再加上时间戳的哈希值
	combinedStr := fmt.Sprintf("%s%s%x", shortUUID, randomDigits, nanoTime)

	// 使用base64编码并去掉填充
	base64Str := base64.RawURLEncoding.EncodeToString([]byte(combinedStr))

	// 截取前22位
	return crypto.MD5([]byte(base64Str))
}

// 输出执行时间,打印出毫秒
func TrackTime(pre time.Time) time.Duration {
	elapsed := time.Since(pre)
	fmt.Printf("elapsed: %d ms\n", elapsed.Milliseconds())

	return elapsed
}
