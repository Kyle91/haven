// @Author Eric
// @Date 2024/6/10 1:00:00
// @Desc 错误码定义
package common

// 定义错误码
const (
	Success      = 0
	InvalidToken = 10001
	AuthFailed   = 10002
	LimitReached = 10003
	NotExist     = 10004
)

// 错误信息映射
var errorMessages = map[int]string{
	Success:      "成功",
	InvalidToken: "token无效",
	AuthFailed:   "鉴权失败",
	LimitReached: "请求次数已达上限",
	NotExist:     "数据不存在",
}

// GetErrorMessage 根据错误码获取错误信息
func GetErrorMessage(code int) string {
	if msg, exists := errorMessages[code]; exists {
		return msg
	}
	return "Unknown Error"
}
