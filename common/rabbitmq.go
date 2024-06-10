// @Author Eric
// @Date 2024/6/10 9:40:00
// @Desc 消息队列相关常量定义
package common

const (
	ExchangeName     = "haven_topic_exchange"
	LoginServerQueue = "login_queue"   // 登录服务器队列
	GatewayQueue     = "gateway_queue" // 网关队列
	DataServerQueue  = "data_queue"    //数据服务器队列
)

const (
	LoginReqRoutingKey = "login.req.%s" // 发送给登录服务器请求的路由键
	LoginResRoutingKey = "login.res.%s" // 登录服务器响应的路由键
	DataReqRoutingKey  = "data.req.%s"  // 发送给数据服务器请求的路由键
	DataResRoutingKey  = "data.res.%s"  // 数据服务器响应的路由键
)

// GetQueueName 根据服务名获取队列名
func GetQueueName(serviceName string) string {
	switch serviceName {
	case "login":
		return LoginServerQueue
	case "data":
		return DataServerQueue
	case "gateway":
		return GatewayQueue
	default:
		return ""
	}
}
