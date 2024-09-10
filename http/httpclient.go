// @Author Eric
// @Date 2024/7/30 23:34:00
// @Desc
package http

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

func PostRequest(url, userAgent string, body []byte) ([]byte, error) {
	// 创建请求和响应对象
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// 设置请求方法、URL和Body
	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetRequestURI(url)
	req.SetBody(body)

	// 设置User-Agent
	req.Header.Set("User-Agent", userAgent)
	req.Header.SetContentType("application/json")

	// 发送请求
	client := &fasthttp.Client{}
	err := client.Do(req, resp)
	if err != nil {
		return nil, fmt.Errorf("发送请求时出错: %v", err)
	}

	// 返回响应Body
	return resp.Body(), nil
}

func GetRequest(url, userAgent string) ([]byte, error) {
	// 创建请求和响应对象
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// 设置请求方法和URL
	req.Header.SetMethod(fasthttp.MethodGet)
	req.SetRequestURI(url)

	// 设置User-Agent
	req.Header.Set("User-Agent", userAgent)

	// 发送请求
	client := &fasthttp.Client{}
	err := client.Do(req, resp)
	if err != nil {
		return nil, fmt.Errorf("发送请求时出错: %v", err)
	}

	// 返回响应Body
	return resp.Body(), nil
}
