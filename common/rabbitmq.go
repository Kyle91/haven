// @Author Eric
// @Date 2024/6/10 9:40:00
// @Desc 消息队列相关常量定义
package common

const (
	ExchangeName     = "haven_topic_exchange"
	LoginServerQueue = "login_queue_%s"   // 登录服务器队列
	GatewayQueue     = "gateway_queue_%s" // 网关队列
	DataServerQueue  = "data_queue_%s"    //数据服务器队列
)

const (
	GatewayRoutingKey     = "gateway.%s" // 网关相关消息的路由键
	LoginServerRoutingKey = "login.%s"   // 登录服务器相关消息的路由键
	DataServerRoutingKey  = "data.%s"    // 数据服务器相关消息的路由键
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

// GetRoutingKey 根据服务名获取路由键
func GetRoutingKey(serviceName string) string {
	switch serviceName {
	case "login":
		return LoginServerRoutingKey
	case "data":
		return DataServerRoutingKey
	case "gateway":
		return GatewayRoutingKey
	default:
		return ""
	}
}
